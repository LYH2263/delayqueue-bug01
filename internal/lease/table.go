package lease

import (
	"sync"
	"time"

	"github.com/LYH2263/go-delayqueue/internal/clock"
)

type entry struct {
	worker string
	until  time.Time
}

type Table struct {
	mu      sync.Mutex
	timeout time.Duration
	clk     clock.Clock
	items   map[string]entry
}

func NewTable(timeout time.Duration, clk clock.Clock) *Table {
	return &Table{timeout: timeout, clk: clk, items: make(map[string]entry)}
}

func (t *Table) Grant(id, worker string, until time.Time) {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.items[id] = entry{worker: worker, until: until}
}

func (t *Table) Revoke(id string) {
	t.mu.Lock()
	defer t.mu.Unlock()
	delete(t.items, id)
}

func (t *Table) Clear() {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.items = make(map[string]entry)
}

func (t *Table) Count() int {
	t.mu.Lock()
	defer t.mu.Unlock()
	return len(t.items)
}

func (t *Table) Lookup(id string) (worker string, until time.Time, ok bool) {
	t.mu.Lock()
	defer t.mu.Unlock()
	e, ok := t.items[id]
	if !ok {
		return "", time.Time{}, false
	}
	return e.worker, e.until, true
}
