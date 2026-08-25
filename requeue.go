package delayqueue

func (b *Broker) RequeueDead(id string) error {
	b.mu.Lock()
	defer b.mu.Unlock()
	if b.closed {
		return ErrClosed
	}
	t, ok := b.tasks[id]
	if !ok || !t.Dead {
		return ErrNotFound
	}
	t.Dead = false
	t.Attempts = 0
	t.LeasedBy = ""
	b.pending = append(b.pending, id)
	for i, x := range b.dead {
		if x == id {
			b.dead = append(b.dead[:i], b.dead[i+1:]...)
			break
		}
	}
	if b.journal != nil {
		_ = b.journal.AppendRequeue(id)
	}
	return nil
}
