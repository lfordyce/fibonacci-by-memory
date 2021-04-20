package log

import (
	"bufio"
	"bytes"
	"context"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/pkgerrors"
	"io"
	"log"
	"os"
	"time"
)

type LoggerOption interface {
	Apply(*Logger)
}

type LoggerOptionFunc func(*Logger)

func (fn LoggerOptionFunc) Apply(logger *Logger) {
	fn(logger)
}

type Logger struct {
	*zerolog.Logger
}

func WithApp() LoggerOption {
	return LoggerOptionFunc(func(l *Logger) {
		al := appLogger{
			Logger: l,
			ticker: time.NewTicker(30 * time.Second),
		}
		al.Initialize(context.Background())
	})
}

// NewLogger returns a new Logger instance
func NewLogger(opts ...LoggerOption) Logger {
	zerolog.TimestampFunc = func() time.Time { return time.Now().UTC() }
	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack

	zl := zerolog.New(os.Stdout)
	l := Logger{&zl}

	for _, opt := range opts {
		opt.Apply(&l)
	}
	return l
}

func (l Logger) Log(iom string) *zerolog.Event {
	return l.Logger.Log().
		Timestamp().
		Str("id", iom)
}

func (l Logger) Writer(iom string) io.Writer {
	var b bytes.Buffer
	scanner := bufio.NewScanner(&b)

	go func(scanner *bufio.Scanner) {
		for scanner.Scan() {
			line := scanner.Text()
			l.Log(iom).Msg(line)
		}
	}(scanner)

	return &b
}

func (l Logger) StandardLog(iom string) *log.Logger {
	return log.New(l.Writer(iom), "", 0)
}
