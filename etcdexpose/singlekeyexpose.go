package etcdexpose

import (
	"bytes"
	"errors"
	"github.com/coreos/go-etcd/etcd"
	"log"
	"text/template"
)

type SingleKeyExpose struct {
	client    *etcd.Client
	template  *template.Template
	ping      *Ping
	namespace string
	key       string
}

type templateValue struct {
	Value string
}

func NewSingleKeyExpose(client *etcd.Client,
	namespace string,
	template *template.Template,
	ping *Ping,
	key string,
) *SingleKeyExpose {
	return &SingleKeyExpose{
		client:    client,
		namespace: namespace,
		template:  template,
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

	val, err := s.renderValue(pick)

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

func (s *SingleKeyExpose) renderValue(node *etcd.Node) (string, error) {
	b := &bytes.Buffer{}
	err := s.template.Execute(b, &templateValue{Value: node.Value})
	return b.String(), err
}
