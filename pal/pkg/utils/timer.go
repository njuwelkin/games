package utils

import (
	"container/heap"
	"sync"
)

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
	mu        sync.Mutex
}

func NewTimer() *TimerManager {
	return &TimerManager{
		queue: timerHeap([]timerNode{}),
	}
}

func (t *TimerManager) AddOneTimeEvent(countDown int, fn func(int)) uint64 {
	t.mu.Lock()
	defer t.mu.Unlock()

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
	t.mu.Lock()
	defer t.mu.Unlock()

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
	t.mu.Lock()
	defer t.mu.Unlock()

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
	t.mu.Lock()
	defer t.mu.Unlock()

	var i int
	for i = range t.queue {
		if t.queue[i].id == id {
			break
		}
	}
	if len(t.queue) > 0 {
		heap.Remove(&t.queue, i)
	}
}

func (t *TimerManager) Update() {
	t.mu.Lock()
	t.count++

	var tasksToExecute []timerNode
	for t.queue.Len() > 0 && t.queue[0].deadline <= t.count {
		head := heap.Pop(&t.queue).(timerNode)
		head.repeat--
		tasksToExecute = append(tasksToExecute, head)
	}
	t.mu.Unlock()

	for _, task := range tasksToExecute {
		task.callback(task.repeat)

		t.mu.Lock()
		if task.until != nil && !task.until() || task.repeat > 0 {
			task.deadline += task.freq
			heap.Push(&t.queue, task)
		}
		t.mu.Unlock()
	}
}
