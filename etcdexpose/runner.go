package etcdexpose

import (
	"github.com/coreos/go-etcd/etcd"
	"log"
)

type Handler interface {
	Perform(*etcd.Response)
}

type DumbHandler struct {
}

func (d *DumbHandler) Perform(e *etcd.Response) {
	log.Printf("HI THERE")
}

type Runner struct {
	Watcher *EtcdWatcher
	Handler Handler
}

func NewRunner(watcher *EtcdWatcher, handler Handler) *Runner {
	return &Runner{
		Watcher: watcher,
		Handler: handler,
	}
}

func (r *Runner) Start() {
	go r.Watcher.Start()
	for {
		select {
		case event := <-r.Watcher.EventChan:
			log.Printf("Received a new event %s", event.Action)
			r.Handler.Perform(event)
		case err := <-r.Watcher.ErrorChan:
			log.Fatal("Error %s", err)
		}
	}
}
