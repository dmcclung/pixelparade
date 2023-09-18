package views

import (
	"bytes"
	"fmt"
	"html/template"
	"io"
	"net/http"

	"github.com/dmcclung/pixelparade/templates"
	"github.com/gorilla/csrf"
)

type Template struct {
	htmlTmpl *template.Template
}

type HeaderData struct {
	Tab string
	Header string
}

func Must(t Template, err error) Template {
	if err != nil {
		panic(err)
	}
	return t
}

func Parse(name ...string) (Template, error) {
	t := template.New(name[0])
	t = t.Funcs(
		template.FuncMap{
			"csrfField": func() (template.HTML, error) {
				return "", fmt.Errorf("csrfField not implemented")
			},
			"header": func(tab, header string) HeaderData {
				return HeaderData{
					Tab: tab,
					Header: header,
				}
			},
		},
	)
	t, err := t.ParseFS(templates.FS, name...)
	if err != nil {
		return Template{}, fmt.Errorf("parsing template: %v", err)
	}
	return Template{
		htmlTmpl: t,
	}, nil
}

func (t Template) Execute(w http.ResponseWriter, r *http.Request, data interface{}) error {
	tmpl, err := t.htmlTmpl.Clone()
	if err != nil {
		return fmt.Errorf("cloning template: %w", err)
	}

	tmpl = tmpl.Funcs(
		template.FuncMap{
			"csrfField": func() template.HTML {
				return csrf.TemplateField(r)
			},
		},
	)

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	var buf bytes.Buffer
	err = tmpl.Execute(&buf, data)
	if err != nil {
		return fmt.Errorf("executing template: %w", err)
	}
	io.Copy(w, &buf)
	return nil
}
