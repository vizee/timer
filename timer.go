package timer

import (
	"sync"
	"time"
)

type Handler interface {
	OnTime()
}

type Task struct {
	expire int64
	hidx   int
	ontime Handler
}

func heapsiftup(h []*Task, i int) {
	v := h[i]
	ex := v.expire
	for i > 0 {
		p := (i - 1) / 2
		if h[p].expire <= ex {
			break
		}
		h[i] = h[p]
		h[i].hidx = i
		i = p
	}
	h[i] = v
	v.hidx = i
}

func heapsiftdown(h []*Task, i int) {
	v := h[i]
	ex := v.expire
	n := len(h) / 2
	for i < n {
		c := 2*i + 1
		r := c + 1
		if r < len(h) && h[r].expire < h[c].expire {
			c = r
		}
		if h[c].expire >= ex {
			break
		}
		h[i] = h[c]
		h[i].hidx = i
		i = c
	}
	h[i] = v
	v.hidx = i
}

type Timer struct {
	mu    sync.Mutex
	heap  []*Task
	timer *time.Timer
}

func (t *Timer) remove(idx int) *Task {
	h := t.heap
	task := h[idx]
	last := len(h) - 1
	h[idx] = h[last]
	h[idx].hidx = idx
	if idx != last {
		heapsiftdown(h[:last], idx)
		heapsiftup(h[:last], idx)
	}
	t.heap = h[:last]
	return task
}

func (t *Timer) popExpired() {
	for {
		var task *Task
		t.mu.Lock()
		if len(t.heap) > 0 {
			if t.heap[0].expire <= time.Now().UnixNano() {
				task = t.remove(0)
				task.hidx = -1
			}
		}
		t.mu.Unlock()
		if task == nil {
			break
		}
		task.ontime.OnTime()
	}
}

func (t *Timer) timeloop() {
	for range t.timer.C {
		t.popExpired()
		t.mu.Lock()
		if len(t.heap) > 0 {
			t.resetTimer(t.heap[0].expire)
		}
		t.mu.Unlock()
	}
}

func (t *Timer) resetTimer(expire int64) {
	// TODO reset timer only if expire changed.
	d := time.Duration(expire - time.Now().UnixNano())
	if t.timer == nil {
		t.timer = time.NewTimer(d)
		go t.timeloop()
	} else {
		t.timer.Reset(d)
	}
}

func (t *Timer) Remove(task *Task) {
	t.mu.Lock()
	if task.hidx >= 0 {
		t.remove(task.hidx)
		task.hidx = -1
		if len(t.heap) > 0 {
			t.resetTimer(t.heap[0].expire)
		}
	}
	t.mu.Unlock()
}

func (t *Timer) Reset(task *Task, expire int64) {
	// assert(t.idx < len(tm.h) && tm.h[t.idx] == t)
	t.mu.Lock()
	task.expire = expire
	if task.hidx < 0 {
		task.hidx = len(t.heap)
		t.heap = append(t.heap, task)
	} else {
		heapsiftdown(t.heap, task.hidx)
	}
	heapsiftup(t.heap, task.hidx)
	t.resetTimer(t.heap[0].expire)
	t.mu.Unlock()
}

func (t *Timer) Add(h Handler, expire int64) *Task {
	task := &Task{
		ontime: h,
		hidx:   -1,
	}
	t.Reset(task, expire)
	return task
}
