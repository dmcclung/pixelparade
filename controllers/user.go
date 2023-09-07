package controllers

import (
	"log"
	"net/http"

	"github.com/dmcclung/pixelparade/models"
)

type UserTemplates struct {
	New Template
	Signin Template
}

type User struct {
	Templates UserTemplates
	UserService models.UserService
}

func (u User) Create(w http.ResponseWriter, r *http.Request) {
	err := u.Templates.New.Execute(w, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (u User) Signin(w http.ResponseWriter, r *http.Request) {
	err := u.Templates.Signin.Execute(w, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (u User) ProcessSignin(w http.ResponseWriter, r *http.Request) {
	email := r.FormValue("email")
	password := r.FormValue("password")
	log.Printf("email %v password %v", email, password)
	user, err := u.UserService.Authenticate(email, password)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}
	log.Printf("User authenticated %v, %v\n", user.Email, user.Password)
	cookie := http.Cookie{
		Name:  "email",
		Value: user.Email,
		Path:  "/",
	}
	http.SetCookie(w, &cookie)
	http.Redirect(w, r, "/", http.StatusSeeOther)
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