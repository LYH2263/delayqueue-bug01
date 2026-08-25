package delayqueue

import (
	"time"

	"github.com/LYH2263/go-delayqueue/internal/clock"
)

type Options struct {
	JournalPath  string
	LeaseTimeout time.Duration
	MaxAttempts  int
	Clock        clock.Clock
	PollInterval time.Duration
}

func (o Options) withDefaults() Options {
	if o.LeaseTimeout <= 0 {
		o.LeaseTimeout = 30 * time.Second
	}
	if o.MaxAttempts <= 0 {
		o.MaxAttempts = 5
	}
	if o.Clock == nil {
		o.Clock = clock.System{}
	}
	if o.PollInterval <= 0 {
		o.PollInterval = 10 * time.Millisecond
	}
	return o
}
