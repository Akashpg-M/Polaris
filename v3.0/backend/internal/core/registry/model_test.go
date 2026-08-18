package registry

import "testing"

func TestLifecycleTransitions(t *testing.T) {
	allowed := [][2]string{{"REGISTERED", "ACTIVE"}, {"ACTIVE", "SUSPENDED"}, {"SUSPENDED", "ACTIVE"}, {"ACTIVE", "DECOMMISSIONED"}}
	for _, v := range allowed {
		if !ValidTransition(v[0], v[1]) {
			t.Fatalf("expected %s -> %s", v[0], v[1])
		}
	}
	denied := [][2]string{{"DECOMMISSIONED", "ACTIVE"}, {"ACTIVE", "REGISTERED"}, {"REGISTERED", "SUSPENDED"}}
	for _, v := range denied {
		if ValidTransition(v[0], v[1]) {
			t.Fatalf("unexpected %s -> %s", v[0], v[1])
		}
	}
}
