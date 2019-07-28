// Copied from github.com/kyma-project/kyma/components/binding-usage-controller/pkg/signal/signal.go

package signal

import (
	"context"
	"os"
	"os/signal"
	"syscall"
)

// SetupChannel registered for SIGTERM and Interrupt. A stop channel is returned
// which is closed on one of these signals. If a second signal is caught, the program
// is terminated with exit code 1.
func SetupChannel() (stopCh <-chan struct{}) {
	stop := make(chan struct{})
	c := make(chan os.Signal, 2)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		close(stop)
		<-c
		os.Exit(1) // second signal. Exit directly.
	}()

	return stop
}

// CancelOnInterrupt calls cancel function when os.Interrupt or SIGTERM is received
func CancelOnInterrupt(stopCh <-chan struct{}, ctx context.Context, cancel context.CancelFunc) {
	go func() {
		select {
		case <-ctx.Done():
		case <-stopCh:
			cancel()
		}
	}()
}
