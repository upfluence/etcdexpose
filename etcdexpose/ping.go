package etcdexpose

import (
	"bytes"
	"net/http"
	"text/template"
)

const URL_TEMPLATE = "http://{{.Value}}{{.Path}}"

type Ping struct {
	client *http.Client
	path   string
	tmpl   *template.Template
}

type urlMembers struct {
	Value string
	Path  string
}

func NewPing(path string) *Ping {
	tmpl, _ := template.New("url").Parse(URL_TEMPLATE)
	return &Ping{
		client: &http.Client{},
		path:   path,
		tmpl:   tmpl,
	}
}

func (p *Ping) Do(value string) (*http.Response, error) {
	url, err := p.renderUrl(&urlMembers{Value: value, Path: p.path})
	if err != nil {
		return nil, err
	}
	return p.client.Get(url)
}

func (p *Ping) renderUrl(value *urlMembers) (string, error) {
	b := &bytes.Buffer{}
	err := p.tmpl.Execute(b, value)
	return b.String(), err
}
