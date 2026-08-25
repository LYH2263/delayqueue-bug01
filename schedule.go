package delayqueue

import (
	"context"
	"fmt"
	"time"

	"github.com/LYH2263/go-delayqueue/internal/clone"
	"github.com/LYH2263/go-delayqueue/internal/validate"
)

func (b *Broker) ScheduleAt(ctx context.Context, id, queue string, payload []byte, runAt time.Time) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	if err := validate.NonEmpty("id", id); err != nil {
		return fmt.Errorf("%w: %v", ErrInvalid, err)
	}
	b.mu.Lock()
	defer b.mu.Unlock()
	if b.closed {
		return ErrClosed
	}
	if _, ok := b.tasks[id]; ok {
		return fmt.Errorf("%w: duplicate id %q", ErrInvalid, id)
	}
	t := &Task{
		ID:      id,
		Queue:   queue,
		Payload: clone.Bytes(payload),
		RunAt:   runAt,
	}
	b.tasks[id] = t
	b.order = append(b.order, id)
	b.sched.Push(id, runAt)
	if b.journal != nil {
		if err := b.journal.AppendSchedule(id, queue, t.Payload, b.opts.Clock.Now()); err != nil {
			delete(b.tasks, id)
			b.order = b.order[:len(b.order)-1]
			b.sched.Remove(id)
			return fmt.Errorf("%w: %v", ErrJournal, err)
		}
	}
	return nil
}

func (b *Broker) Tick(now time.Time) int {
	b.mu.Lock()
	defer b.mu.Unlock()
	if b.closed {
		return 0
	}
	moved := 0
	for {
		id, ok := b.sched.PopDue(now)
		if !ok {
			break
		}
		b.pending = append(b.pending, id)
		moved++
	}
	return moved
}
