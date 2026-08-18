-- 1. Enable PostGIS (Crucial for advanced spatial math in Postgres)
CREATE EXTENSION IF NOT EXISTS postgis;

CREATE TABLE IF NOT EXISTS telemetry_history (
    id SERIAL PRIMARY KEY,
    tenant_id VARCHAR(50) NOT NULL,
    device_id VARCHAR(128) NOT NULL,
	device_boot_id VARCHAR(128) NOT NULL,
	sequence_number BIGINT NOT NULL,
	event_id VARCHAR(64) NOT NULL,
    
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
	recorded_at TIMESTAMP NOT NULL,
	observed_at TIMESTAMPTZ NOT NULL,
	ingested_at TIMESTAMPTZ NOT NULL,
	schema_version INT NOT NULL,
	correlation_id VARCHAR(128) NOT NULL
);

-- Forward migration for Phase 0 databases. Docker entrypoint scripts only run
-- on a fresh volume; smoke-test.ps1 reapplies this idempotent file explicitly.
DO $$ BEGIN
  IF EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name='telemetry_history' AND column_name='node_id')
     AND NOT EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name='telemetry_history' AND column_name='device_id') THEN
    ALTER TABLE telemetry_history RENAME COLUMN node_id TO device_id;
  END IF;
END $$;
ALTER TABLE telemetry_history ADD COLUMN IF NOT EXISTS event_id VARCHAR(64);
ALTER TABLE telemetry_history ADD COLUMN IF NOT EXISTS device_boot_id VARCHAR(128);
ALTER TABLE telemetry_history ADD COLUMN IF NOT EXISTS sequence_number BIGINT;
ALTER TABLE telemetry_history ADD COLUMN IF NOT EXISTS observed_at TIMESTAMPTZ;
ALTER TABLE telemetry_history ADD COLUMN IF NOT EXISTS ingested_at TIMESTAMPTZ;
ALTER TABLE telemetry_history ADD COLUMN IF NOT EXISTS schema_version INT;
ALTER TABLE telemetry_history ADD COLUMN IF NOT EXISTS correlation_id VARCHAR(128);
UPDATE telemetry_history SET
  event_id=COALESCE(event_id, 'legacy:' || id),
  device_boot_id=COALESCE(device_boot_id, 'legacy'),
  sequence_number=COALESCE(sequence_number, id),
  observed_at=COALESCE(observed_at, recorded_at AT TIME ZONE 'UTC'),
  ingested_at=COALESCE(ingested_at, recorded_at AT TIME ZONE 'UTC'),
  schema_version=COALESCE(schema_version, 0),
  correlation_id=COALESCE(correlation_id, 'legacy:' || id);
ALTER TABLE telemetry_history ALTER COLUMN event_id SET NOT NULL;
ALTER TABLE telemetry_history ALTER COLUMN device_boot_id SET NOT NULL;
ALTER TABLE telemetry_history ALTER COLUMN sequence_number SET NOT NULL;
ALTER TABLE telemetry_history ALTER COLUMN observed_at SET NOT NULL;
ALTER TABLE telemetry_history ALTER COLUMN ingested_at SET NOT NULL;
ALTER TABLE telemetry_history ALTER COLUMN schema_version SET NOT NULL;
ALTER TABLE telemetry_history ALTER COLUMN correlation_id SET NOT NULL;
CREATE UNIQUE INDEX IF NOT EXISTS uq_telemetry_event_id ON telemetry_history(event_id);
CREATE UNIQUE INDEX IF NOT EXISTS uq_telemetry_device_sequence ON telemetry_history(tenant_id, device_id, device_boot_id, sequence_number);

-- Index for fast historical reporting (e.g., "Where was Drone-1001 yesterday?")
CREATE INDEX IF NOT EXISTS idx_telemetry_node_time ON telemetry_history(device_id, recorded_at DESC);

-- Index for predictive queries using standard floats
CREATE INDEX IF NOT EXISTS idx_telemetry_predictive ON telemetry_history(recorded_at, lat, lon);

-- NEW V3 INDEX: A GIST index makes PostGIS spatial queries (like bounding boxes) lightning fast
CREATE INDEX IF NOT EXISTS idx_telemetry_geom ON telemetry_history USING GIST (geom);
