package time

import (
	"time"

	iface "github.com/upfluence/etcdexpose/watcher"
)

type watcher struct {
	interval   time.Duration
	stopChan   chan bool
	bufferSize int
}

func NewWatcher(interval time.Duration, bufferSize int) iface.Watcher {
	return &watcher{
		interval:   interval,
		bufferSize: bufferSize,
	}
}

func (t *watcher) Start() <-chan bool {
	out := make(chan bool, t.bufferSize)
	ticker := time.NewTicker(t.interval)
	t.stopChan = make(chan bool)
	go t.run(out, ticker)
	return out
}

func (t *watcher) Stop() {
	t.stopChan <- true
}

func (t *watcher) run(evt chan<- bool, ticker *time.Ticker) {
	for {
		select {
		case <-ticker.C:
			evt <- true
		case <-t.stopChan:
			ticker.Stop()
			close(evt)
			return
		}
	}
}
