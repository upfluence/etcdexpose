package etcdexpose

import (
	"bytes"
	"net/http"
	"text/template"
)

const URL_TEMPLATE = "http://{{.Value}}:{{.Port}}{{.Path}}"

type HealthCheck struct {
	client *http.Client
	path   string
	port   uint
	tmpl   *template.Template
}

type urlMembers struct {
	Value string
	Path  string
	Port  uint
}

func NewHealthCheck(path string, port uint) *HealthCheck {
	tmpl, _ := template.New("url").Parse(URL_TEMPLATE)
	return &HealthCheck{
		client: &http.Client{},
		path:   path,
		tmpl:   tmpl,
		port:   port,
	}
}

func (p *HealthCheck) Do(value string) (*http.Response, error) {
	url, err := p.renderUrl(
		&urlMembers{Value: value, Path: p.path, Port: p.port},
	)

	if err != nil {
		return nil, err
	}
	return p.client.Get(url)
}

func (p *HealthCheck) renderUrl(value *urlMembers) (string, error) {
	b := &bytes.Buffer{}
	err := p.tmpl.Execute(b, value)
	return b.String(), err
}
