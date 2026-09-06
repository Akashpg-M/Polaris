import { useEffect, useRef, useState } from 'react';
import L from 'leaflet';
// @ts-ignore
import 'leaflet.heat';
import 'leaflet/dist/leaflet.css';
import type { ZonePrediction } from '../types/polaris';
import { engineApi, gatewayWs } from '../config';

interface MapState {
  map: L.Map;
  markersLayer: L.LayerGroup;
  heatLayer: any;
  hotspotsLayer: L.LayerGroup;
}

// Track markers by Node ID in-memory to prevent drawing duplicates or flashing layers
interface MarkerCache {
  [nodeId: string]: {
    marker: L.Marker;
    lastUpdated: number;
  };
}

export default function MapDashboard() {
  const containerRef = useRef<HTMLDivElement | null>(null);
  const mapRef = useRef<MapState | null>(null);
  const markersCacheRef = useRef<MarkerCache>({});
  const [activeNodes, setActiveNodes] = useState(0);
  const [latestNode, setLatestNode] = useState<{ id: string; lat: number; lon: number; profile: string } | null>(null);

  useEffect(() => {
    if (!containerRef.current || mapRef.current) return;

    // Initialize Map with dark custom tiles matching your ops center theme
    const map = L.map(containerRef.current, { zoomControl: false }).setView([13.04, 80.24], 12);
    L.control.zoom({ position: 'topright' }).addTo(map);

    L.tileLayer('https://{s}.basemaps.cartocdn.com/dark_all/{z}/{x}/{y}{r}.png', {
        attribution: '&copy; CARTO'
    }).addTo(map);

    const markersLayer = L.layerGroup().addTo(map);
    const hotspotsLayer = L.layerGroup().addTo(map);

    mapRef.current = {
      map,
      markersLayer,
      heatLayer: null,
      hotspotsLayer
    };

    const waitForSize = () => {
      const size = map.getSize();
      if (size.x === 0 || size.y === 0) {
        requestAnimationFrame(waitForSize);
        return;
      }
      const heatLayer = (L as any).heatLayer([], { radius: 25, blur: 15, maxZoom: 14 }).addTo(map);
      if (mapRef.current) mapRef.current.heatLayer = heatLayer;
      map.invalidateSize();
    };

    waitForSize();

    // 1. ASYMMETRIC STREAM: WebSocket Real-Time Downstream Connection
    let ws: WebSocket | null = null;
	let disposed=false;
    const connectStreamingGateway = async () => {
	  const operatorToken=localStorage.getItem('polaris_operator_token');if(!operatorToken||disposed)return;
	  const ticketResponse=await fetch(`${engineApi}/dashboard-ticket`,{method:'POST',headers:{Authorization:`Bearer ${operatorToken}`,'X-Tenant-ID':'alpha_logistics','Content-Type':'application/json'},body:'{}'});if(!ticketResponse.ok||disposed)return;
	  const ticket=(await ticketResponse.json()).data.ticket;
	  ws = new WebSocket(`${gatewayWs}/ws/dashboard?ticket=${encodeURIComponent(ticket)}`);

      ws.onmessage = (event) => {
        try {
          // Parse downstream JSON packet coming directly via Redis PubSub from the Engine
          const node = JSON.parse(event.data);
          if (!mapRef.current) return;

          const { markersLayer, heatLayer } = mapRef.current;
          const cache = markersCacheRef.current;
          const profile = node.type === 7 ? 'STATIC' : node.type === 6 ? 'GROUND_ROBOT' : node.type === 5 ? 'AERIAL_DRONE' : 'ROAD_VEHICLE';
          const markerGlyph = profile === 'STATIC' ? '●' : profile === 'ROAD_VEHICLE' ? '◆' : profile === 'GROUND_ROBOT' ? '■' : '▲';
          const droneIcon = L.divIcon({ html: markerGlyph, className: 'custom-div-icon text-cyan-300', iconSize: [24, 24] });

          const now = Date.now();
          setLatestNode({ id: node.id, lat: node.lat, lon: node.lon, profile });

          if (cache[node.id]) {
            // Smoothly move the existing marker instead of deleting it (No layout flashing!)
            cache[node.id].marker.setLatLng([node.lat, node.lon]);
            cache[node.id].lastUpdated = now;
          } else {
            // Create a brand new marker if this node is spinning up for the first time
            const marker = L.marker([node.lat, node.lon], { icon: droneIcon })
              .bindPopup(`<strong style="color:#10b981">Node: ${node.id}</strong><br>Speed: ${node.velocity_mps?.toFixed(1) || 0} m/s<br>Battery: ${node.energy_percent}%`)
              .addTo(markersLayer);

            cache[node.id] = { marker, lastUpdated: now };
          }

          // Update Heatmap layers based on cached locations
          const heatData: [number, number, number][] = [];
          Object.keys(cache).forEach(id => {
            const position = cache[id].marker.getLatLng();
            heatData.push([position.lat, position.lng, 1]);
          });

          if (heatLayer) {
            heatLayer.setLatLngs(heatData);
          }

          setActiveNodes(Object.keys(cache).length);
        } catch (err) {
          // slog.error("Failed to parse incoming streaming telemetry frame", err);
        }
      };

      ws.onclose = () => {
        // Safe self-healing reconnection strategy if the gateway bounces
        setTimeout(connectStreamingGateway, 3000);
      };
    };

    connectStreamingGateway();

    // 2. Cold-Storage / Predictive Engine Aggregations (Keep as an interval poll)
    const fetchPredictedZones = async () => {
      try {
		const operatorToken = localStorage.getItem('polaris_operator_token');
		if (!operatorToken) return;
        const res = await fetch(`${engineApi}/zones/predicted`, {
		  headers: { Authorization: `Bearer ${operatorToken}`, 'X-Tenant-ID': 'alpha_logistics' }
		});
        const json = await res.json();

        const zones: ZonePrediction[] = Array.isArray(json.data) ? json.data : [];

        if (mapRef.current) {
          const { hotspotsLayer } = mapRef.current;
          hotspotsLayer.clearLayers();

          zones.forEach(zone => {
            L.circle([zone.Lat, zone.Lon], {
              color: '#ef4444',
              fillColor: '#ef4444',
              fillOpacity: 0.2,
              radius: zone.RadiusKm * 1000
            }).bindPopup(`<b>🤖 AI Prediction:</b> High Demand<br>Zone: ${zone.ID}`).addTo(hotspotsLayer);
          });
        }
      } catch (err) {
        // Silent fail on dev disconnects
      }
    };

    const i2 = setInterval(fetchPredictedZones, 10000);
    fetchPredictedZones(); 

    // Cache Janitor: Clean up stale nodes that haven't pinged in 10 seconds
    const janitorInterval = setInterval(() => {
      const cache = markersCacheRef.current;
      const now = Date.now();
      Object.keys(cache).forEach(id => {
        if (now - cache[id].lastUpdated > 10000) {
          if (mapRef.current) mapRef.current.markersLayer.removeLayer(cache[id].marker);
          delete cache[id];
        }
      });
      setActiveNodes(Object.keys(cache).length);
    }, 5000);

    return () => {
	  disposed=true;
      if (ws) ws.close();
      clearInterval(i2);
      clearInterval(janitorInterval);
      if (mapRef.current) {
        mapRef.current.map.remove();
        mapRef.current = null;
      }
      markersCacheRef.current = {};
    };
  }, []);

  return (
    <div className="w-full h-full relative">
      <div className="absolute top-4 left-4 z-[1000] bg-slate-800/90 p-4 rounded-lg border border-slate-700 shadow-lg">
        <div className="text-3xl font-bold text-emerald-500">{activeNodes}</div>
        <div className="text-xs text-slate-400 uppercase tracking-widest">Active Nodes</div>
      </div>
      {latestNode && (
        <div data-testid="latest-vehicle" className="absolute top-4 left-40 z-[1000] bg-slate-800/95 p-4 rounded-lg border border-slate-700 shadow-lg">
          <div className="text-xs text-slate-400 uppercase tracking-widest">Latest Vehicle</div>
          <div data-testid="latest-vehicle-id" className="font-mono font-bold text-cyan-400">{latestNode.id}</div>
          <div data-testid="latest-vehicle-coordinates" className="font-mono text-xs text-slate-300">
            {latestNode.lat.toFixed(6)}, {latestNode.lon.toFixed(6)}
          </div>
          <div className="mt-1 text-[10px] tracking-wide text-emerald-300">{latestNode.profile}</div>
        </div>
      )}
      <div ref={containerRef} className="w-full h-full" />
    </div>
  );
}
