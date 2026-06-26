-- 1. Enable PostGIS (Crucial for advanced spatial math in Postgres)
CREATE EXTENSION IF NOT EXISTS postgis;

CREATE TABLE IF NOT EXISTS telemetry_history (
    id SERIAL PRIMARY KEY,
    tenant_id VARCHAR(50) NOT NULL,
    node_id VARCHAR(50) NOT NULL,
    
    -- V2 UPDATES: Changed to INT to natively store Protobuf Enums
    asset_type INT NOT NULL, 
    status INT NOT NULL,     
    
    -- RAW GPS
    lat DOUBLE PRECISION NOT NULL,
    lon DOUBLE PRECISION NOT NULL,
    
    -- V3 POSTGIS UPDATE: Native spatial column (EPSG:4326 is standard GPS)
    geom GEOMETRY(Point, 4326), 
    
    -- V3 PHYSICS UPDATES: Required for recreating traffic and handover simulations
    velocity_mps DOUBLE PRECISION,
    heading_deg DOUBLE PRECISION,
    
    battery INT,
    recorded_at TIMESTAMP NOT NULL
);

-- Index for fast historical reporting (e.g., "Where was Drone-1001 yesterday?")
CREATE INDEX idx_telemetry_node_time ON telemetry_history(node_id, recorded_at DESC);

-- Index for predictive queries using standard floats
CREATE INDEX idx_telemetry_predictive ON telemetry_history(recorded_at, lat, lon);

-- NEW V3 INDEX: A GIST index makes PostGIS spatial queries (like bounding boxes) lightning fast
CREATE INDEX idx_telemetry_geom ON telemetry_history USING GIST (geom);