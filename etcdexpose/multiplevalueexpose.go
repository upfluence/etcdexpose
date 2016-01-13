package etcdexpose

import (
	"errors"
	"log"
	"strings"

	"github.com/coreos/go-etcd/etcd"
)

type MultipleValueExpose struct {
	client      *EtcdClient
	renderer    *ValueRenderer
	healthCheck *HealthCheck
}

func NewMutlipleValueExpose(
	client *EtcdClient,
	renderer *ValueRenderer,
	healthCheck *HealthCheck,
) *MultipleValueExpose {
	return &MultipleValueExpose{
		client:      client,
		renderer:    renderer,
		healthCheck: healthCheck,
	}
}

func (m *MultipleValueExpose) Perform() error {
	resp, err := m.client.ReadNamespace()
	if err != nil {
		return err
	}

	picks := m.filterNodes(resp.Node.Nodes)

	if picks.Len() == 0 {
		return errors.New("Failed to find any valid node in given namespace")
	}

	val, err := m.formatNodes(picks)

	if err != nil {
		return err
	}

	_, setErr := m.client.WriteValue(val)

	if setErr != nil {
		return setErr
	}

	return nil
}

func (m *MultipleValueExpose) filterNodes(nodes etcd.Nodes) etcd.Nodes {
	var selection etcd.Nodes
	for _, node := range nodes {
		err := m.healthCheck.Do(node.Value)
		if err == nil {
			log.Printf("Node %s marked as valid", node.Key)
			selection = append(selection, node)
		} else {
			log.Printf("Node %s: Error: %s", node.Key, err.Error())
		}
	}
	return selection
}

func (m *MultipleValueExpose) formatNodes(nodes etcd.Nodes) (string, error) {
	urls := []string{}
	for _, node := range nodes {
		val, err := m.renderer.Perform(node.Value)
		if err != nil {
			return "", err
		}
		urls = append(urls, val)
	}
	return strings.Join(urls, ","), nil
}
