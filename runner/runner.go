package runner

import (
	"log"
	"reflect"

	"github.com/upfluence/etcdexpose/handler"
	"github.com/upfluence/etcdexpose/watcher"
)

type Runner struct {
	handler  handler.Handler
	watchers []watcher.Watcher
}

func NewRunner(handler handler.Handler, watchers []watcher.Watcher) *Runner {
	return &Runner{
		watchers: watchers,
		handler:  handler,
	}
}

func (r *Runner) Start() {
	err := r.handler.Perform()

	if err != nil {
		log.Print(err)
	}

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
			return
		}

		err := r.handler.Perform()

		if err != nil {
			log.Println(err)
		}
	}
}

func (r *Runner) Stop() {
	for _, watcher := range r.watchers {
		watcher.Stop()
	}
}
