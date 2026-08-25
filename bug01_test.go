package delayqueue

import (
	"bytes"
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/LYH2263/go-delayqueue/internal/clock"
)

func TestBug01_ScheduleJournalRunAtCrossLayer(t *testing.T) {
	path := filepath.Join(t.TempDir(), "j.log")
	now := time.Unix(1_700_000_000, 0)
	b, err := New(Options{JournalPath: path, Clock: clock.Fixed{T: now}})
	if err != nil {
		t.Fatal(err)
	}
	runAt := now.Add(time.Hour)
	if err := b.ScheduleAt(context.Background(), "later", "q", []byte("x"), runAt); err != nil {
		t.Fatal(err)
	}
	if err := b.journal.Flush(); err != nil {
		t.Fatal(err)
	}
	raw, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	_ = b.Close()
	var rec struct {
		Op    string `json:"op"`
		RunAt int64  `json:"run_at"`
	}
	line := bytes.TrimSpace(raw)
	if err := json.Unmarshal(line, &rec); err != nil {
		t.Fatal(err)
	}
	if rec.Op != "schedule" {
		t.Fatalf("op=%s", rec.Op)
	}
	if rec.RunAt != runAt.UnixNano() {
		t.Fatalf("journal run_at=%d want %d (cross-layer drift)", rec.RunAt, runAt.UnixNano())
	}
}
