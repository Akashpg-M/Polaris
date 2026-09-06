package stream

import (
	"context"
	"fmt"
	"log/slog"
	"sort"
	"sync/atomic"
	"time"

	"github.com/Akashpg-M/polaris/backend/internal/application/spatial"
	"github.com/Akashpg-M/polaris/backend/internal/core/events"
	"github.com/redis/go-redis/v9"
	"github.com/segmentio/kafka-go"
)

const DashboardUpdatesChannel = "spatial:updates"

type pendingTelemetry struct{ message kafka.Message }
type partitionBatch struct {
	items       []pendingTelemetry
	firstQueued time.Time
}

type telemetryReader interface {
	FetchMessage(context.Context) (kafka.Message, error)
	CommitMessages(context.Context, ...kafka.Message) error
	Close() error
}
type telemetryWriter interface {
	WriteMessages(context.Context, ...kafka.Message) error
	Close() error
}
type stateApplier interface {
	ApplyEnvelope(*events.TelemetryEnvelope) spatial.Classification
}
type latestProjector interface {
	Apply(context.Context, *events.TelemetryEnvelope) (spatial.Classification, error)
	Ready(context.Context) error
}

type KafkaConsumer struct {
	reader        telemetryReader
	dlq           telemetryWriter
	brokerURL     string
	engine        stateApplier
	projector     latestProjector
	batchSize     int
	flushInterval time.Duration
	maxRetries    int
	done          chan struct{}
	lastProgress  atomic.Int64
}

func NewKafkaConsumer(brokerURL string, engine stateApplier, redisClient *redis.Client) *KafkaConsumer {
	c := &KafkaConsumer{
		reader:    kafka.NewReader(kafka.ReaderConfig{Brokers: []string{brokerURL}, Topic: KafkaTelemetryTopic, GroupID: "polaris_engine_group", CommitInterval: 0}),
		dlq:       &kafka.Writer{Addr: kafka.TCP(brokerURL), Topic: DeadLetterTopic, Balancer: &kafka.Hash{}},
		brokerURL: brokerURL, engine: engine, projector: NewRedisProjector(redisClient),
		batchSize: 1000, flushInterval: 150 * time.Millisecond, maxRetries: 5, done: make(chan struct{}),
	}
	c.lastProgress.Store(time.Now().UnixMilli())
	return c
}

func partitionID(message kafka.Message) string {
	return fmt.Sprintf("%s:%d", message.Topic, message.Partition)
}

func (c *KafkaConsumer) sendToDLQ(ctx context.Context, msg kafka.Message, reason string) error {
	return c.dlq.WriteMessages(ctx, kafka.Message{Key: msg.Key, Value: msg.Value, Headers: []kafka.Header{
		{Key: "error_reason", Value: []byte(reason)}, {Key: "source_topic", Value: []byte(msg.Topic)},
		{Key: "source_partition", Value: []byte(fmt.Sprint(msg.Partition))}, {Key: "source_offset", Value: []byte(fmt.Sprint(msg.Offset))},
		{Key: "failed_at", Value: []byte(time.Now().UTC().Format(time.RFC3339Nano))},
	}})
}

func (c *KafkaConsumer) process(ctx context.Context, item pendingTelemetry) bool {
	envelope, err := events.Unmarshal(item.message.Value)
	if err != nil {
		if dlqErr := c.sendToDLQ(ctx, item.message, err.Error()); dlqErr != nil {
			slog.Error("permanent telemetry failure could not reach DLQ", "partition", item.message.Partition, "offset", item.message.Offset, "error", dlqErr)
			return false
		}
		slog.Warn("permanent telemetry failure sent to DLQ", "partition", item.message.Partition, "offset", item.message.Offset, "error", err)
		return true
	}

	var lastErr error
	for attempt := 1; attempt <= c.maxRetries; attempt++ {
		redisClass, projectionErr := c.projector.Apply(ctx, envelope)
		if projectionErr == nil {
			memoryClass := spatial.Classification("NOT_APPLIED")
			// A Redis DUPLICATE may be the replay that rebuilds an empty engine
			// after restart. Stale/retired/conflicting events must never enter it.
			if redisClass == spatial.Accepted || redisClass == spatial.NewBoot || redisClass == spatial.Duplicate {
				memoryClass = c.engine.ApplyEnvelope(envelope)
			}
			// Per-event classification is useful for diagnosis but is too noisy for
			// the steady-state INFO path under fleet load.
			slog.Debug("telemetry state classified", "event_id", envelope.EventID, "spatial", memoryClass, "redis", redisClass)
			return true
		}
		lastErr = projectionErr
		slog.Warn("transient Redis projection failure", "event_id", envelope.EventID, "attempt", attempt, "error", projectionErr)
		if attempt < c.maxRetries {
			select {
			case <-time.After(time.Duration(attempt*50) * time.Millisecond):
			case <-ctx.Done():
				return false
			}
		}
	}
	if err := c.sendToDLQ(ctx, item.message, "retry_exhausted: "+lastErr.Error()); err != nil {
		slog.Error("retry-exhausted telemetry could not reach DLQ", "offset", item.message.Offset, "error", err)
		return false
	}
	return true
}

