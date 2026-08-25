package delayqueue

import "errors"

var (
	ErrClosed    = errors.New("delayqueue: broker closed")
	ErrNotFound  = errors.New("delayqueue: task not found")
	ErrInvalid   = errors.New("delayqueue: invalid argument")
	ErrLeaseLost = errors.New("delayqueue: lease expired or lost")
	ErrNoClock   = errors.New("delayqueue: clock not configured")
	ErrJournal   = errors.New("delayqueue: journal failure")
)
