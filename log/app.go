package log

import (
	"context"
	"github.com/rs/zerolog"
	"time"
)

type appLogger struct {
	*Logger
	ticker *time.Ticker
}

// Initialize initializes the appLogger
func (l appLogger) Initialize(ctx context.Context) {
	go func(ctx context.Context) {
		tlap := time.Now().UTC()
		for {
			select {
			case t := <-l.ticker.C:
				l.Heartbeat(t.Sub(tlap))
				tlap = t
			case <-ctx.Done():
				l.ticker.Stop()
				return
			}
		}
	}(ctx)

	l.Log("app.startup").
		Str("level", zerolog.GlobalLevel().String()).Msg("")
}

// Heartbeat for the appLogger
func (l appLogger) Heartbeat(dur time.Duration) {
	e := l.Log("app.heartbeat")
	if dur > 0 {
		e = e.Float64("interval", dur.Truncate(time.Millisecond/10).Seconds())
	}
	e.Msg("")
}
