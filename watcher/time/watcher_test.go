package time

import (
	"errors"
	"testing"
	"time"
)

var (
	timeoutError = errors.New("Timeout error")
)

func waitNticks(
	sync chan<- error,
	in <-chan bool,
	expectedTicks int,
	maxWait time.Duration,
) {
	total := 0
	timer := time.NewTimer(maxWait)
	for {
		select {
		case <-in:
			total += 1
			if total == expectedTicks {
				timer.Stop()
				sync <- nil
				return
			}
		case <-timer.C:
			sync <- timeoutError
			return
		}
	}
}

func TestRestart(t *testing.T) {
	var sync chan error
	var out <-chan bool
	watcher := NewWatcher(500 * time.Millisecond)

	for i := 0; i < 10; i++ {
		sync = make(chan error)
		out = watcher.Start()

		go waitNticks(sync, out, 2, 2*time.Second)

		err := <-sync

		if err != nil {
			t.Errorf("Expected no error, got [%v]\n", err)
		}

		watcher.Stop()
	}
}
