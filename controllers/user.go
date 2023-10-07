package controllers

import (
	"fmt"
	"log"
	"net/http"
	"net/url"

	"github.com/dmcclung/pixelparade/context"
	"github.com/dmcclung/pixelparade/errors"
	"github.com/dmcclung/pixelparade/models"
)

type UserTemplates struct {
	SignUp         Template
	SignIn         Template
	Me             Template
	ForgotPassword Template
	CheckEmail     Template
	ResetPassword  Template
}

type User struct {
	Templates            UserTemplates
	UserService          *models.UserService
	SessionService       *models.SessionService
	PasswordResetService *models.PasswordResetService
	EmailService         *models.EmailService
}

func (u User) SignUp(w http.ResponseWriter, r *http.Request) {
	err := u.Templates.SignUp.Execute(w, r, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (u User) SignIn(w http.ResponseWriter, r *http.Request) {
	err := u.Templates.SignIn.Execute(w, r, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (u User) ResetPassword(w http.ResponseWriter, r *http.Request) {
	var data struct {
		Token  string
		Email  string
		Tab    string
		Header string
	}
	data.Token = r.FormValue("token")
	data.Email = r.FormValue("email")
	data.Tab = "Dashboard"
	data.Header = ""
	err := u.Templates.ResetPassword.Execute(w, r, data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (u User) ProcessResetPassword(w http.ResponseWriter, r *http.Request) {
	var data struct {
		Token    string
		Password string
	}
	data.Token = r.FormValue("token")
	data.Password = r.FormValue("password")

	user, err := u.PasswordResetService.Validate(data.Token)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		return
	}

	err = u.UserService.UpdatePassword(user.ID, data.Password)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Error updating password", http.StatusInternalServerError)
		return
	}

	session, err := u.SessionService.Create(user.ID)
	if err != nil {
		fmt.Println(err)
		http.Redirect(w, r, "/signin", http.StatusFound)
	}
	setCookie(w, CookieSession, session.Token)
	http.Redirect(w, r, "/users/me", http.StatusFound)
}

func (u User) ForgotPassword(w http.ResponseWriter, r *http.Request) {
	err := u.Templates.ForgotPassword.Execute(w, r, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (u User) CheckEmail(w http.ResponseWriter, r *http.Request) {
	err := u.Templates.CheckEmail.Execute(w, r, nil)
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

func (u User) ProcessForgotPassword(w http.ResponseWriter, r *http.Request) {
	email := r.FormValue("email")

	log.Printf("forgot password email %v", email)

	reset, err := u.PasswordResetService.Create(email)
	if err != nil {
		log.Printf("create password reset token: %v", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	queryParams := url.Values{
		"email": {email},
		"token": {reset.Token},
	}
	resetLink := fmt.Sprintf("http://localhost:3000/reset-password?%v", queryParams.Encode())

	err = u.EmailService.SendResetEmail(email, resetLink)
	if err != nil {
		log.Printf("send reset email: %v", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/check-email", http.StatusSeeOther)
}

func (u User) ProcessSignIn(w http.ResponseWriter, r *http.Request) {
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

func (u User) ProcessSignOut(w http.ResponseWriter, r *http.Request) {
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

func (u User) ProcessSignUp(w http.ResponseWriter, r *http.Request) {
	email := r.FormValue("email")
	password := r.FormValue("password")
	user, err := u.UserService.Create(email, password)
	if err != nil {
		if errors.Is(err, models.ErrEmailTaken) {
			err = errors.Public(err, "That email address is already associated with an account")
		}
		err := u.Templates.SignUp.Execute(w, r, nil, err)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
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
