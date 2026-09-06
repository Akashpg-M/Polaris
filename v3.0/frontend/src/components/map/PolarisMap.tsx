import { useEffect, useRef } from 'react';
import L from 'leaflet';
import 'leaflet/dist/leaflet.css';
import type { MobilityPosition, MobilityProfile, PredictedZone, RouteResult } from '../../types/domain';

export interface MapDeviceMarker {
  id: string;
  label: string;
  position: MobilityPosition;
  profile: MobilityProfile;
  connectivity?: string;
}

function glyph(profile: MobilityProfile) {
  return profile === 'AERIAL_DRONE' ? '\u25B2' : profile === 'GROUND_ROBOT' ? '\u25A0' : profile === 'STATIC' ? '\u25CF' : '\u25C6';
}

export function PolarisMap({ devices = [], selectedId, onDeviceSelect, selectionMode, onMapSelect, origin, destination, searchRadiusMeters, routes = [], zones = [], compact = false, fitKey = '' }: {
  devices?: MapDeviceMarker[];
  selectedId?: string;
  onDeviceSelect?: (id: string) => void;
  selectionMode?: 'SEARCH' | 'ORIGIN' | 'DESTINATION';
  onMapSelect?: (position: MobilityPosition) => void;
  origin?: MobilityPosition;
  destination?: MobilityPosition;
  searchRadiusMeters?: number;
  routes?: Array<{ route: RouteResult; tone?: 'primary' | 'alternate' }>;
  zones?: PredictedZone[];
  compact?: boolean;
  fitKey?: string;
}) {
  const container = useRef<HTMLDivElement | null>(null);
  const map = useRef<L.Map | null>(null);
  const overlays = useRef<L.LayerGroup | null>(null);
  const lastFit = useRef('');
  const selectCallback = useRef(onMapSelect);
  useEffect(() => { selectCallback.current = onMapSelect; }, [onMapSelect]);

  useEffect(() => {
    if (!container.current || map.current) return;
    const instance = L.map(container.current, { zoomControl: !compact, attributionControl: !compact }).setView([13.04, 80.24], compact ? 13 : 11);
    L.tileLayer('https://{s}.basemaps.cartocdn.com/light_all/{z}/{x}/{y}{r}.png', { attribution: '&copy; OpenStreetMap &copy; CARTO', maxZoom: 19 }).addTo(instance);
    overlays.current = L.layerGroup().addTo(instance);
    map.current = instance;
    const click = (event: L.LeafletMouseEvent) => selectCallback.current?.({ latitude: event.latlng.lat, longitude: event.latlng.lng });
    instance.on('click', click);
    const resize = new ResizeObserver(() => instance.invalidateSize());
    resize.observe(container.current);
    return () => { resize.disconnect(); instance.off('click', click); instance.remove(); map.current = null; overlays.current = null; };
  }, [compact]);

  useEffect(() => {
    const instance = map.current; const layer = overlays.current;
    if (!instance || !layer) return;
    layer.clearLayers();
    const bounds: L.LatLngExpression[] = [];
    devices.forEach(device => {
      const point: L.LatLngExpression = [device.position.latitude, device.position.longitude]; bounds.push(point);
      const selected = device.id === selectedId ? ' selected' : '';
      const status = device.connectivity?.toLowerCase() || 'unknown';
      L.marker(point, { keyboard: true, title: device.label, icon: L.divIcon({ className: `fleet-marker ${status}${selected}`, html: `<span>${glyph(device.profile)}</span>`, iconSize: [32, 32], iconAnchor: [16, 16] }) })
        .on('click', event => { L.DomEvent.stopPropagation(event); onDeviceSelect?.(device.id); }).addTo(layer);
    });
    if (origin) {
      const point: L.LatLngExpression = [origin.latitude, origin.longitude]; bounds.push(point);
      L.marker(point, { icon: L.divIcon({ className: 'coordinate-marker origin', html: '<span>O</span>', iconSize: [30, 30], iconAnchor: [15, 15] }), title: 'Origin' }).addTo(layer);
      if (searchRadiusMeters && searchRadiusMeters > 0) L.circle(point, { radius: searchRadiusMeters, className: 'search-radius', color: '#31d7ca', fillOpacity: .08, weight: 2 }).addTo(layer);
    }
    if (destination) {
      const point: L.LatLngExpression = [destination.latitude, destination.longitude]; bounds.push(point);
      L.marker(point, { icon: L.divIcon({ className: 'coordinate-marker destination', html: '<span>D</span>', iconSize: [30, 30], iconAnchor: [15, 15] }), title: 'Destination' }).addTo(layer);
    }
    routes.forEach(({ route, tone = 'primary' }) => {
      const points = route.waypoints.map(point => [point.latitude, point.longitude] as [number, number]);
      if (points.length) { points.forEach(point => bounds.push(point)); L.polyline(points, { color: tone === 'alternate' ? '#ae82ff' : '#31d7ca', weight: tone === 'alternate' ? 4 : 5, opacity: .9, dashArray: tone === 'alternate' ? '9 7' : undefined, className: `route-line ${tone}` }).addTo(layer); }
    });
    zones.forEach(zone => {
      const point: L.LatLngExpression = [zone.lat, zone.lon]; bounds.push(point);
      L.circle(point, { radius: zone.radius_km * 1000, color: '#f4b860', fillColor: '#f4b860', fillOpacity: .12, weight: 2, dashArray: '7 6' }).bindTooltip(zone.id).addTo(layer);
    });
    if (bounds.length && fitKey && lastFit.current !== fitKey) {
      instance.fitBounds(L.latLngBounds(bounds).pad(.2), { maxZoom: compact ? 14 : 16, animate: false });
      lastFit.current = fitKey;
    }
  }, [devices, selectedId, onDeviceSelect, origin, destination, searchRadiusMeters, routes, zones, compact, fitKey]);

  return <div className={`polaris-map fleet-map-canvas ${compact ? 'compact' : ''} ${selectionMode ? 'selecting' : ''}`} ref={container} aria-label={selectionMode ? `Map selecting ${selectionMode.toLowerCase()} coordinates` : 'Polaris map'} />;
}
