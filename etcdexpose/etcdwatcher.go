package etcdexpose

import (
	"github.com/coreos/go-etcd/etcd"
	"log"
)

type EtcdWatcher struct {
	Namespace string
	EventChan chan *etcd.Response
	ErrorChan chan error
	client    *etcd.Client
	stopChan  chan bool
}

func NewEtcdWatcher(namespace string, cli *etcd.Client) *EtcdWatcher {
	return &EtcdWatcher{
		Namespace: namespace,
		EventChan: make(chan *etcd.Response),
		ErrorChan: make(chan error, 1),
		stopChan:  make(chan bool),
		client:    cli,
	}
}

func (e *EtcdWatcher) Start() {
	log.Printf("Begining to watch key %s", e.Namespace)
	for {
		_, err := e.client.Watch(
			e.Namespace,
			0,
			true,
			e.EventChan,
			e.stopChan)
		e.ErrorChan <- err
	}
}
