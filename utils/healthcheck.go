package utils

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
	retryDelay, timeout time.Duration) (*HealthCheck, error) {

	tmpl, err := template.New("url").Parse(URL_TEMPLATE)
	if err != nil {
		return nil, err
	}

	return &HealthCheck{
		client:     &http.Client{Timeout: timeout},
		path:       path,
		tmpl:       tmpl,
		port:       port,
		retry:      retry,
		retryDelay: retryDelay,
	}, nil
}

func (p *HealthCheck) Do(host string) error {
	url, err := p.renderUrl(
		&urlMembers{Value: host, Path: p.path, Port: p.port},
	)

	if err != nil {
		return err
	}

	var attempt uint = 0

	for attempt < p.retry {
		log.Printf(
			"Performing retry at url [%s], attempt [%d]/[%d]\n",
			url,
			attempt+1,
			p.retry,
		)

		err = p.test(url)

		if err == nil {
			break
		}

		time.Sleep(p.retryDelay)
		attempt += 1
	}

	return err
}

func (p *HealthCheck) test(url string) error {
	res, err := p.client.Get(url)

	if err != nil {
		return err
	}

	defer res.Body.Close()

	return err
}

func (p *HealthCheck) renderUrl(value *urlMembers) (string, error) {
	b := &bytes.Buffer{}
	err := p.tmpl.Execute(b, value)
	return b.String(), err
}
