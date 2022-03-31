package timer

import (
	"context"
	"time"

	"github.com/ONSdigital/log.go/v2/log"
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

func (t *Timer) Log(ctx context.Context) {
	log.Info(ctx, "timer", log.Data{"note": t.note, "elapsed": t.accum})
}