// flushPartition commits only the highest contiguous successfully processed offset.
func (c *KafkaConsumer) flushPartition(ctx context.Context, key string, batch *partitionBatch) bool {
	if len(batch.items) == 0 {
		return true
	}
	sort.Slice(batch.items, func(i, j int) bool { return batch.items[i].message.Offset < batch.items[j].message.Offset })
	slog.Info("partition batch flush started", "partition", key, "messages", len(batch.items), "queue_age_ms", time.Since(batch.firstQueued).Milliseconds())
	succeeded := 0
	for _, item := range batch.items {
		if !c.process(ctx, item) {
			break
		}
		succeeded++
	}
	if succeeded == 0 {
		return false
	}
	highest := batch.items[succeeded-1].message
	if err := c.reader.CommitMessages(ctx, highest); err != nil {
		slog.Error("Kafka commit failed; successful effects will replay", "partition", key, "offset", highest.Offset, "error", err)
		return false
	}
	c.lastProgress.Store(time.Now().UnixMilli())
	batch.items = batch.items[succeeded:]
	if len(batch.items) == 0 {
		batch.firstQueued = time.Time{}
	} else {
		batch.firstQueued = time.Now()
	}
	slog.Info("partition batch committed", "partition", key, "messages", succeeded, "highest_offset", highest.Offset)
	return true
}

func (c *KafkaConsumer) Start(ctx context.Context, workerID string) {
	defer close(c.done)
	defer c.reader.Close()
	defer c.dlq.Close()
	slog.Info("Kafka partition-aware spatial consumer started", "worker_id", workerID, "batch_size", c.batchSize, "flush_interval", c.flushInterval)
	ticker := time.NewTicker(5 * time.Millisecond)
	defer ticker.Stop()
	batches := make(map[string]*partitionBatch)
	fetched := make(chan kafka.Message)
	fetchErrors := make(chan error, 1)
	go func() {
		defer close(fetched)
		for {
			message, err := c.reader.FetchMessage(ctx)
			if err != nil {
				if ctx.Err() == nil {
					fetchErrors <- err
				}
				return
			}
			select {
			case fetched <- message:
			case <-ctx.Done():
				return
			}
		}
	}()
	flushDue := func(flushCtx context.Context, force bool) {
		now := time.Now()
		for key, batch := range batches {
			if force || len(batch.items) >= c.batchSize || (!batch.firstQueued.IsZero() && now.Sub(batch.firstQueued) >= c.flushInterval) {
				c.flushPartition(flushCtx, key, batch)
				if len(batch.items) == 0 {
					delete(batches, key)
				}
			}
		}
	}
	shutdownFlush := func() {
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		flushDue(shutdownCtx, true)
		cancel()
		pending := 0
		for _, batch := range batches {
			pending += len(batch.items)
		}
		if pending > 0 {
			slog.Error("spatial consumer stopped with uncommitted messages", "pending", pending)
		} else {
			slog.Info("spatial consumer shutdown flush complete")
		}
	}
	for {
		select {
		case message, ok := <-fetched:
			if !ok {
				shutdownFlush()
				return
			}
			key := partitionID(message)
			batch := batches[key]
			if batch == nil {
				batch = &partitionBatch{firstQueued: time.Now()}
				batches[key] = batch
			}
			batch.items = append(batch.items, pendingTelemetry{message: message})
		case err := <-fetchErrors:
			slog.Error("Kafka fetch loop stopped", "error", err)
			shutdownFlush()
			return
		case <-ticker.C:
			c.lastProgress.Store(time.Now().UnixMilli())
			flushDue(ctx, false)
		case <-ctx.Done():
			shutdownFlush()
			return
		}
	}
}

func (c *KafkaConsumer) Ready(ctx context.Context) error {
	conn, err := kafka.DialContext(ctx, "tcp", c.brokerURL)
	if err != nil {
		return err
	}
	defer conn.Close()
	return c.projector.Ready(ctx)
}
func (c *KafkaConsumer) Healthy() bool {
	select {
	case <-c.done:
		return false
	default:
		return time.Since(time.UnixMilli(c.lastProgress.Load())) < 2*time.Minute
	}
}
func (c *KafkaConsumer) Wait(ctx context.Context) error {
	select {
	case <-c.done:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}
