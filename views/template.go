package views

import (
	"bytes"
	"errors"
	"fmt"
	"html/template"
	"io"
	"net/http"
	"path/filepath"

	"github.com/dmcclung/pixelparade/context"
	"github.com/dmcclung/pixelparade/models"
	"github.com/dmcclung/pixelparade/templates"
	"github.com/gorilla/csrf"
)

type Template struct {
	htmlTmpl *template.Template
}

type public interface {
	Public() string
}

func Must(t Template, err error) Template {
	if err != nil {
		panic(err)
	}
	return t
}

func Parse(patterns ...string) (Template, error) {
	t := template.New(filepath.Base(patterns[0]))
	t = t.Funcs(
		template.FuncMap{
			"csrfField": func() (template.HTML, error) {
				return "", fmt.Errorf("csrfField not implemented")
			},
			"currentUser": func() (*models.User, error) {
				return nil, fmt.Errorf("currentUser not implemented")
			},
			"errors": func() error {
				return fmt.Errorf("errors not implemented")
			},
		},
	)
	t, err := t.ParseFS(templates.FS, patterns...)
	if err != nil {
		return Template{}, fmt.Errorf("parsing template: %v", err)
	}
	return Template{
		htmlTmpl: t,
	}, nil
}

func errMessages(errs []error) []string {
	var errMsgs []string
	for _, err := range errs {
		var pubErr public
		if errors.As(err, &pubErr) {
			errMsgs = append(errMsgs, pubErr.Public())
		} else {
			fmt.Println(err)
			errMsgs = append(errMsgs, "Something went wrong")
		}
	}
	return errMsgs
}

func (t Template) Execute(w http.ResponseWriter, r *http.Request, data interface{}, errs ...error) error {
	tmpl, err := t.htmlTmpl.Clone()
	if err != nil {
		return fmt.Errorf("cloning template: %w", err)
	}

	msgs := errMessages(errs)

	tmpl = tmpl.Funcs(
		template.FuncMap{
			"csrfField": func() template.HTML {
				return csrf.TemplateField(r)
			},
			"currentUser": func() *models.User {
				return context.User(r.Context())
			},
			"errors": func() []string {
				return msgs
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
