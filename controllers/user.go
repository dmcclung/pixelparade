package controllers

import (
	"fmt"
	"net/http"

	"github.com/dmcclung/pixelparade/models"
)

type User struct {
	Templates struct {
		New Template
	}
	UserService models.UserService
}

func (u User) Create(w http.ResponseWriter, r *http.Request) {
	err := u.Templates.New.Execute(w, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (u User) New(w http.ResponseWriter, r *http.Request) {
	email := r.FormValue("email")
	password := r.FormValue("password")
	fmt.Printf("User signup %v, %v\n", email, password)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}