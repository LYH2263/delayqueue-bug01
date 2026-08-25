package delayqueue

import (
	"sync"

	"github.com/LYH2263/go-delayqueue/internal/journal"
	"github.com/LYH2263/go-delayqueue/internal/lease"
	"github.com/LYH2263/go-delayqueue/internal/retry"
	"github.com/LYH2263/go-delayqueue/internal/schedule"
)

type Broker struct {
	mu      sync.Mutex
	closed  bool
	opts    Options
	tasks   map[string]*Task
	order   []string
	dead    []string
	pending []string
	journal *journal.Journal
	leases  *lease.Table
	sched   *schedule.Heap
	retry   retry.Policy
}

func New(opts Options) (*Broker, error) {
	opts = opts.withDefaults()
	b := &Broker{
		opts:   opts,
		tasks:  make(map[string]*Task),
		leases: lease.NewTable(opts.LeaseTimeout, opts.Clock),
		sched:  schedule.NewHeap(opts.Clock),
		retry:  retry.DefaultPolicy(),
	}
	if opts.JournalPath != "" {
		j, err := journal.Open(opts.JournalPath)
		if err != nil {
			return nil, err
		}
		b.journal = j
	}
	return b, nil
}

func (b *Broker) Stats() Stats {
	b.mu.Lock()
	defer b.mu.Unlock()
	leased := 0
	for _, t := range b.tasks {
		if t.LeasedBy != "" {
			leased++
		}
	}
	return Stats{
		Pending:   len(b.pending),
		Leased:    leased,
		Dead:      len(b.dead),
		Scheduled: b.sched.Len(),
	}
}
