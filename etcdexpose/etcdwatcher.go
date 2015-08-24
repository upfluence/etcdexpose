package etcdexpose

import (
	"github.com/coreos/go-etcd/etcd"
	"log"
)

type EtcdWatcher struct {
	Namespace string
	client    *etcd.Client
	stopChan  chan bool
}

func NewEtcdWatcher(namespace string, cli *etcd.Client) *EtcdWatcher {
	return &EtcdWatcher{
		Namespace: namespace,
		client:    cli,
	}
}

func (e *EtcdWatcher) Start(eventChan chan *etcd.Response) {
	log.Printf("Begining to watch key %s", e.Namespace)

	e.stopChan = make(chan bool, 1)

	_, err := e.client.Watch(
		e.Namespace,
		0,
		true,
		eventChan,
		e.stopChan)

	log.Printf("EtcdWatcher error: %s", err.Error())
	return
}

func (e *EtcdWatcher) Stop() {
	e.stopChan <- true
}
