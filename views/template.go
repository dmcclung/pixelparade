package views

import (
	"bytes"
	"fmt"
	"html/template"
	"io"
	"net/http"

	"github.com/dmcclung/pixelparade/context"
	"github.com/dmcclung/pixelparade/models"
	"github.com/dmcclung/pixelparade/templates"
	"github.com/gorilla/csrf"
)

type Template struct {
	htmlTmpl *template.Template
}

type HeaderData struct {
	Tab    string
	Header string
	User   *models.User
}

func Must(t Template, err error) Template {
	if err != nil {
		panic(err)
	}
	return t
}

func dict(values ...interface{}) (map[string]interface{}, error) {
	if len(values)%2 != 0 {
		return nil, fmt.Errorf("dictionary parameters must be key value pairs")
	}
	dict := make(map[string]interface{})
	for i := 0; i < len(values) - 1; i += 2 {
		key, ok := values[i].(string)
		if !ok {
			return nil, fmt.Errorf("dictionary keys must be strings")
		}
		dict[key] = values[i+1]
	}
	return dict, nil
}

func Parse(name ...string) (Template, error) {
	t := template.New(name[0])
	t = t.Funcs(
		template.FuncMap{
			"csrfField": func() (template.HTML, error) {
				return "", fmt.Errorf("csrfField not implemented")
			},
			"currentUser": func() (*models.User, error) {
				return nil, fmt.Errorf("currentUser not implemented")
			},
			"dict": dict,
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
			"currentUser": func() *models.User {
				return context.User(r.Context())
			},
			"dict": dict,
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
