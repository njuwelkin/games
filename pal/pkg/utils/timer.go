package utils

import "container/heap"

type timerNode struct {
	freq     int
	deadline int
	repeat   int
	callback func(int)
	until    func() bool

	id uint64
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

	idCounter uint64
}

func NewTimer() *TimerManager {
	return &TimerManager{
		queue: timerHeap([]timerNode{}),
	}
}

func (t *TimerManager) AddOneTimeEvent(countDown int, fn func(int)) uint64 {
	t.idCounter++
	heap.Push(&t.queue, timerNode{
		deadline: t.count + countDown,
		repeat:   1,
		callback: fn,
		id:       t.idCounter,
	})
	return t.idCounter
}

func (t *TimerManager) AddRepeatEvent(freq int, repeat int, fn func(int)) uint64 {
	t.idCounter++
	heap.Push(&t.queue, timerNode{
		deadline: t.count + freq,
		freq:     freq,
		repeat:   repeat,
		callback: fn,
		id:       t.idCounter,
	})
	return t.idCounter
}

func (t *TimerManager) RepeatUntil(freq int, fn func(int), until func() bool) uint64 {
	t.idCounter++
	heap.Push(&t.queue, timerNode{
		deadline: t.count + freq,
		freq:     freq,
		until:    until,
		callback: fn,
		id:       t.idCounter,
	})
	return t.idCounter
}

func (t *TimerManager) RemoveEvent(id uint64) {
	var i int
	for i = range t.queue {
		if t.queue[i].id == id {
			break
		}
	}
	heap.Remove(&t.queue, i)
}

func (t *TimerManager) Update() {
	for t.queue.Len() > 0 && t.queue[0].deadline <= t.count {
		head := t.queue[0]
		head.repeat--
		head.callback(head.repeat)
		heap.Pop(&t.queue)

		if head.until != nil && !head.until() ||
			head.repeat > 0 {
			head.deadline += head.freq
			heap.Push(&t.queue, head)
		}
	}
	t.count++
}
