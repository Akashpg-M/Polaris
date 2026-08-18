package task

import "testing"

func TestTaskTransitions(t *testing.T) {
	allowed := [][2]Status{{Pending, Assigning}, {Assigning, Assigned}, {Assigned, InProgress}, {InProgress, Completed}, {Assigned, Failed}}
	for _, pair := range allowed {
		if !ValidTransition(pair[0], pair[1]) {
			t.Fatalf("expected %s -> %s", pair[0], pair[1])
		}
	}
	if ValidTransition(Completed, Pending) || ValidTransition(Cancelled, InProgress) {
		t.Fatal("terminal task state was reversible")
	}
}
