package etcdexpose

import (
	"github.com/coreos/go-etcd/etcd"
	"log"
	"reflect"
)

type Handler interface {
	Perform() error
}

type Watcher interface {
	Start(eventChan chan *etcd.Response)
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
	err := r.handler.Perform()

	if err != nil {
		log.Print(err)
	}

	event_cases := make([]reflect.SelectCase, len(r.watchers))

	for i, watcher := range r.watchers {
		localEventChan := make(chan *etcd.Response)
		event_cases[i] = reflect.SelectCase{
			Dir:  reflect.SelectRecv,
			Chan: reflect.ValueOf(localEventChan),
		}

		go watcher.Start(localEventChan)
	}

	for {
		_, _, ok := reflect.Select(event_cases)
		log.Printf("Received a new event")
		if !ok {
			log.Printf("Spotted a chan close, returning")
			return
		}
		r.handler.Perform()
	}
}

func (r *Runner) Stop() {
	for _, watcher := range r.watchers {
		watcher.Stop()
	}
}
