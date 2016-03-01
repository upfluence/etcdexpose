package etcd

import (
	"github.com/coreos/etcd/client"
	iface "github.com/upfluence/etcdexpose/watcher"
	"golang.org/x/net/context"
	"log"
)

type watcher struct {
	keys       client.KeysAPI
	namespace  string
	bufferSize int
	cancel     context.CancelFunc
}

func NewWatcher(keys client.KeysAPI, namespace string, bufferSize int) iface.Watcher {
	return &watcher{keys, namespace, bufferSize, nil}
}

func (w *watcher) Start() <-chan bool {
	out := make(chan bool, w.bufferSize)

	ctx, cancel := context.WithCancel(context.Background())
	w.cancel = cancel

	etcdWatcher := w.keys.Watcher(
		w.namespace,
		&client.WatcherOptions{0, true},
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
