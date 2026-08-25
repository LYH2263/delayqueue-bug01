package delayqueue

import "github.com/LYH2263/go-delayqueue/internal/clone"

func (b *Broker) SnapshotPending() []string {
	b.mu.Lock()
	defer b.mu.Unlock()
	return clone.Strings(b.pending)
}

func (b *Broker) SnapshotDead() []string {
	b.mu.Lock()
	defer b.mu.Unlock()
	return clone.Strings(b.dead)
}

func (b *Broker) ListPendingViews() []PendingView {
	b.mu.Lock()
	defer b.mu.Unlock()
	out := make([]PendingView, 0, len(b.pending))
	for _, id := range b.pending {
		t := b.tasks[id]
		out = append(out, PendingView{ID: id, Queue: t.Queue, RunAt: t.RunAt, Attempt: t.Attempts})
	}
	return out
}
