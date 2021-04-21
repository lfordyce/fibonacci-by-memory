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

// LoggerOption is an abstraction for providing additional logging configuration options.
type LoggerOption interface {
	Apply(*Logger)
}

// LoggerOptionFunc concrete type to implement the LoggerOption interface.
type LoggerOptionFunc func(*Logger)

// Apply method satisfies the LoggerOption interface
func (fn LoggerOptionFunc) Apply(logger *Logger) {
	fn(logger)
}

// Logger is a convince wrapper around the zerolog implementation.
type Logger struct {
	*zerolog.Logger
}

// WithApp option to include app logging functionality
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

// Log returns a zerolog.Event reference
func (l Logger) Log(attributeID string) *zerolog.Event {
	return l.Logger.Log().
		Timestamp().
		Str("id", attributeID)
}

// Writer returns an io.Writer instance
func (l Logger) Writer(attributeID string) io.Writer {
	var b bytes.Buffer
	scanner := bufio.NewScanner(&b)

	go func(scanner *bufio.Scanner) {
		for scanner.Scan() {
			line := scanner.Text()
			l.Log(attributeID).Msg(line)
		}
	}(scanner)

	return &b
}

// StandardLog returns the standard library log Logger
func (l Logger) StandardLog(attributeID string) *log.Logger {
	return log.New(l.Writer(attributeID), "", 0)
}
