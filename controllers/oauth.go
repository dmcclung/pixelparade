package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/gorilla/csrf"
	"golang.org/x/oauth2"
)

type Oauth struct {
	ProviderConfigs map[string]*oauth2.Config
}

func redirectURI(r *http.Request, provider string) string {
	if r.Host == "localhost:3000" {
		return fmt.Sprintf("http://localhost:3000/oauth/%s/redirect", provider)
	}
	return fmt.Sprintf("https://pixelparade.xyz/oauth/%s/redirect",  provider)
}

func (oa Oauth) Connect(w http.ResponseWriter, r *http.Request) {
	provider := chi.URLParam(r, "provider")
	provider = strings.ToLower(provider)

	config, ok := oa.ProviderConfigs[provider];
	if !ok {
		http.Error(w, "Unsupported OAuth2 service", http.StatusBadRequest)
		return
	}
	
	// use PKCE to protect against CSRF attacks
	// https://www.ietf.org/archive/id/draft-ietf-oauth-security-topics-22.html#name-countermeasures-6
	// TODO: Store this in the database for exchange?
	// verifier := oauth2.GenerateVerifier()

	state := csrf.Token(r)
	setCookie(w, "oauth_state", state)
	url := config.AuthCodeURL(
		state, 
		oauth2.SetAuthURLParam("redirect_uri", redirectURI(r, provider)),
		oauth2.AccessTypeOffline, 
		oauth2.SetAuthURLParam("token_access_type", "offline"),
		// oauth2.S256ChallengeOption(verifier),
	)
	http.Redirect(w, r, url, http.StatusSeeOther)
}

func (oa Oauth) Redirect(w http.ResponseWriter, r *http.Request) {
	provider := chi.URLParam(r, "provider")
	provider = strings.ToLower(provider)

	config, ok := oa.ProviderConfigs[provider]
	if !ok {
		http.Error(w, "Unsupported OAuth2 service", http.StatusBadRequest)
		return
	}

	state := r.FormValue("state")
	cookieState, err := readCookie(r, "oauth_state")
	if err != nil || cookieState != state {
		if err != nil {
			fmt.Println(err)
		}
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}
	deleteCookie(w, "oauth_state")

	code := r.FormValue("code")
	token, err := config.Exchange(r.Context(), code)
	if err != nil {
		http.Error(w, "Something went wrong", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	enc.Encode(token)
}
