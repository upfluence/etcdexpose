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

	eventCases := make([]reflect.SelectCase, len(r.watchers))

	for i, watcher := range r.watchers {
		localEventChan := make(chan *etcd.Response)
		eventCases[i] = reflect.SelectCase{
			Dir:  reflect.SelectRecv,
			Chan: reflect.ValueOf(localEventChan),
		}

		go watcher.Start(localEventChan)
	}

	for {
		chosen, _, ok := reflect.Select(eventCases)
		log.Printf("Received a new event")
		if !ok {
			log.Printf("Spotted a chan close at %d, returning\n", chosen)
			return
		}
		err := r.handler.Perform()

		if err != nil {
			log.Print(err)
		}
	}
}

func (r *Runner) Stop() {
	for _, watcher := range r.watchers {
		watcher.Stop()
	}
}
