package main

import (
	"fmt"
	"log"

	"github.com/dmcclung/pixelparade/jwt"
)

func main() {
	fmt.Printf("Paste ID token here and press enter: ")

	var idToken string
	if _, err := fmt.Scan(&idToken); err != nil {
		log.Fatal(err)
		return
	}

	subject, err := jwt.GetSubFromJWT(idToken)
	if err != nil {
		log.Fatal(err)
		return
	}

	log.Printf("User ID: %v\n", subject)
}
