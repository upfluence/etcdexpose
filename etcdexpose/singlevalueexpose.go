package etcdexpose

import (
	"errors"
	"github.com/coreos/go-etcd/etcd"
	"log"
)

type SingleValueExpose struct {
	client      *EtcdClient
	renderer    *ValueRenderer
	healthCheck *HealthCheck
}

func NewSingleValueExpose(
	client *EtcdClient,
	renderer *ValueRenderer,
	healthCheck *HealthCheck,
) *SingleValueExpose {
	return &SingleValueExpose{
		client:      client,
		renderer:    renderer,
		healthCheck: healthCheck,
	}
}

func (s *SingleValueExpose) Perform() error {
	resp, err := s.client.ReadNamespace()
	if err != nil {
		return err
	}

	pick := s.pickNode(resp.Node.Nodes)

	if pick == nil {
		s.client.RemoveKey()
		return errors.New("Unable to find a valid node in given namespace")
	}

	val, err := s.renderer.Perform(pick.Value)

	if err != nil {
		return err
	}

	_, err = s.client.WriteValue(val)

	if err != nil {
		return err
	}

	return nil
}

func (s *SingleValueExpose) pickNode(nodes etcd.Nodes) *etcd.Node {
	var pick *etcd.Node = nil
	for _, node := range nodes {
		_, err := s.healthCheck.Do(node.Value)
		if err == nil {
			pick = node
			log.Printf(
				"Picked node %s at address %s",
				node.Key,
				node.Value,
			)
			break
		} else {
			log.Printf("Node %s: Error: %s", node.Key, err.Error())
		}
	}
	return pick
}
