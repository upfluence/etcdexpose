package runner

import (
	"testing"
	"time"

	"github.com/upfluence/etcdexpose/watcher"

	mock_handler "github.com/upfluence/etcdexpose/handler/mock"
	mock_watcher "github.com/upfluence/etcdexpose/watcher/mock"
	time_watcher "github.com/upfluence/etcdexpose/watcher/time"
)

func TestStartShouldCallHandlerOnStart(t *testing.T) {
	watchers := []watcher.Watcher{
		mock_watcher.NewChanCloser(1 * time.Millisecond),
	}

	handler := mock_handler.NewHandler(nil)

	runner := NewRunner(handler, watchers)

	runner.Start()

	if handler.CallCount != 1 {
		t.Errorf("Expected 1 call on Start, got [%d] \n", handler.CallCount)
	}
}

func TestStartShouldExitOnChanClose(t *testing.T) {
	watchers := []watcher.Watcher{
		time_watcher.NewWatcher(100*time.Millisecond, 5),
		time_watcher.NewWatcher(200*time.Millisecond, 5),
		mock_watcher.NewChanCloser(500 * time.Millisecond),
	}

	handler := mock_handler.NewHandler(nil)

	runner := NewRunner(handler, watchers)

	runner.Start()

	if handler.CallCount == 0 {
		t.Errorf("Expected at least a call [%d] \n", handler.CallCount)
	}

	runner.Stop()
}

func TestStartShouldBeResistantToRestart(t *testing.T) {
	watchers := []watcher.Watcher{
		time_watcher.NewWatcher(100*time.Millisecond, 5),
		time_watcher.NewWatcher(200*time.Millisecond, 5),
		mock_watcher.NewChanCloser(500 * time.Millisecond),
	}

	handler := mock_handler.NewHandler(nil)
	runner := NewRunner(handler, watchers)

	for i := 0; i < 10; i++ {
		handler.CallCount = 0
		runner.Start()
		if handler.CallCount == 0 {
			t.Errorf("Expected at least a call [%d] \n", handler.CallCount)
		}
		runner.Stop()
	}

}
