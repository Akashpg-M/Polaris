package auth

import "testing"

func TestGeneratedTokenIsHashedAndVerifiable(t *testing.T) {
	raw, prefix, hash, err := GenerateToken("dev")
	if err != nil {
		t.Fatal(err)
	}
	if len(hash) != 32 || len(raw) < 80 || prefix == "" {
		t.Fatalf("weak token shape: %d %d", len(raw), len(hash))
	}
	parsed, err := TokenPrefix(raw)
	if err != nil || parsed != prefix {
		t.Fatalf("prefix: %q %v", parsed, err)
	}
	if !Verify(raw, hash) {
		t.Fatal("valid token did not verify")
	}
	if Verify(raw+"x", hash) {
		t.Fatal("modified token verified")
	}
}
func TestRolePermissions(t *testing.T) {
	if !Can(PlatformAdmin, "mutate") || !Can(TenantAdmin, "mutate") || Can(Operator, "mutate") || Can(Viewer, "mutate") {
		t.Fatal("mutation permission matrix violated")
	}
	if !Can(Viewer, "read") || Can(Viewer, "audit") {
		t.Fatal("viewer permission matrix violated")
	}
	if !Can(Operator, "orchestrate") || Can(Viewer, "orchestrate") || Can(Operator, "admin_retry") || !Can(TenantAdmin, "admin_retry") {
		t.Fatal("orchestration permission matrix violated")
	}
}
