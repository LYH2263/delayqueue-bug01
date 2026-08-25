package schedule

import (
	"container/heap"
	"sync"
	"time"

	"github.com/LYH2263/go-delayqueue/internal/clock"
)

type item struct {
	id    string
	runAt time.Time
	index int
}

type innerHeap []*item

func (h innerHeap) Len() int           { return len(h) }
func (h innerHeap) Less(i, j int) bool { return h[i].runAt.Before(h[j].runAt) }
func (h innerHeap) Swap(i, j int) {
	h[i], h[j] = h[j], h[i]
	h[i].index = i
	h[j].index = j
}
func (h *innerHeap) Push(x any) {
	n := len(*h)
	it := x.(*item)
	it.index = n
	*h = append(*h, it)
}
func (h *innerHeap) Pop() any {
	old := *h
	n := len(old)
	it := old[n-1]
	*h = old[:n-1]
	return it
}

type Heap struct {
	mu  sync.Mutex
	clk clock.Clock
	h   innerHeap
	idx map[string]int
}

func NewHeap(clk clock.Clock) *Heap {
	h := &Heap{clk: clk, idx: make(map[string]int)}
	heap.Init(&h.h)
	return h
}

func (h *Heap) Len() int {
	h.mu.Lock()
	defer h.mu.Unlock()
	return len(h.h)
}

func (h *Heap) Push(id string, runAt time.Time) {
	h.mu.Lock()
	defer h.mu.Unlock()
	if i, ok := h.idx[id]; ok {
		h.h[i].runAt = runAt
		heap.Fix(&h.h, i)
		return
	}
	heap.Push(&h.h, &item{id: id, runAt: runAt})
	h.idx[id] = len(h.h) - 1
}

func (h *Heap) Remove(id string) {
	h.mu.Lock()
	defer h.mu.Unlock()
	i, ok := h.idx[id]
	if !ok {
		return
	}
	heap.Remove(&h.h, i)
	delete(h.idx, id)
}

func (h *Heap) PopDue(now time.Time) (string, bool) {
	h.mu.Lock()
	defer h.mu.Unlock()
	if len(h.h) == 0 {
		return "", false
	}
	if h.h[0].runAt.After(now) {
		return "", false
	}
	it := heap.Pop(&h.h).(*item)
	delete(h.idx, it.id)
	return it.id, true
}
