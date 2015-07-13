package etcdexpose

import (
	"github.com/coreos/go-etcd/etcd"
	"time"
)

type TimeWatcher struct {
	interval time.Duration
	unit     time.Duration
}

func NewTimeWatcher(interval time.Duration, unit time.Duration) *TimeWatcher {
	return &TimeWatcher{
		interval: interval,
		unit:     unit,
	}
}

func (t *TimeWatcher) Start(eventChan chan *etcd.Response, _ chan error) {
	for {
		time.Sleep(t.interval * t.unit)
		eventChan <- nil
	}
}
