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
	Watcher Watcher
	Handler Handler
}

func NewRunner(watcher Watcher, handler Handler) *Runner {
	return &Runner{
		Watcher: watcher,
		Handler: handler,
	}
}

func (r *Runner) Start() {
	eventChan := make(chan bool)
	failureChan := make(chan error)
	go r.Watcher.Start(eventChan, failureChan)
	for {
		select {
		case <-eventChan:
			log.Printf("Received a new event ")
			err := r.Handler.Perform()
			if err != nil {
				log.Print(err)
			}
			log.Print("Processed event")
		case err := <-failureChan:
			log.Fatal("Error %s", err)
		}
	}
}
