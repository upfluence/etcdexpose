package utils

import (
	"bytes"
	"text/template"
)

type ValueRenderer struct {
	template *template.Template
	port     uint
}

type templateValue struct {
	Value string
	Port  uint
}

func NewValueRenderer(tmpl string, port uint) (*ValueRenderer, error) {
	t := template.New("value template")
	t, err := t.Parse(tmpl)

	if err != nil {
		return nil, err
	}

	return &ValueRenderer{
		template: t,
		port:     port,
	}, nil
}

func (v *ValueRenderer) Perform(value string) (string, error) {
	b := &bytes.Buffer{}
	err := v.template.Execute(b, &templateValue{Value: value, Port: v.port})
	return b.String(), err
}
