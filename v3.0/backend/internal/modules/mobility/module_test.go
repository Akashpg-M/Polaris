package mobility

import (
	"context"
	"testing"

	"github.com/Akashpg-M/polaris/backend/internal/core/extension"
)

func TestMissingRoadGraphDegradesOptionalModule(t *testing.T) {
	cfg := DefaultConfig()
	cfg.RoadGraphPath = "does-not-exist.osm.pbf"
	cfg.Required = false
	m := New(cfg, nil)
	if err := m.Start(context.Background()); err != nil {
		t.Fatal(err)
	}
	status := m.Ready(context.Background())
	if status.State != extension.ModuleDegraded || status.Components["spatial"].State != extension.ModuleReady || status.Components["routing"].State != extension.ModuleFailed {
		t.Fatalf("unexpected degradation: %#v", status)
	}
	_ = m.Close(context.Background())
}
func TestMissingRoadGraphFailsRequiredModule(t *testing.T) {
	cfg := DefaultConfig()
	cfg.RoadGraphPath = "does-not-exist.osm.pbf"
	cfg.Required = true
	m := New(cfg, nil)
	if err := m.Start(context.Background()); err == nil {
		t.Fatal("mandatory Mobility accepted a missing graph")
	}
}
func TestConfigurationRejectsUnsafeBounds(t *testing.T) {
	cfg := DefaultConfig()
	cfg.MaxRoutedCandidates = cfg.MaxRawCandidates + 1
	if err := cfg.Validate(); err == nil {
		t.Fatal("invalid candidate fanout accepted")
	}
	cfg = DefaultConfig()
	cfg.H3ShardResolution = cfg.H3Resolution + 1
	if err := cfg.Validate(); err == nil {
		t.Fatal("invalid H3 hierarchy accepted")
	}
}
