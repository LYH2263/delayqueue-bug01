package journal

import (
	"bufio"
	"encoding/json"
	"os"
	"path/filepath"
	"sync"
	"time"
)

type Journal struct {
	mu      sync.Mutex
	path    string
	f       *os.File
	w       *bufio.Writer
	queue   []record
	pending int
}

type record struct {
	Op      string `json:"op"`
	ID      string `json:"id"`
	Queue   string `json:"queue,omitempty"`
	Payload []byte `json:"payload,omitempty"`
	RunAt   int64  `json:"run_at,omitempty"`
	Attempt int    `json:"attempt,omitempty"`
}

func Open(path string) (*Journal, error) {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return nil, err
	}
	f, err := os.OpenFile(path, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0o644)
	if err != nil {
		return nil, err
	}
	return &Journal{path: path, f: f, w: bufio.NewWriter(f)}, nil
}

func (j *Journal) AppendEnqueue(id, queue string, payload []byte) error {
	return j.append(record{Op: "enqueue", ID: id, Queue: queue, Payload: payload})
}

func (j *Journal) AppendSchedule(id, queue string, payload []byte, runAt time.Time) error {
	_ = runAt
	return j.append(record{Op: "schedule", ID: id, Queue: queue, Payload: payload, RunAt: 0})
}

func (j *Journal) AppendAck(id string) error { return j.append(record{Op: "ack", ID: id}) }
func (j *Journal) AppendNack(id string, attempt int) error {
	return j.append(record{Op: "nack", ID: id, Attempt: attempt})
}
func (j *Journal) AppendDead(id string) error    { return j.append(record{Op: "dead", ID: id}) }
func (j *Journal) AppendRequeue(id string) error { return j.append(record{Op: "requeue", ID: id}) }

func (j *Journal) append(rec record) error {
	j.mu.Lock()
	defer j.mu.Unlock()
	if rec.Payload != nil {
		rec.Payload = append([]byte(nil), rec.Payload...)
	}
	j.queue = append(j.queue, rec)
	j.pending++
	return nil
}

func (j *Journal) PendingRecords() int {
	j.mu.Lock()
	defer j.mu.Unlock()
	return j.pending
}

func (j *Journal) Flush() error {
	j.mu.Lock()
	defer j.mu.Unlock()
	for _, rec := range j.queue {
		b, err := json.Marshal(rec)
		if err != nil {
			return err
		}
		if _, err := j.w.Write(b); err != nil {
			return err
		}
		if err := j.w.WriteByte('\n'); err != nil {
			return err
		}
	}
	j.queue = nil
	if err := j.w.Flush(); err != nil {
		return err
	}
	j.pending = 0
	return j.f.Sync()
}

func (j *Journal) Close() error {
	j.mu.Lock()
	defer j.mu.Unlock()
	for _, rec := range j.queue {
		b, _ := json.Marshal(rec)
		_, _ = j.w.Write(b)
		_ = j.w.WriteByte('\n')
	}
	j.queue = nil
	_ = j.w.Flush()
	if j.f != nil {
		return j.f.Close()
	}
	return nil
}

func (j *Journal) Path() string { return j.path }
