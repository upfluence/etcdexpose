package etcdexpose

import (
	"github.com/coreos/go-etcd/etcd"
	"log"
)

type EtcdWatcher struct {
	Namespace string
	client    *etcd.Client
	stopChan  chan bool
	running   bool
}

func NewEtcdWatcher(namespace string, cli *etcd.Client) *EtcdWatcher {
	return &EtcdWatcher{
		Namespace: namespace,
		client:    cli,
		running:   false,
	}
}

func (e *EtcdWatcher) Start(eventChan chan *etcd.Response, errorChan chan error) {
	log.Printf("Begining to watch key %s", e.Namespace)

	e.stopChan = make(chan bool)
	e.running = true

	_, err := e.client.Watch(
		e.Namespace,
		0,
		true,
		eventChan,
		e.stopChan)
	e.running = false
	errorChan <- err
	return
}

func (e *EtcdWatcher) Stop() {
	if e.running {
		e.stopChan <- true
	}
}
