package controllers

import (
	"log"
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
	user, err := u.UserService.Create(email, password)
	// TODO: error handling could be better here
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	log.Printf("User signup %v, %v\n", user.Email, user.Password)
	// TODO: How to redirect and save sesssion?
	http.Redirect(w, r, "/", http.StatusSeeOther)
}