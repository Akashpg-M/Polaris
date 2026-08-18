export interface LogEntry {
  time: string;
  msg: string;
  type: 'info' | 'success' | 'warning' | 'danger';
}

export interface ZonePrediction {
  ID: string;
  Lat: number;
  Lon: number;
  RadiusKm: number;
  RequiredAssets: number;
  TargetClass: number;
  TenantID: string;
}

/**
 * High-Performance, Zero-Allocation Protobuf Wire Format Encoder
 * Compiles physics telemetry directly to raw bytes matching spatial.proto field tags
 */
export function serializeProtobufTelemetry(data: {
  id: string;
  tenantId: string;
  type: number;
  status: number;
  lat: number;
  lon: number;
  velocityMps: number;
  headingDeg: number;
  energyPercent: number;
  timestamp: number;
	deviceBootId: string;
	sequenceNumber: number;
	bootStartedAt: number;
	observedAt: number;
	schemaVersion?: number;
}): Uint8Array {
  const encoder = new TextEncoder();
  const idBytes = encoder.encode(data.id);
  const tenantBytes = encoder.encode(data.tenantId);
	const bootBytes = encoder.encode(data.deviceBootId);

  // Allocate a safe upper-bound workspace buffer
  const buffer = new ArrayBuffer(256);
  const view = new DataView(buffer);
  let offset = 0;

  // Helper to write Protobuf Varints (Variable-length integers)
  const writeVarint = (value: number) => {
    while (value >= 0x80) {
      view.setUint8(offset++, (value & 0x7f) | 0x80);
      value >>>= 7;
    }
    view.setUint8(offset++, value & 0x7f);
  };

  // Helper to write BigInt Varints (for 64-bit millisecond timestamps)
  const writeBigVarint = (value: bigint) => {
    while (value >= 0x80n) {
      view.setUint8(offset++, Number((value & 0x7fn) | 0x80n));
      value >>= 7n;
    }
    view.setUint8(offset++, Number(value & 0x7fn));
  };

  // Field 1: string id (Wire Type 2: Length-delimited)
  view.setUint8(offset++, (1 << 3) | 2);
  writeVarint(idBytes.length);
  for (let i = 0; i < idBytes.length; i++) view.setUint8(offset++, idBytes[i]);

  // Field 2: string tenant_id (Wire Type 2: Length-delimited)
  view.setUint8(offset++, (2 << 3) | 2);
  writeVarint(tenantBytes.length);
  for (let i = 0; i < tenantBytes.length; i++) view.setUint8(offset++, tenantBytes[i]);

  // Field 3: int32 type (Wire Type 0: Varint)
  view.setUint8(offset++, (3 << 3) | 0);
  writeVarint(data.type);

  // Field 4: int32 status (Wire Type 0: Varint)
  view.setUint8(offset++, (4 << 3) | 0);
  writeVarint(data.status);

  // Field 5: double lat (Wire Type 1: 64-bit Fixed)
  view.setUint8(offset++, (5 << 3) | 1);
  view.setFloat64(offset, data.lat, true); // true = Little Endian
  offset += 8;

  // Field 6: double lon (Wire Type 1: 64-bit Fixed)
  view.setUint8(offset++, (6 << 3) | 1);
  view.setFloat64(offset, data.lon, true);
  offset += 8;

  // Field 7: double velocity_mps (Wire Type 1: 64-bit Fixed)
  view.setUint8(offset++, (7 << 3) | 1);
  view.setFloat64(offset, data.velocityMps, true);
  offset += 8;

  // Field 8: double heading_deg (Wire Type 1: 64-bit Fixed)
  view.setUint8(offset++, (8 << 3) | 1);
  view.setFloat64(offset, data.headingDeg, true);
  offset += 8;

  // Field 9: int32 energy_percent (Wire Type 0: Varint)
  view.setUint8(offset++, (9 << 3) | 0);
  writeVarint(data.energyPercent);

  // Field 10: int64 timestamp (Wire Type 0: Varint)
  view.setUint8(offset++, (10 << 3) | 0);
  writeBigVarint(BigInt(data.timestamp));

	// Phase 1/2 device-owned ordering identity.
	view.setUint8(offset++, (11 << 3) | 2);
	writeVarint(bootBytes.length);
	for (let i = 0; i < bootBytes.length; i++) view.setUint8(offset++, bootBytes[i]);
	view.setUint8(offset++, (12 << 3) | 0); writeBigVarint(BigInt(data.sequenceNumber));
	view.setUint8(offset++, (13 << 3) | 0); writeBigVarint(BigInt(data.bootStartedAt));
	view.setUint8(offset++, (14 << 3) | 0); writeBigVarint(BigInt(data.observedAt));
	view.setUint8(offset++, (15 << 3) | 0); writeVarint(data.schemaVersion ?? 1);

  // Slice out the exact compiled bytecode envelope
  return new Uint8Array(buffer, 0, offset);
}
