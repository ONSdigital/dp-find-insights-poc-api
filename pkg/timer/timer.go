package timer

import (
	"log"
	"time"
)

type Timer struct {
	note  string
	start time.Time
}

func New(note string) *Timer {
	return &Timer{
		note:  note,
		start: time.Now(),
	}
}

func (t *Timer) Stop() {
	log.Printf("%s: %s\n", t.note, time.Since(t.start))
}
