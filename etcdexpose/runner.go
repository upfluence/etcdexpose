package etcdexpose

import (
	"log"
)

type Handler interface {
	Perform() error
}

type Watcher interface {
	Start(eventChan chan bool, failureChan chan error)
	Stop()
}

type Runner struct {
	handler  Handler
	watchers []Watcher
}

func NewRunner(handler Handler) *Runner {
	return &Runner{
		watchers: []Watcher{},
		handler:  handler,
	}
}

func (r *Runner) AddWatcher(watcher Watcher) {
	r.watchers = append(r.watchers, watcher)
}

func (r *Runner) Start() {
	eventChan := make(chan bool)
	failureChan := make(chan error)

	for _, watcher := range r.watchers {
		go watcher.Start(eventChan, failureChan)
	}

	for {
		select {
		case <-eventChan:
			log.Printf("Received a new event ")
			err := r.handler.Perform()
			if err != nil {
				log.Print(err)
			}
			log.Print("Processed event")
		case err := <-failureChan:
			log.Fatal("Error %s", err)
		}
	}
}

func (r *Runner) Stop() {
	for _, watcher := range r.watchers {
		watcher.Stop()
	}
}
