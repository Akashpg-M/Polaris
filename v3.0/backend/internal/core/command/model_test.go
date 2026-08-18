package command

import "testing"

func TestCommandTransitionsAndCapabilities(t *testing.T) {
	if !ValidTransition(Pending, Delivered) || !ValidTransition(Delivered, Acknowledged) || !ValidTransition(Acknowledged, Completed) {
		t.Fatal("legal command lifecycle rejected")
	}
	if ValidTransition(Completed, Pending) || ValidTransition(Pending, Acknowledged) {
		t.Fatal("illegal command lifecycle accepted")
	}
	if RequiredCapability("RELOCATE") != "receive_relocation_command" || RequiredCapability("CAPTURE_IMAGE") != "capture_image" {
		t.Fatal("command capability mapping is incorrect")
	}
}
