package etcdexpose

import (
	"github.com/coreos/go-etcd/etcd"
	"log"
)

type Handler interface {
	Perform() error
}

type Watcher interface {
	Start(eventChan chan *etcd.Response, failureChan chan error)
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
	eventChan := make(chan *etcd.Response)
	failureChan := make(chan error)

	err := r.handler.Perform()
	if err != nil {
		log.Print(err)
	}

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
