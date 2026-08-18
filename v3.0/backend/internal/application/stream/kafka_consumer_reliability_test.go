package stream

import (
	"context"
	"errors"
	"sync"
	"testing"
	"time"

	pb "github.com/Akashpg-M/polaris/backend/api/proto/v1"
	"github.com/Akashpg-M/polaris/backend/internal/application/spatial"
	"github.com/Akashpg-M/polaris/backend/internal/core/events"
	"github.com/segmentio/kafka-go"
)

type fakeState struct {
	mu    sync.Mutex
	calls int
}

func (f *fakeState) ApplyEnvelope(*events.TelemetryEnvelope) spatial.Classification {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.calls++
	return spatial.Accepted
}

type fakeProjection struct {
	mu             sync.Mutex
	failures       int
	calls          int
	classification spatial.Classification
}

func (f *fakeProjection) Apply(context.Context, *events.TelemetryEnvelope) (spatial.Classification, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.calls++
	if f.calls <= f.failures {
		return "", errors.New("redis unavailable")
	}
	return f.classification, nil
}
func (*fakeProjection) Ready(context.Context) error { return nil }

type fakeWriter struct {
	mu       sync.Mutex
	messages []kafka.Message
}

func (f *fakeWriter) WriteMessages(_ context.Context, m ...kafka.Message) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.messages = append(f.messages, m...)
	return nil
}
func (*fakeWriter) Close() error { return nil }

type fakeReader struct {
	mu        sync.Mutex
	source    chan kafka.Message
	commits   []kafka.Message
	commitErr error
}

func (f *fakeReader) FetchMessage(ctx context.Context) (kafka.Message, error) {
	select {
	case m := <-f.source:
		return m, nil
	case <-ctx.Done():
		return kafka.Message{}, ctx.Err()
	}
}
func (f *fakeReader) CommitMessages(_ context.Context, m ...kafka.Message) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	if f.commitErr != nil {
		return f.commitErr
	}
	f.commits = append(f.commits, m...)
	return nil
}
func (*fakeReader) Close() error { return nil }

func streamEnvelope(device string, sequence uint64) *events.TelemetryEnvelope {
	now := time.Now().UTC()
	p := &pb.SpatialObject{Id: device, TenantId: "tenant-1", DeviceBootId: "boot-1", SequenceNumber: sequence,
		BootStartedAt: now.Add(-time.Minute).UnixMilli(), ObservedAt: now.UnixMilli(), SchemaVersion: 1, Type: pb.NodeType_NODE_TYPE_DRONE,
		Status: pb.NodeStatus_NODE_STATUS_ACTIVE, Lat: 13, Lon: 80, EnergyPercent: 50}
	return events.NewTelemetryEnvelope(p, now, "", "", "")
}
func kafkaEnvelope(t *testing.T, partition int, offset int64, device string) kafka.Message {
	t.Helper()
	data, err := streamEnvelope(device, uint64(offset+1)).Marshal()
	if err != nil {
		t.Fatal(err)
	}
	return kafka.Message{Topic: KafkaTelemetryTopic, Partition: partition, Offset: offset, Key: []byte("tenant-1:" + device), Value: data}
}

func TestRedisFailureDoesNotAdvanceSpatialState(t *testing.T) {
	state := &fakeState{}
	projection := &fakeProjection{failures: 2, classification: spatial.Accepted}
	writer := &fakeWriter{}
	c := &KafkaConsumer{engine: state, projector: projection, dlq: writer, maxRetries: 3}
	if !c.process(context.Background(), pendingTelemetry{message: kafkaEnvelope(t, 0, 0, "device-retry")}) {
		t.Fatal("event should succeed after retry")
	}
	if state.calls != 1 {
		t.Fatalf("spatial applied %d times before/after Redis recovery; want once", state.calls)
	}
}

func TestUnsupportedSchemaReachesDLQ(t *testing.T) {
	e := streamEnvelope("device-poison", 1)
	e.SchemaVersion = 99
	data, err := e.Marshal()
	if err != nil {
		t.Fatal(err)
	}
	writer := &fakeWriter{}
	c := &KafkaConsumer{engine: &fakeState{}, projector: &fakeProjection{}, dlq: writer, maxRetries: 1}
	if !c.process(context.Background(), pendingTelemetry{message: kafka.Message{Topic: KafkaTelemetryTopic, Partition: 2, Offset: 7, Value: data}}) {
		t.Fatal("successful DLQ is terminal")
	}
	if len(writer.messages) != 1 || string(writer.messages[0].Value) != string(data) {
		t.Fatal("original event was not preserved in DLQ")
	}
}

func TestCommitFailureLeavesSuccessfulBatchForReplay(t *testing.T) {
	reader := &fakeReader{commitErr: errors.New("broker unavailable")}
	batch := &partitionBatch{firstQueued: time.Now(), items: []pendingTelemetry{{message: kafkaEnvelope(t, 1, 10, "device-commit")}}}
	c := &KafkaConsumer{reader: reader, engine: &fakeState{}, projector: &fakeProjection{classification: spatial.Accepted}, dlq: &fakeWriter{}, maxRetries: 1}
	if c.flushPartition(context.Background(), "telemetry.ingress:1", batch) {
		t.Fatal("commit failure must not report success")
	}
	if len(batch.items) != 1 {
		t.Fatal("successful effects must remain queued for replay when commit fails")
	}
}

func TestGracefulShutdownFlushesPartitionsIndependently(t *testing.T) {
	reader := &fakeReader{source: make(chan kafka.Message, 2)}
	reader.source <- kafkaEnvelope(t, 0, 20, "device-p0")
	reader.source <- kafkaEnvelope(t, 2, 30, "device-p2")
	c := &KafkaConsumer{reader: reader, dlq: &fakeWriter{}, engine: &fakeState{}, projector: &fakeProjection{classification: spatial.Accepted}, batchSize: 1000, flushInterval: time.Hour, maxRetries: 1, done: make(chan struct{})}
	ctx, cancel := context.WithCancel(context.Background())
	go c.Start(ctx, "test")
	time.Sleep(30 * time.Millisecond)
	cancel()
	waitCtx, waitCancel := context.WithTimeout(context.Background(), time.Second)
	defer waitCancel()
	if err := c.Wait(waitCtx); err != nil {
		t.Fatal(err)
	}
	reader.mu.Lock()
	defer reader.mu.Unlock()
	if len(reader.commits) != 2 {
		t.Fatalf("committed %d partitions; want 2", len(reader.commits))
	}
	seen := map[int]bool{}
	for _, m := range reader.commits {
		seen[m.Partition] = true
	}
	if !seen[0] || !seen[2] {
		t.Fatalf("partition commits not independent: %#v", seen)
	}
}
