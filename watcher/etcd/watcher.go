package etcd

import (
	iface "github.com/upfluence/etcdexpose/watcher"
	"log"

	"github.com/coreos/etcd/client"
	"golang.org/x/net/context"
)

type watcher struct {
	namespace string
	keys      client.KeysAPI
	cancel    context.CancelFunc
}

func NewWatcher(namespace string, keys client.KeysAPI) iface.Watcher {
	return &watcher{namespace, keys, nil}
}

func (w *watcher) Start() <-chan bool {
	// Avoid blocking watcher
	out := make(chan bool, 5)

	ctx, cancel := context.WithCancel(context.Background())
	w.cancel = cancel

	etcdWatcher := w.keys.Watcher(
		w.namespace,
		&client.WatcherOptions{0, false},
	)

	go w.run(etcdWatcher, ctx, out)

	return out
}

func (w *watcher) Stop() {
	if w.cancel != nil {
		w.cancel()
	}
}

func (w *watcher) run(watch client.Watcher, ctx context.Context, out chan<- bool) {
	for {
		_, err := watch.Next(ctx)
		if err != nil {
			log.Printf(
				"Got an error from etcd [%s], closing chan, exiting\n",
				err,
			)
			close(out)
			return
		}
		out <- true
	}
}
