package timer

import (
	"log"
	"time"
)

type Timer struct {
	note  string
	accum time.Duration
	start time.Time
}

func New(note string) *Timer {
	return &Timer{
		note: note,
	}
}

func (t *Timer) Start() {
	t.start = time.Now()
}

func (t *Timer) Stop() {
	t.accum += time.Since(t.start)
}

func (t *Timer) Print() {
	log.Printf("%s: %s", t.note, t.accum)
}
