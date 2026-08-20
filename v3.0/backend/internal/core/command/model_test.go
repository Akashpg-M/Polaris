package command

import (
	"encoding/json"
	"testing"
	"time"
)

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

func TestDeliveryObservationIsVolatile(t *testing.T) {
	record := Record{CommandID: "command", TenantID: "tenant", DeviceID: "device", Payload: json.RawMessage(`{}`)}
	envelope := record.Envelope()
	if envelope.DeliveryObservation != nil {
		t.Fatal("durable record unexpectedly retained volatile delivery timing")
	}
	envelope.DeliveryObservation = &DeliveryObservation{RelayPublishedAt: time.Now().UTC()}
	if replay := record.Envelope(); replay.DeliveryObservation != nil || replay.CommandID != envelope.CommandID || string(replay.Payload) != string(envelope.Payload) {
		t.Fatal("delivery observation mutated durable command identity")
	}
}
