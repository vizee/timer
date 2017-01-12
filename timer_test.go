package timer

import (
	"log"
	"sync"
	"testing"
	"time"
)

type testTimer struct {
	name string
	now  time.Time
	wg   *sync.WaitGroup
}

func (t *testTimer) OnTime() {
	log.Println(t.name, time.Now().Sub(t.now))
	t.wg.Done()
}

func TestTimer(t *testing.T) {
	var (
		tm Timer
		wg sync.WaitGroup
	)
	wg.Add(4)
	now := time.Now()
	log.Println(now)
	tm.Add(&testTimer{"a", now, &wg}, now.Add(5*time.Second).UnixNano())
	log.Println("add a 5s")

	tm.Add(&testTimer{"b", now, &wg}, now.Add(3*time.Second).UnixNano())
	log.Println("add b 3s")

	tm.Add(&testTimer{"c", now, &wg}, now.Add(5*time.Second).UnixNano()+1000000)
	log.Println("add c 5.000001s")

	tm.Add(&testTimer{"d", now, &wg}, now.Add(10*time.Second).UnixNano())
	log.Println("add d 10s")

	time.AfterFunc(5*time.Second, func() {
		log.Println("stop d", tm.heap)
		wg.Done()
		tm.Stop()
	})

	wg.Wait()
	log.Println("try wait stoped task")
	wg.Add(1)

	tm.Add(&testTimer{"expired", now, &wg}, time.Now().Add(6*time.Second).UnixNano())
	wg.Wait()
	log.Println("task heap", tm.heap)
}
