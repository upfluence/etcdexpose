package multiple

import (
	"github.com/coreos/etcd/client"
	"github.com/upfluence/etcdexpose/utils"

	"errors"
	"log"
	"strings"
)

type MultipleValueExpose struct {
	client      *utils.EtcdClient
	renderer    *utils.ValueRenderer
	healthCheck *utils.HealthCheck
}

func NewMutlipleValueExpose(
	client *utils.EtcdClient,
	renderer *utils.ValueRenderer,
	healthCheck *utils.HealthCheck,
) *MultipleValueExpose {
	return &MultipleValueExpose{
		client:      client,
		renderer:    renderer,
		healthCheck: healthCheck,
	}
}

func (m *MultipleValueExpose) Run(in <-chan bool) {
	go func() {
		for {
			_, ok := <-in
			if !ok {
				log.Println("in chan closed, exiting")
				return
			}
			err := m.perform()
			if err != nil {
				log.Println(err)
			}
		}
	}()
}

func (m *MultipleValueExpose) perform() error {
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

	_, err = m.client.WriteValue(val)

	if err != nil {
		return err
	}

	return nil
}

func (m *MultipleValueExpose) filterNodes(nodes client.Nodes) client.Nodes {
	var selection client.Nodes
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

func (m *MultipleValueExpose) formatNodes(nodes client.Nodes) (string, error) {
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
