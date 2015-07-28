package etcdexpose

import (
	"github.com/coreos/go-etcd/etcd"
	"time"
)

type TimeWatcher struct {
	ticker   *time.Ticker
	interval time.Duration
	unit     time.Duration
	stopChan chan bool
}

func NewTimeWatcher(interval time.Duration, unit time.Duration) *TimeWatcher {
	return &TimeWatcher{
		interval: interval,
		unit:     unit,
	}
}

func (t *TimeWatcher) Start(eventChan chan *etcd.Response) {
	t.ticker = time.NewTicker(t.interval * t.unit)
	t.stopChan = make(chan bool)

	for {
		select {
		case <-t.ticker.C:
			eventChan <- nil
		case <-t.stopChan:
			return
		}
	}
}

func (t *TimeWatcher) Stop() {
	t.stopChan <- true
	t.ticker.Stop()
}
