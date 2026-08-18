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

-- Phase 2: durable registry, credentials, audit and transactional outbox.
CREATE EXTENSION IF NOT EXISTS pgcrypto;
CREATE TABLE IF NOT EXISTS tenants (
  tenant_id TEXT PRIMARY KEY, display_name TEXT NOT NULL,
  status TEXT NOT NULL CHECK(status IN ('ACTIVE','SUSPENDED','DEACTIVATED')),
  metadata JSONB NOT NULL DEFAULT '{}', created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(), updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE TABLE IF NOT EXISTS projects (
  project_id UUID PRIMARY KEY, tenant_id TEXT NOT NULL REFERENCES tenants(tenant_id), name TEXT NOT NULL,
  description TEXT, status TEXT NOT NULL DEFAULT 'ACTIVE', metadata JSONB NOT NULL DEFAULT '{}',
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(), updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(), UNIQUE(tenant_id,name)
);
CREATE TABLE IF NOT EXISTS device_types (
  device_type_id TEXT PRIMARY KEY, display_name TEXT NOT NULL, category TEXT NOT NULL, description TEXT,
  telemetry_schema TEXT NOT NULL DEFAULT 'spatial.v1', metadata JSONB NOT NULL DEFAULT '{}'
);
CREATE TABLE IF NOT EXISTS devices (
  tenant_id TEXT NOT NULL REFERENCES tenants(tenant_id), device_id TEXT NOT NULL, project_id UUID REFERENCES projects(project_id),
  device_type_id TEXT NOT NULL REFERENCES device_types(device_type_id), display_name TEXT NOT NULL,
  lifecycle_status TEXT NOT NULL CHECK(lifecycle_status IN ('REGISTERED','ACTIVE','SUSPENDED','DECOMMISSIONED')),
  firmware_version TEXT, software_version TEXT, model_version TEXT, metadata JSONB NOT NULL DEFAULT '{}',
  registered_at TIMESTAMPTZ NOT NULL DEFAULT NOW(), updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(), deactivated_at TIMESTAMPTZ,
  PRIMARY KEY(tenant_id,device_id)
);
CREATE TABLE IF NOT EXISTS capabilities (
  capability_id TEXT PRIMARY KEY, display_name TEXT NOT NULL, description TEXT, schema JSONB NOT NULL DEFAULT '{}'
);
CREATE TABLE IF NOT EXISTS device_capabilities (
  tenant_id TEXT NOT NULL, device_id TEXT NOT NULL, capability_id TEXT NOT NULL REFERENCES capabilities(capability_id),
  configuration JSONB NOT NULL DEFAULT '{}', enabled BOOLEAN NOT NULL DEFAULT TRUE, created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  PRIMARY KEY(tenant_id,device_id,capability_id), FOREIGN KEY(tenant_id,device_id) REFERENCES devices(tenant_id,device_id)
);
CREATE TABLE IF NOT EXISTS device_credentials (
  credential_id UUID PRIMARY KEY, tenant_id TEXT NOT NULL, device_id TEXT NOT NULL, token_prefix TEXT NOT NULL UNIQUE,
  token_hash BYTEA NOT NULL, status TEXT NOT NULL CHECK(status IN ('ACTIVE','REVOKED','EXPIRED')),
  issued_at TIMESTAMPTZ NOT NULL DEFAULT NOW(), expires_at TIMESTAMPTZ, last_used_at TIMESTAMPTZ, revoked_at TIMESTAMPTZ,
  FOREIGN KEY(tenant_id,device_id) REFERENCES devices(tenant_id,device_id)
);
CREATE TABLE IF NOT EXISTS operator_api_keys (
  api_key_id UUID PRIMARY KEY, tenant_id TEXT REFERENCES tenants(tenant_id), name TEXT NOT NULL, token_prefix TEXT NOT NULL UNIQUE,
  token_hash BYTEA NOT NULL, role TEXT NOT NULL CHECK(role IN ('PLATFORM_ADMIN','TENANT_ADMIN','OPERATOR','VIEWER')),
  status TEXT NOT NULL CHECK(status IN ('ACTIVE','REVOKED','EXPIRED')), issued_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  expires_at TIMESTAMPTZ, last_used_at TIMESTAMPTZ, revoked_at TIMESTAMPTZ
);
CREATE TABLE IF NOT EXISTS connection_tickets (
  ticket_prefix TEXT PRIMARY KEY, ticket_hash BYTEA NOT NULL, tenant_id TEXT NOT NULL, device_id TEXT NOT NULL,
  credential_id UUID NOT NULL, expires_at TIMESTAMPTZ NOT NULL, consumed_at TIMESTAMPTZ,
  FOREIGN KEY(tenant_id,device_id) REFERENCES devices(tenant_id,device_id)
);
CREATE TABLE IF NOT EXISTS operator_connection_tickets (
  ticket_prefix TEXT PRIMARY KEY, ticket_hash BYTEA NOT NULL, api_key_id UUID NOT NULL,
  tenant_id TEXT, role TEXT NOT NULL, expires_at TIMESTAMPTZ NOT NULL, consumed_at TIMESTAMPTZ,
  FOREIGN KEY(api_key_id) REFERENCES operator_api_keys(api_key_id)
);
CREATE TABLE IF NOT EXISTS outbox_events (
  outbox_id UUID PRIMARY KEY, aggregate_type TEXT NOT NULL, aggregate_id TEXT NOT NULL, tenant_id TEXT NOT NULL,
  event_id TEXT NOT NULL UNIQUE, event_type TEXT NOT NULL, schema_version INTEGER NOT NULL, payload JSONB NOT NULL,
  status TEXT NOT NULL DEFAULT 'PENDING', attempt_count INTEGER NOT NULL DEFAULT 0, next_attempt_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(), published_at TIMESTAMPTZ, last_error TEXT
);
CREATE INDEX IF NOT EXISTS idx_outbox_pending ON outbox_events(status,next_attempt_at,created_at);
CREATE TABLE IF NOT EXISTS audit_events (
  audit_id UUID PRIMARY KEY, tenant_id TEXT, actor_type TEXT NOT NULL, actor_id TEXT NOT NULL, action TEXT NOT NULL,
  resource_type TEXT NOT NULL, resource_id TEXT NOT NULL, request_id TEXT, outcome TEXT NOT NULL,
  metadata JSONB NOT NULL DEFAULT '{}', created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX IF NOT EXISTS idx_audit_tenant_time ON audit_events(tenant_id,created_at DESC);

INSERT INTO device_types(device_type_id,display_name,category,description) VALUES
 ('delivery_drone','Delivery drone','MOBILE','Spatial delivery aircraft'),
 ('ground_robot','Ground robot','MOBILE','Autonomous ground robot'),
 ('connected_vehicle','Connected vehicle','MOBILE','Connected road vehicle'),
 ('fixed_iot_sensor','Fixed IoT sensor','STATIC','Fixed telemetry sensor') ON CONFLICT DO NOTHING;
INSERT INTO capabilities(capability_id,display_name,description) VALUES
 ('navigate','Navigate','Autonomous navigation'),
 ('receive_relocation_command','Receive relocation command','Accept relocation directives'),
 ('capture_image','Capture image','Capture still imagery'),
 ('measure_temperature','Measure temperature','Report ambient temperature') ON CONFLICT DO NOTHING;

-- Phase 3: durable task, assignment and command orchestration.
CREATE TABLE IF NOT EXISTS tasks (
  task_id UUID PRIMARY KEY, tenant_id TEXT NOT NULL REFERENCES tenants(tenant_id), project_id UUID REFERENCES projects(project_id),
  task_type TEXT NOT NULL, status TEXT NOT NULL CHECK(status IN ('PENDING','ASSIGNING','ASSIGNED','IN_PROGRESS','COMPLETED','FAILED','CANCELLED','EXPIRED')),
  priority TEXT NOT NULL CHECK(priority IN ('LOW','NORMAL','HIGH','CRITICAL')), requirements JSONB NOT NULL DEFAULT '{}', target JSONB NOT NULL DEFAULT '{}',
  assigned_device_id TEXT, correlation_id TEXT NOT NULL, created_by UUID NOT NULL REFERENCES operator_api_keys(api_key_id),
  version BIGINT NOT NULL DEFAULT 1, created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(), updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  assigned_at TIMESTAMPTZ, started_at TIMESTAMPTZ, completed_at TIMESTAMPTZ, failed_at TIMESTAMPTZ,
  expires_at TIMESTAMPTZ NOT NULL, failure_reason TEXT,
  FOREIGN KEY(tenant_id,assigned_device_id) REFERENCES devices(tenant_id,device_id)
);
CREATE INDEX IF NOT EXISTS idx_tasks_tenant_status ON tasks(tenant_id,status,priority,created_at);

CREATE TABLE IF NOT EXISTS device_assignments (
  assignment_id UUID PRIMARY KEY, tenant_id TEXT NOT NULL, device_id TEXT NOT NULL, task_id UUID NOT NULL REFERENCES tasks(task_id),
  status TEXT NOT NULL CHECK(status IN ('ACTIVE','RELEASED','EXPIRED')), lease_expires_at TIMESTAMPTZ NOT NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(), updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  FOREIGN KEY(tenant_id,device_id) REFERENCES devices(tenant_id,device_id), UNIQUE(task_id)
);
CREATE UNIQUE INDEX IF NOT EXISTS uq_device_active_assignment ON device_assignments(tenant_id,device_id) WHERE status='ACTIVE';

CREATE TABLE IF NOT EXISTS device_command_sequences (
  tenant_id TEXT NOT NULL, device_id TEXT NOT NULL, last_sequence BIGINT NOT NULL DEFAULT 0,
  PRIMARY KEY(tenant_id,device_id), FOREIGN KEY(tenant_id,device_id) REFERENCES devices(tenant_id,device_id)
);

CREATE TABLE IF NOT EXISTS commands (
  command_id UUID PRIMARY KEY, tenant_id TEXT NOT NULL, device_id TEXT NOT NULL, task_id UUID NOT NULL REFERENCES tasks(task_id),
  command_type TEXT NOT NULL, payload JSONB NOT NULL, status TEXT NOT NULL CHECK(status IN ('PENDING','DELIVERED','ACKNOWLEDGED','COMPLETED','FAILED','EXPIRED','CANCELLED')),
  sequence_number BIGINT NOT NULL, correlation_id TEXT NOT NULL, causation_id TEXT NOT NULL, attempt_count INTEGER NOT NULL DEFAULT 0,
  max_attempts INTEGER NOT NULL DEFAULT 5 CHECK(max_attempts>0), version BIGINT NOT NULL DEFAULT 1,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(), available_at TIMESTAMPTZ NOT NULL DEFAULT NOW(), sent_at TIMESTAMPTZ,
  acknowledged_at TIMESTAMPTZ, completed_at TIMESTAMPTZ, expires_at TIMESTAMPTZ NOT NULL,
  ack_status TEXT, result JSONB, last_error TEXT,
  FOREIGN KEY(tenant_id,device_id) REFERENCES devices(tenant_id,device_id), UNIQUE(tenant_id,device_id,sequence_number)
);
CREATE INDEX IF NOT EXISTS idx_commands_dispatch ON commands(status,available_at,expires_at);
CREATE INDEX IF NOT EXISTS idx_commands_device_order ON commands(tenant_id,device_id,sequence_number);
CREATE INDEX IF NOT EXISTS idx_commands_task ON commands(tenant_id,task_id);

CREATE TABLE IF NOT EXISTS command_attempts (
  attempt_id UUID PRIMARY KEY, command_id UUID NOT NULL REFERENCES commands(command_id), attempt_number INTEGER NOT NULL,
  gateway_id TEXT NOT NULL, ownership_epoch BIGINT NOT NULL, started_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  completed_at TIMESTAMPTZ, result TEXT, error TEXT, UNIQUE(command_id,attempt_number)
);
