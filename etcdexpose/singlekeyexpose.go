package etcdexpose

import (
	"errors"
	"github.com/coreos/go-etcd/etcd"
	"log"
)

type SingleKeyExpose struct {
	client    *etcd.Client
	renderer  *ValueRenderer
	ping      *Ping
	namespace string
	key       string
}

func NewSingleKeyExpose(client *etcd.Client,
	namespace string,
	renderer *ValueRenderer,
	ping *Ping,
	key string,
) *SingleKeyExpose {
	return &SingleKeyExpose{
		client:    client,
		namespace: namespace,
		renderer:  renderer,
		ping:      ping,
		key:       key}
}

func (s *SingleKeyExpose) Perform(e *etcd.Response) error {
	resp, err := s.client.Get(s.namespace, false, false)
	if err != nil {
		return err
	}

	if resp.Node.Nodes.Len() == 0 {
		return errors.New("No key to expose in given namespace")
	}

	pick := s.pickNode(resp.Node.Nodes)

	if pick == nil {
		return errors.New("Unable to find a valid node in given namespace")
	}

	val, err := s.renderer.Perform(pick.Value)

	if err != nil {
		return err
	}

	_, setErr := s.client.Set(s.key, val, 0)

	if setErr != nil {
		return setErr
	}

	log.Printf("Updated %s to %s", s.key, val)

	return nil
}

func (s *SingleKeyExpose) pickNode(nodes etcd.Nodes) *etcd.Node {
	var pick *etcd.Node = nil
	for _, node := range nodes {
		_, err := s.ping.Do(node.Value)
		if err == nil {
			pick = node
			break
		}
	}
	return pick
}
