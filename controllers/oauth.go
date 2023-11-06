package controllers

import (
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/gorilla/csrf"
	"golang.org/x/oauth2"
)

type Oauth struct {
	ProviderConfigs map[string]*oauth2.Config
}

func (oa Oauth) Connect(w http.ResponseWriter, r *http.Request) {
	provider := chi.URLParam(r, "provider")
	provider = strings.ToLower(provider)

	config, ok := oa.ProviderConfigs[provider];
	if !ok {
		http.Error(w, "Unsupported OAuth2 service", http.StatusBadRequest)
		return
	}

	state := csrf.Token(r)
	setCookie(w, "oauth_state", state)
	url := config.AuthCodeURL(
		state, 
		oauth2.SetAuthURLParam("redirect_uri", "http://localhost:3000/oauth/dropbox/redirect"),
	)
	http.Redirect(w, r, url, http.StatusSeeOther)
}
