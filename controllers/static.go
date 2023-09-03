package controllers

import (
	"net/http"
)

func Static(tmplt Template) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := tmplt.Execute(w, nil)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}
