package etcdexpose

import (
	"time"
)

type TimeWatcher struct {
	interval time.Duration
	unit     time.Duration
	stopChan chan bool
}

func NewTimeWatcher(interval time.Duration, unit time.Duration) *TimeWatcher {
	return &TimeWatcher{
		interval: interval,
		unit:     unit,
		stopChan: make(chan bool),
	}
}

func (t *TimeWatcher) Start(eventChan chan bool, _ chan error) {
	for {
		time.Sleep(t.interval * t.unit)
		eventChan <- true
	}
}
