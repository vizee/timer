package timer

import (
	"sync"
	"testing"
	"time"
)

func TestTimer(t *testing.T) {
	var (
		tm Timer
		wg sync.WaitGroup
	)
	wg.Add(3)
	now := time.Now()
	t.Log(now)
	tm.Add(func() {
		t.Log("a", time.Now().Sub(now))
		wg.Done()
	}, now.Add(5*time.Second).UnixNano())
	tm.Add(func() {
		t.Log("b", time.Now().Sub(now))
		wg.Done()
	}, now.Add(3*time.Second).UnixNano())
	tm.Add(func() {
		t.Log("c", time.Now().Sub(now))
		wg.Done()
	}, now.Add(5*time.Second).UnixNano()+1000)
	wg.Wait()
	tm.Stop()
}
