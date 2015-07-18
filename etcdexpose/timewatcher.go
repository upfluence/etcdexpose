package etcdexpose

import (
	"github.com/coreos/go-etcd/etcd"
	"time"
)

type TimeWatcher struct {
	ticker *time.Ticker
}

func NewTimeWatcher(interval time.Duration, unit time.Duration) *TimeWatcher {
	return &TimeWatcher{
		ticker: time.NewTicker(interval * unit),
	}
}

func (t *TimeWatcher) Start(eventChan chan *etcd.Response, _ chan error) {
	for {
		<-t.ticker.C
		eventChan <- nil
	}
}

func (t *TimeWatcher) Stop() {
	t.ticker.Stop()
}
