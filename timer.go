package timer

import (
	"sync"
	"time"
)

type TaskFunc func()

type task struct {
	fn     TaskFunc
	expire int64
}

type Timer struct {
	m      sync.Mutex
	heap   []task
	t      *time.Timer
	expire int64
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
	n := (len(h) + 3) / 4
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
	if t.expire == expire {
		return
	}
	t.expire = expire
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
		var fn TaskFunc
		t.m.Lock()
		if len(t.heap) > 0 && t.heap[0].expire <= time.Now().UnixNano() {
			fn = t.heap[0].fn
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
		fn()
	}
}

func (t *Timer) timeloop(tt *time.Timer) {
	for range tt.C {
		if t.t != tt {
			break
		}
		t.removeExpired()
		t.m.Lock()
		if len(t.heap) > 0 {
			t.resetTimer(t.heap[0].expire)
		}
		t.m.Unlock()
	}
	tt.Stop()
}

func (t *Timer) Add(f TaskFunc, expire int64) {
	t.m.Lock()
	defer t.m.Unlock()
	t.heap = append(t.heap, task{f, expire})
	t.siftup(len(t.heap) - 1)
	t.resetTimer(t.heap[0].expire)
}

func (t *Timer) Stop() {
	t.m.Lock()
	defer t.m.Unlock()
	tt := t.t
	t.t = nil
	t.expire = 0
	tt.Reset(0)
}
