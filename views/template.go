package views

import (
	"fmt"
	"html/template"
	"net/http"
)

type Template struct {
	htmlTmpl *template.Template
}

func Parse(tPath string) (Template, error) {
  t, err := template.ParseFiles(tPath)
  if err != nil {
		return Template{}, fmt.Errorf("error parsing template: %v", err)
	}
  return Template{
    htmlTmpl: t,
  }, nil
}

func (t Template) Execute(w http.ResponseWriter, tData interface{}) (error) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	err := t.htmlTmpl.Execute(w, tData)
	if err != nil {
    return fmt.Errorf("error executing template: %v", err)
	}
  return nil
}


