package controllers

import (
  "net/http"
  "github.com/dmcclung/pixelparade/views"
)

func Static(tmplt views.Template) http.HandlerFunc {
  return func(w http.ResponseWriter, r *http.Request) {
    err := tmplt.Execute(w, nil)
    if err != nil {
      http.Error(w, err.Error(), http.StatusInternalServerError)
    }
  }
}


