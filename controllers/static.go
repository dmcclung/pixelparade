package controllers

import (
  "net/http"
  "github.com/dmcclung/pixelparade/views"
)

type Static struct {
	HtmlTmpl views.Template
}

func (s Static) Get(w http.ResponseWriter, r *http.Request) {
  err := s.HtmlTmpl.Execute(w, nil)
  if err != nil {
    http.Error(w, err.Error(), http.StatusInternalServerError)
  }
}


