package main

import (
	"fmt"
	"log"

	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
	// Generate a bcrypt hash of the password,
	// with a work factor (cost) of 14.
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		return "", err
	}

	// Return the string representation of the hashed password.
	return string(hashedPassword), nil
}

func main() {
	password := "admin"
	hashedPassword, err := HashPassword(password)
	if err != nil {
		log.Fatalf("Error hashing password: %v", err)
	}

	fmt.Printf("Original Password: %s\n", password)
	fmt.Printf("Hashed Password: %s\n", hashedPassword)
}
