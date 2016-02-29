package etcdexpose

import (
	"bytes"
	"log"
	"net/http"
	"text/template"
	"time"
)

const URL_TEMPLATE = "http://{{.Value}}:{{.Port}}{{.Path}}"

type HealthCheck struct {
	client     *http.Client
	path       string
	port       uint
	retry      uint
	retryDelay time.Duration
	tmpl       *template.Template
}

type urlMembers struct {
	Value string
	Path  string
	Port  uint
}

func NewHealthCheck(
	path string,
	port, retry uint,
	retryDelay,
	timeout time.Duration) *HealthCheck {
	tmpl, _ := template.New("url").Parse(URL_TEMPLATE)
	return &HealthCheck{
		client:     &http.Client{Timeout: 5 * time.Second},
		path:       path,
		tmpl:       tmpl,
		port:       port,
		retry:      retry,
		retryDelay: retryDelay,
	}
}

func (p *HealthCheck) Do(value string) error {
	url, err := p.renderUrl(
		&urlMembers{Value: value, Path: p.path, Port: p.port},
	)

	if err != nil {
		return err
	}

	return p.test(url, 0)
}

func (p *HealthCheck) test(url string, attempt uint) error {
	log.Printf(
		"Performing healthcheck at url [%s], attempt [%d]/[%d]",
		url,
		attempt,
		p.retry,
	)
	_, err := p.client.Get(url)

	if err != nil && attempt < p.retry {
		time.Sleep(p.retryDelay)
		attempt += 1
		err = p.test(url, attempt)
	}

	return err
}

func (p *HealthCheck) renderUrl(value *urlMembers) (string, error) {
	b := &bytes.Buffer{}
	err := p.tmpl.Execute(b, value)
	return b.String(), err
}
