package delayqueue

import "time"

type Task struct {
	ID         string
	Payload    []byte
	Queue      string
	Attempts   int
	RunAt      time.Time
	LeasedBy   string
	LeaseUntil time.Time
	Dead       bool
}

type Stats struct {
	Pending   int
	Leased    int
	Dead      int
	Scheduled int
}

type PendingView struct {
	ID      string
	Queue   string
	RunAt   time.Time
	Attempt int
}
