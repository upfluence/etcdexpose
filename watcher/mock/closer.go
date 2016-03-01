package mock

import (
	iface "github.com/upfluence/etcdexpose/watcher"
)

type closer struct {
}

func NewChanCloser() iface.Watcher {
	return &closer{}
}

func (w *closer) Start() <-chan bool {
	out := make(chan bool)
	go func() {
		close(out)
		return
	}()
	return out
}

func (w *closer) Stop() {}
