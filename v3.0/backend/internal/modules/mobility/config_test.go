package mobility

import "testing"

func TestTrafficScopeIsExplicitAndValidated(t *testing.T) {
	cfg := DefaultConfig()
	if cfg.TrafficScope != "SHARED_TRUSTED" || cfg.Validate() != nil {
		t.Fatalf("unexpected default traffic policy: %#v", cfg)
	}
	cfg.TrafficScope = "TENANT_PRIVATE"
	if cfg.Validate() == nil {
		t.Fatal("unimplemented tenant-private traffic policy was silently accepted")
	}
}
