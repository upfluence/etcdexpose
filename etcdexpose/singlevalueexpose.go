package etcdexpose

import (
	"errors"
	"github.com/coreos/go-etcd/etcd"
)

type SingleValueExpose struct {
	client   *EtcdClient
	renderer *ValueRenderer
	ping     *Ping
}

func NewSingleValueExpose(
	client *EtcdClient,
	renderer *ValueRenderer,
	ping *Ping,
) *SingleValueExpose {
	return &SingleValueExpose{
		client:   client,
		renderer: renderer,
		ping:     ping,
	}
}

func (s *SingleValueExpose) Perform() error {
	resp, err := s.client.ReadNamespace()
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

	_, setErr := s.client.WriteValue(val)

	if setErr != nil {
		return setErr
	}

	return nil
}

func (s *SingleValueExpose) pickNode(nodes etcd.Nodes) *etcd.Node {
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
