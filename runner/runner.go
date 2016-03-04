package runner

import (
	"log"
	"reflect"

	"github.com/upfluence/etcdexpose/handler"
	"github.com/upfluence/etcdexpose/watcher"
)

type Runner struct {
	watchers   []watcher.Watcher
	handler    handler.Handler
	bufferSize int
}

func NewRunner(handler handler.Handler, watchers []watcher.Watcher, bufferSize int) *Runner {
	return &Runner{watchers, handler, bufferSize}
}

func (r *Runner) Start() {
	consumerChan := make(chan bool, r.bufferSize)
	r.handler.Run(consumerChan)

	// Init the consumer
	consumerChan <- true

	eventCases := make([]reflect.SelectCase, len(r.watchers))

	for i, watcher := range r.watchers {
		eventCases[i] = reflect.SelectCase{
			Dir: reflect.SelectRecv,
			Chan: reflect.ValueOf(
				watcher.Start(),
			),
		}

	}

	for {
		chosen, _, ok := reflect.Select(eventCases)
		log.Printf("Received a new event")
		if !ok {
			log.Printf(
				"Spotted a chan close at %d, returning\n",
				chosen,
			)
			close(consumerChan)
			return
		}
		consumerChan <- true
	}
}

func (r *Runner) Stop() {
	for _, watcher := range r.watchers {
		watcher.Stop()
	}
}
