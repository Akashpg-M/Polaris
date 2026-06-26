package repository

import (
	"context"
	"encoding/json"
	"log/slog"
)

type MockEventPublisher struct{}

func NewMockEventPublisher() *MockEventPublisher {
	return &MockEventPublisher{}
}

// PublishEvent marshals the actor's internal mutation event record and outputs to system log
func (m *MockEventPublisher) PublishEvent(ctx context.Context, topic string, event interface{}) error {
	bytes, err := json.Marshal(event)
	if err != nil {
		return err
	}
	
	// Print directly as a structured info log to verify our asynchronous projection flow is humming!
	slog.Info("[Mock Event Mesh Bus]", "topic", topic, "event_stream_payload", string(bytes))
	return nil
}