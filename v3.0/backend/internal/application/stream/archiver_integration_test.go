//go:build integration

package stream

import (
	"context"
	"os"
	"testing"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func TestArchiveReplayIsIdempotent(t *testing.T) {
	postgresURL := os.Getenv("POSTGRES_URL")
	if postgresURL == "" {
		postgresURL = "postgres://polaris_user:polaris_password@localhost:5432/polaris_core?sslmode=disable"
	}
	db, err := sqlx.Connect("postgres", postgresURL)
	if err != nil {
		t.Skipf("PostgreSQL integration dependency unavailable: %v", err)
	}
	defer db.Close()
	e := streamEnvelope("archive-replay-integration", 1)
	a := &KafkaPostgresArchiver{db: db}
	ctx := context.Background()
	defer db.ExecContext(ctx, "DELETE FROM telemetry_history WHERE event_id=$1", e.EventID)
	if err := a.archive(ctx, e); err != nil {
		t.Fatal(err)
	}
	if err := a.archive(ctx, e); err != nil {
		t.Fatal(err)
	}
	var count int
	if err := db.GetContext(ctx, &count, "SELECT count(*) FROM telemetry_history WHERE event_id=$1", e.EventID); err != nil {
		t.Fatal(err)
	}
	if count != 1 {
		t.Fatalf("duplicate replay produced %d rows", count)
	}
}
