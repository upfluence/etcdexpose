package etcdexpose

import (
	"errors"
	"github.com/coreos/go-etcd/etcd"
	"log"
	"strings"
)

type MultipleValueExpose struct {
	client    *etcd.Client
	renderer  *ValueRenderer
	ping      *Ping
	namespace string
	key       string
}

func NewMutlipleValueExpose(client *etcd.Client,
	renderer *ValueRenderer,
	ping *Ping,
	namespace string,
	key string,
) *MultipleValueExpose {
	return &MultipleValueExpose{
		client:    client,
		namespace: namespace,
		renderer:  renderer,
		ping:      ping,
		key:       key,
	}
}

func (m *MultipleValueExpose) Perform(e *etcd.Response) error {
	resp, err := m.client.Get(m.namespace, false, false)
	if err != nil {
		return err
	}

	if resp.Node.Nodes.Len() == 0 {
		return errors.New("No key to expose in given namespace")
	}

	picks := m.filterNodes(resp.Node.Nodes)

	if picks.Len() == 0 {
		return errors.New("Failed to find any valid node in given namespace")
	}

	val, err := m.formatNodes(picks)

	if err != nil {
		return err
	}

	_, setErr := m.client.Set(m.key, val, 0)

	if setErr != nil {
		return setErr
	}

	log.Printf("Updated %s to %s", m.key, val)

	return nil
}

func (m *MultipleValueExpose) filterNodes(nodes etcd.Nodes) etcd.Nodes {
	var selection etcd.Nodes
	for _, node := range nodes {
		_, err := m.ping.Do(node.Value)
		if err == nil {
			selection = append(selection, node)
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
