package controllers

import (
	"fmt"
	"net/http"
)

const (
	CookieSession = "session"
)

func newCookie(name, value string) *http.Cookie {
	c := http.Cookie{
		Name:     name,
		Value:    value,
		Path:     "/",
		HttpOnly: true,
	}
	return &c
}

func setCookie(w http.ResponseWriter, name, value string) {
	c := newCookie(name, value)
	http.SetCookie(w, c)
}

func deleteCookie(w http.ResponseWriter, name string) {
	c := newCookie(name, "")
	c.MaxAge = -1
	http.SetCookie(w, c)
}

func readCookie(r *http.Request, name string) (string, error) {
	c, err := r.Cookie(name)
	if err != nil {
		return "", fmt.Errorf("%s: %w", name, err)
	}
	return c.Value, nil
}
