package timer

import (
	"sync"
	"time"
)

type TimeTask interface {
	OnTime()
}

type task struct {
	tt     TimeTask
	expire int64
}

type Timer struct {
	m    sync.Mutex
	heap []task
	t    *time.Timer
}

func (t *Timer) siftup(i int) {
	h := t.heap
	v := h[i]

	for i > 0 {
		p := (i - 1) / 4

		if h[p].expire <= v.expire {
			break
		}

		h[i] = h[p]
		i = p
	}

	h[i] = v
}

func (t *Timer) siftdown(i int) {
	h := t.heap
	v := h[i]
	n := (len(h) + 2) / 4

	for i < n {
		b := 4*i + 1
		min := b

		for j := b + 1; j < b+4 && j < len(h); j++ {
			if h[j].expire < h[min].expire {
				min = j
			}
		}

		if h[min].expire >= v.expire {
			break
		}

		h[i] = h[min]
		i = min
	}

	h[i] = v
}

func (t *Timer) resetTimer(expire int64) {
	d := time.Duration(expire - time.Now().UnixNano())
	if t.t != nil {
		t.t.Reset(d)
	} else {
		t.t = time.NewTimer(d)
		go t.timeloop(t.t)
	}
}

func (t *Timer) removeExpired() {
	for {
		var fn TimeTask

		t.m.Lock()
		now := time.Now().UnixNano()
		if len(t.heap) > 0 && t.heap[0].expire <= now {
			// heap pop
			fn = t.heap[0].tt
			n := len(t.heap) - 1
			t.heap[0] = t.heap[n]
			t.heap = t.heap[:n]
			if n > 1 {
				t.siftdown(0)
			}
		}
		t.m.Unlock()

		if fn == nil {
			break
		}

		fn.OnTime()
	}
}

func (t *Timer) timeloop(tm *time.Timer) {
	for range tm.C {
		if t.t != tm {
			break
		}

		t.removeExpired()

		t.m.Lock()
		if len(t.heap) > 0 {
			t.resetTimer(t.heap[0].expire)
		}
		t.m.Unlock()
	}
	tm.Stop()
}

func (t *Timer) Add(tt TimeTask, expire int64) {
	t.m.Lock()
	// heap push
	t.heap = append(t.heap, task{tt, expire})
	t.siftup(len(t.heap) - 1)
	t.resetTimer(t.heap[0].expire)
	t.m.Unlock()
}

func (t *Timer) Stop() {
	t.m.Lock()
	tm := t.t
	t.heap = nil
	t.t = nil
	tm.Reset(0)
	t.m.Unlock()
}
