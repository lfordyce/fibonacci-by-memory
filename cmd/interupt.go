package cmd

import (
	"os"
	"os/signal"
	"syscall"
)

// RegisterInterrupt listens for CTRL-C and notifies the main application over {progChan}
func RegisterInterrupt() <-chan struct{} {
	intrChan := make(chan os.Signal, 1)
	done := make(chan struct{}, 1)

	signal.Notify(intrChan, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)

	go func() {
		<-intrChan
		done <- struct{}{}
	}()

	return done
}
