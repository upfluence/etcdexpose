package etcdexpose

import (
	"github.com/coreos/etcd/client"
	"github.com/upfluence/etcdexpose/utils"

	"errors"
	"log"
)

type SingleValueExpose struct {
	client      *utils.EtcdClient
	renderer    *utils.ValueRenderer
	healthCheck *utils.HealthCheck
}

func NewSingleValueExpose(
	client *utils.EtcdClient,
	renderer *utils.ValueRenderer,
	healthCheck *utils.HealthCheck,
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

func (s *SingleValueExpose) pickNode(nodes client.Nodes) *client.Node {
	var pick *client.Node = nil
	for _, node := range nodes {
		err := s.healthCheck.Do(node.Value)
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
