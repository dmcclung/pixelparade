package controllers

import (
	"fmt"
	"log"
	"net/http"

	"github.com/dmcclung/pixelparade/models"
)

type UserTemplates struct {
	Signup Template
	Signin Template
}

type User struct {
	Templates      UserTemplates
	UserService    *models.UserService
	SessionService *models.SessionService
}

func (u User) GetSignup(w http.ResponseWriter, r *http.Request) {
	err := u.Templates.Signup.Execute(w, r, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (u User) GetSignin(w http.ResponseWriter, r *http.Request) {
	err := u.Templates.Signin.Execute(w, r, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (u User) CurrentUser(w http.ResponseWriter, r *http.Request) {
	token, err := readCookie(r, CookieSession)
	
	if err != nil {
		fmt.Println(err)
		http.Redirect(w, r, "/signin", http.StatusSeeOther)
		return
	}
	user, err := u.SessionService.User(token)
	if err != nil {
		fmt.Println(err)
		http.Redirect(w, r, "/signin", http.StatusSeeOther)
		return
	}

	fmt.Fprintf(w, "Current user: %s\n", user.Email)
}

func (u User) PostSignin(w http.ResponseWriter, r *http.Request) {
	email := r.FormValue("email")
	password := r.FormValue("password")
	
	log.Printf("email %v password %v", email, password)
	
	user, err := u.UserService.Authenticate(email, password)
	if err != nil {
		log.Printf("error authenticating %v", err)
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	log.Printf("User authenticated %v, %v\n", user.Email, user.Password)

	session, err := u.SessionService.Create(user.ID)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	setCookie(w, "session", session.Token)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (u User) PostSignup(w http.ResponseWriter, r *http.Request) {
	email := r.FormValue("email")
	password := r.FormValue("password")
	user, err := u.UserService.Create(email, password)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	session, err := u.SessionService.Create(user.ID)
	if err != nil {
		fmt.Println(err)
		http.Redirect(w, r, "/signin", http.StatusSeeOther)
		return
	}

	log.Printf("User signup %v, %v\n", user.Email, user.Password)

	setCookie(w, "session", session.Token)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
