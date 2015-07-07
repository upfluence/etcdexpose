package etcdexpose

import (
	"github.com/coreos/go-etcd/etcd"
)

type Handler interface {
	Handle(*etcd.Response)
}

type Runner struct {
	Watcher *EtcdWatcher
	Handler *Handler
}

func NewRunner(watcher *EtcdWatcher, handler *Handler) (*Runner, error) {
	return &Runner{Watcher: nil, Handler: nil}, nil
}
