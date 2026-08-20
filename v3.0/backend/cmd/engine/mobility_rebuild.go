package main

import (
	"context"
	"encoding/json"
	"strconv"
	"strings"
	"time"

	"github.com/Akashpg-M/polaris/backend/internal/adapter/repository"
	twincore "github.com/Akashpg-M/polaris/backend/internal/core/twin"
	"github.com/Akashpg-M/polaris/backend/internal/modules/mobility/model"
	"github.com/redis/go-redis/v9"
)

type spatialPayloadV1 struct {
	Latitude        float64               `json:"latitude"`
	Longitude       float64               `json:"longitude"`
	HeadingDegrees  *float64              `json:"heading_degrees"`
	SpeedMPS        *float64              `json:"speed_mps"`
	MobilityProfile model.MobilityProfile `json:"mobility_profile"`
}

func mobilityRebuildLoader(client *redis.Client, store *repository.RegistryStore) func(context.Context) ([]model.SpatialState, error) {
	return func(ctx context.Context) ([]model.SpatialState, error) {
		var cursor uint64
		out := []model.SpatialState{}
		for {
			keys, next, err := client.Scan(ctx, cursor, "polaris:twin:*", 200).Result()
			if err != nil {
				return nil, err
			}
			for _, key := range keys {
				if strings.HasSuffix(key, ":retired_boots") {
					continue
				}
				identity := strings.TrimPrefix(key, "polaris:twin:")
				parts := strings.SplitN(identity, ":", 2)
				if len(parts) != 2 {
					continue
				}
				tenantID, deviceID := parts[0], parts[1]
				values, err := client.HGetAll(ctx, key).Result()
				if err != nil || values["connectivity_status"] != "ONLINE" {
					continue
				}
				device, err := store.GetDevice(ctx, tenantID, deviceID)
				if err != nil || device.LifecycleStatus != "ACTIVE" {
					continue
				}
				tenant, err := store.GetTenant(ctx, tenantID)
				if err != nil || tenant.Status != "ACTIVE" {
					continue
				}
				var envelope twincore.ComponentEnvelope
				if json.Unmarshal([]byte(values["component:spatial/v1"]), &envelope) != nil {
					continue
				}
				var payload spatialPayloadV1
				if json.Unmarshal(envelope.Payload, &payload) != nil {
					continue
				}
				bootStarted, _ := strconv.ParseInt(values["boot_started_at"], 10, 64)
				out = append(out, model.SpatialState{TenantID: tenantID, DeviceID: deviceID, ReportedPosition: model.Position{Latitude: payload.Latitude, Longitude: payload.Longitude}, HeadingDegrees: payload.HeadingDegrees, SpeedMPS: payload.SpeedMPS, MobilityProfile: payload.MobilityProfile, ObservedAt: envelope.ObservedAt, BootID: envelope.BootID, BootStartedAt: time.UnixMilli(bootStarted).UTC(), SequenceNumber: envelope.SequenceNumber})
			}
			cursor = next
			if cursor == 0 {
				break
			}
		}
		return out, nil
	}
}
