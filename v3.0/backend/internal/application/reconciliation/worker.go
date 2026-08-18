package reconciliation

import (
	"context"
	"errors"
	"log/slog"
	"time"

	"github.com/Akashpg-M/polaris/backend/internal/adapter/repository"
	"github.com/Akashpg-M/polaris/backend/internal/application/orchestration"
)

type Worker struct {
	store      *repository.RegistryStore
	service    *orchestration.Service
	owners     *repository.ConnectionOwnershipStore
	interval   time.Duration
	ackTimeout time.Duration
}

func New(store *repository.RegistryStore, service *orchestration.Service, owners *repository.ConnectionOwnershipStore, interval, ackTimeout time.Duration) *Worker {
	if interval < 100*time.Millisecond {
		interval = time.Second
	}
	return &Worker{store: store, service: service, owners: owners, interval: interval, ackTimeout: ackTimeout}
}

func (w *Worker) Start(ctx context.Context) {
	ticker := time.NewTicker(w.interval)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			w.run(ctx)
		}
	}
}

func (w *Worker) run(ctx context.Context) {
	if err := w.store.ReconcileCommands(ctx, w.ackTimeout); err != nil {
		slog.Error("command reconciliation failed", "error", err)
	}
	if err := w.store.FailExpiredPendingTasks(ctx); err != nil {
		slog.Error("task expiry reconciliation failed", "error", err)
	}
	if err := w.owners.CleanExpired(ctx); err != nil {
		slog.Error("connection lease reconciliation failed", "error", err)
	}
	pending, err := w.store.PendingTasks(ctx, 50)
	if err != nil {
		slog.Error("pending task scan failed", "error", err)
		return
	}
	for _, v := range pending {
		_, err = w.service.Assign(ctx, v, "reconciler", "")
		if err != nil && !errors.Is(err, orchestration.ErrNoEligibleDevice) && !errors.Is(err, repository.ErrConflict) && !errors.Is(err, repository.ErrInvalidTransition) {
			slog.Error("pending task assignment failed", "task_id", v.TaskID, "error", err)
		}
	}
}
