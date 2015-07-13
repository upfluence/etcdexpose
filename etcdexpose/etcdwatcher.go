package etcdexpose

import (
	"github.com/coreos/go-etcd/etcd"
	"log"
)

type EtcdWatcher struct {
	Namespace       string
	client          *etcd.Client
	stopChan        chan bool
	convertStopChan chan bool
}

func NewEtcdWatcher(namespace string, cli *etcd.Client) *EtcdWatcher {
	return &EtcdWatcher{
		Namespace:       namespace,
		client:          cli,
		stopChan:        make(chan bool),
		convertStopChan: make(chan bool),
	}
}

func (e *EtcdWatcher) Start(eventChan chan bool, errorChan chan error) {
	log.Printf("Begining to watch key %s", e.Namespace)

	respChan := make(chan *etcd.Response)

	// Meh, should find a better way to convert an etcdresponse to a bool
	go func(
		respChan chan *etcd.Response,
		eventChan chan bool,
		stopChan chan bool,
	) {
		for {
			<-respChan
			eventChan <- true

		}
	}(respChan, eventChan, e.convertStopChan)

	for {
		_, err := e.client.Watch(
			e.Namespace,
			0,
			true,
			respChan,
			e.stopChan)
		errorChan <- err
	}
}

func (e *EtcdWatcher) Stop() {
	e.stopChan <- true
	e.convertStopChan <- true
}
