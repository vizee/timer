package timer

import (
	"log"
	"testing"
	"time"
)

type testtask string

func (t testtask) OnTime() {
	log.Println("> task", string(t))
}

func TestTimerAdd(t *testing.T) {
	tm := Timer{}
	log.Println("start")
	now := time.Now()
	tm.Add(testtask("1s"), now.Add(time.Second).UnixNano())
	tm.Add(testtask("3s"), now.Add(time.Second*3).UnixNano())
	tm.Add(testtask("2s"), now.Add(time.Second*2).UnixNano())
	time.Sleep(time.Second * 5)
}

func TestTimerCancel(t *testing.T) {
	tm := Timer{}
	log.Println("start")

	now := time.Now()
	tm.Add(testtask("1s"), now.Add(time.Second).UnixNano())
	t3s := tm.Add(testtask("3s"), now.Add(time.Second*3).UnixNano())
	tm.Add(testtask("2s"), now.Add(time.Second*2).UnixNano())
	log.Println("wait 1s")
	time.Sleep(time.Second)
	tm.Remove(t3s)
	log.Println("cancel 3s")
	time.Sleep(time.Second * 5)
}

type testtask2 struct {
	name string
	tm   *Timer
	t    *Task
}

func (t *testtask2) OnTime() {
	log.Println("> task", t.name)
	if t.name == "1s" && t.tm != nil {
		log.Println("> reset 1s next 4s")
		t.tm.Reset(t.t, time.Now().Add(4*time.Second).UnixNano())
		t.tm = nil
	}
}

func TestTimerReset(t *testing.T) {
	tm := &Timer{}
	log.Println("start")
	now := time.Now()
	t1s := &testtask2{
		name: "1s",
		tm:   tm,
	}
	t1s.t = tm.Add(t1s, now.Add(time.Second).UnixNano())
	tm.Add(testtask("3s"), now.Add(time.Second*3).UnixNano())
	tm.Add(testtask("2s"), now.Add(time.Second*2).UnixNano())
	time.Sleep(time.Second * 10)
}

func init() {
	log.SetFlags(log.Ltime | log.Lmicroseconds)
}
