package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/dmcclung/pixelparade/jwt"
	"github.com/joho/godotenv"
	"golang.org/x/oauth2"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
		return
	}

	var (
		ok           bool
		clientID     string
		clientSecret string
	)

	if clientID, ok = os.LookupEnv("APPLE_APP_ID"); !ok {
		log.Fatal("Apple app ID not configured")
		return
	}

	if clientSecret, err = jwt.GenerateJWT(); err != nil {
		log.Fatal(err)
		return
	}

	ctx := context.Background()
	conf := &oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		Scopes:       []string{},
		Endpoint: oauth2.Endpoint{
			AuthURL:  "https://appleid.apple.com/auth/authorize",
			TokenURL: "https://appleid.apple.com/auth/token",
		},
		RedirectURL: "https://pixelparade.xyz/oauth/apple/redirect",
	}

	// Redirect user to consent page to ask for permission
	// for the scopes specified above.
	url := conf.AuthCodeURL("state", oauth2.AccessTypeOffline)
	fmt.Printf("Visit the URL for the auth dialog: %v\n", url)
	fmt.Printf("Once you have a code, paste it and press enter: ")

	// Use the authorization code that is pushed to the redirect
	// URL. Exchange will do the handshake to retrieve the
	// initial access token. The HTTP Client returned by
	// conf.Client will refresh the token as necessary.
	var code string
	if _, err := fmt.Scan(&code); err != nil {
		log.Fatal(err)
	}
	token, err := conf.Exchange(ctx, code)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Access Token: %v", token.AccessToken)
	log.Printf("Expiry: %v", token.Expiry)
	log.Printf("Refresh Token: %v", token.RefreshToken)
	log.Printf("Token Type: %v", token.TokenType)

	idToken := token.Extra("id_token")
	if idToken == nil {
		log.Fatal("Apple ID Token not found")
		return
	}

	userID, err := jwt.GetSubFromJWT(idToken.(string))
	if err != nil {
		log.Fatal(err)
		return
	}

	log.Printf("User ID: %v", userID)
}
