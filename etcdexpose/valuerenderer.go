package etcdexpose

import (
	"bytes"
	"text/template"
)

type ValueRenderer struct {
	template *template.Template
}

type templateValue struct {
	Value string
}

func NewValueRenderer(tmpl string) (*ValueRenderer, error) {
	t := template.New("value template")
	t, err := t.Parse(tmpl)

	if err != nil {
		return nil, err
	}

	return &ValueRenderer{
		template: t,
	}, nil
}

func (v *ValueRenderer) Perform(value string) (string, error) {
	b := &bytes.Buffer{}
	err := v.template.Execute(b, &templateValue{Value: value})
	return b.String(), err
}
