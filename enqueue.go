package delayqueue

import (
	"context"
	"fmt"

	"github.com/LYH2263/go-delayqueue/internal/clone"
	"github.com/LYH2263/go-delayqueue/internal/validate"
)

func (b *Broker) Enqueue(ctx context.Context, id, queue string, payload []byte) error {
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
		RunAt:   b.opts.Clock.Now(),
	}
	b.tasks[id] = t
	b.order = append(b.order, id)
	b.pending = append(b.pending, id)
	if b.journal != nil {
		if err := b.journal.AppendEnqueue(id, queue, t.Payload); err != nil {
			delete(b.tasks, id)
			b.order = b.order[:len(b.order)-1]
			b.pending = b.pending[:len(b.pending)-1]
			return fmt.Errorf("%w: %v", ErrJournal, err)
		}
	}
	return nil
}
