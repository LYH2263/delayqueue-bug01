package delayqueue

func (b *Broker) Close() error {
	b.mu.Lock()
	defer b.mu.Unlock()
	if b.closed {
		return nil
	}
	if b.journal != nil {
		_ = b.journal.Flush()
	}
	b.leases.Clear()
	b.closed = true
	if b.journal != nil {
		return b.journal.Close()
	}
	return nil
}

func (b *Broker) CloseFlushCount() int {
	b.mu.Lock()
	defer b.mu.Unlock()
	if b.closed {
		return 0
	}
	n := 0
	if b.journal != nil {
		n = b.journal.PendingRecords()
		_ = b.journal.Flush()
	}
	b.leases.Clear()
	b.closed = true
	if b.journal != nil {
		_ = b.journal.Close()
	}
	return n
}
