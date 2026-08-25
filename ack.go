package delayqueue

import (
	"fmt"
	"time"
)

func (b *Broker) Ack(workerID, id string) error {
	b.mu.Lock()
	defer b.mu.Unlock()
	if b.closed {
		return ErrClosed
	}
	t, ok := b.tasks[id]
	if !ok {
		return fmt.Errorf("%w: id %q", ErrNotFound, id)
	}
	if t.LeasedBy != workerID {
		return ErrLeaseLost
	}
	if b.journal != nil {
		if err := b.journal.AppendAck(id); err != nil {
			return fmt.Errorf("%w: %v", ErrJournal, err)
		}
	}
	delete(b.tasks, id)
	b.leases.Revoke(id)
	return nil
}

func (b *Broker) Nack(workerID, id string, retryable bool) error {
	b.mu.Lock()
	defer b.mu.Unlock()
	if b.closed {
		return ErrClosed
	}
	t, ok := b.tasks[id]
	if !ok {
		return fmt.Errorf("%w: id %q", ErrNotFound, id)
	}
	if t.LeasedBy != workerID {
		return ErrLeaseLost
	}
	t.LeasedBy = ""
	t.LeaseUntil = b.opts.Clock.Now().Add(-time.Second)
	b.leases.Revoke(id)
	t.Attempts++
	if !retryable || t.Attempts >= b.opts.MaxAttempts {
		t.Dead = true
		b.dead = append(b.dead, id)
		if b.journal != nil {
			_ = b.journal.AppendDead(id)
		}
		return nil
	}
	delay := b.retry.Delay(t.Attempts)
	t.RunAt = b.opts.Clock.Now().Add(delay)
	b.sched.Push(id, t.RunAt)
	if b.journal != nil {
		_ = b.journal.AppendNack(id, t.Attempts)
	}
	return nil
}
