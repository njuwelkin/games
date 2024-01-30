package tank

import (
	"container/heap"
)

func Max(x, y int) int {
	if x > y {
		return x
	}
	return y
}

func Min(x, y int) int {
	if x < y {
		return x
	}
	return y
}

func Abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func LowerBound(x, divisor int) int {
	return x - x%divisor
}

func UpperBound(x, divisor int) int {
	return LowerBound(x+divisor-1, divisor)
}

type timerNode struct {
	deadline int
	callback func()
}

type timerHeap []timerNode

func (h timerHeap) Swap(i, j int) {
	h[i], h[j] = h[j], h[i]
}

func (h timerHeap) Len() int {
	return len(h)
}

func (h timerHeap) Less(i, j int) bool {
	return h[i].deadline < h[j].deadline
}

func (h *timerHeap) Push(node interface{}) {
	*h = append(*h, node.(timerNode))
}

func (h *timerHeap) Pop() interface{} {
	tmp := (*h)[len(*h)-1]
	*h = (*h)[:len(*h)-1]
	return tmp
}

type TimerManager struct {
	queue timerHeap
	count int
}

func NewTimer() *TimerManager {
	return &TimerManager{
		queue: timerHeap([]timerNode{}),
	}
}

func (t *TimerManager) AddEvent(countDown int, fn func()) {
	heap.Push(&t.queue, timerNode{
		deadline: t.count + countDown,
		callback: fn,
	})
}

func (t *TimerManager) Update() {
	for t.queue.Len() > 0 && t.queue[0].deadline <= t.count {
		t.queue[0].callback()
		heap.Pop(&t.queue)
	}
	t.count++
}
