package mock

import (
	"time"

	iface "github.com/upfluence/etcdexpose/watcher"
)

type closer struct {
	after time.Duration
}

func NewChanCloser(after time.Duration) iface.Watcher {
	return &closer{after}
}

func (w *closer) Start() <-chan bool {
	out := make(chan bool)
	go func() {
		<-time.After(w.after)
		close(out)
		return
	}()
	return out
}

func (w *closer) Stop() {}
