package delayqueue

import "time"

func (b *Broker) ExtendLease(workerID, id string, extra time.Duration) error {
	b.mu.Lock()
	defer b.mu.Unlock()
	if b.closed {
		return ErrClosed
	}
	t, ok := b.tasks[id]
	if !ok {
		return ErrNotFound
	}
	if t.LeasedBy != workerID {
		return ErrLeaseLost
	}
	base := b.opts.Clock.Now()
	t.LeaseUntil = base.Add(extra)
	b.leases.Grant(id, workerID, t.LeaseUntil)
	return nil
}
