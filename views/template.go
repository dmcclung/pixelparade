package views

import (
	"fmt"
	"html/template"
	"net/http"

	"github.com/dmcclung/pixelparade/templates"
)

type Template struct {
	htmlTmpl *template.Template
}

func Must(t Template, err error) Template {
	if err != nil {
		panic(err)
	}
	return t
}

func Parse(tName ...string) (Template, error) {
	t, err := template.ParseFS(templates.FS, tName...)
	if err != nil {
		return Template{}, fmt.Errorf("parsing template: %v", err)
	}
	return Template{
		htmlTmpl: t,
	}, nil
}

func (t Template) Execute(w http.ResponseWriter, tData interface{}) error {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	err := t.htmlTmpl.Execute(w, tData)
	if err != nil {
		return fmt.Errorf("executing template: %v", err)
	}
	return nil
}
