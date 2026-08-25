package delayqueue

import (
	"context"
	"time"
)

func (b *Broker) Dequeue(ctx context.Context, workerID string) (*Task, error) {
	if workerID == "" {
		return nil, ErrInvalid
	}
	for {
		if err := ctx.Err(); err != nil {
			return nil, err
		}
		b.mu.Lock()
		if b.closed {
			b.mu.Unlock()
			return nil, ErrClosed
		}
		if len(b.pending) == 0 {
			b.mu.Unlock()
			select {
			case <-ctx.Done():
				return nil, ctx.Err()
			case <-time.After(b.opts.PollInterval):
				continue
			}
		}
		id := b.pending[0]
		b.pending = b.pending[1:]
		t := b.tasks[id]
		until := b.opts.Clock.Now().Add(b.opts.LeaseTimeout)
		t.LeasedBy = workerID
		t.LeaseUntil = until
		b.leases.Grant(id, workerID, until)
		out := *t
		out.Payload = append([]byte(nil), t.Payload...)
		b.mu.Unlock()
		return &out, nil
	}
}
