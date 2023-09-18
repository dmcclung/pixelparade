package controllers

import (
	"fmt"
	"log"
	"net/http"

	"github.com/dmcclung/pixelparade/context"
	"github.com/dmcclung/pixelparade/models"
)

type UserTemplates struct {
	Signup Template
	Signin Template
	Me Template
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
	user := context.User(r.Context())
	log.Printf("Current user: %s\n", user.Email)
	err := u.Templates.Me.Execute(w, r, user.Email)
	if err != nil {
		log.Printf("error rendering me: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
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

	setCookie(w, CookieSession, session.Token)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (u User) PostSignout(w http.ResponseWriter, r *http.Request) {
	token, err := readCookie(r, CookieSession)
	if err != nil {
		log.Printf("signout: %v\n", err)
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	err = u.SessionService.Delete(token)
	if err != nil {
		log.Printf("deleting session: %v\n", err)
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	log.Printf("deleted session: %v\n", token)
	deleteCookie(w, "session")
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

type UserMiddleware struct {
	SessionService *models.SessionService
}

func (umw UserMiddleware) RequireUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user := context.User(r.Context())
		if user == nil {
			http.Redirect(w, r, "/signin", http.StatusSeeOther)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func (umw UserMiddleware) SetUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token, err := readCookie(r, CookieSession)
		if err != nil {
			next.ServeHTTP(w, r)
			return
		}

		user, err := umw.SessionService.User(token)
		if err != nil {
			next.ServeHTTP(w, r)
			return
		}

		ctx := r.Context()
		ctx = context.WithUser(ctx, user)
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	})
}
