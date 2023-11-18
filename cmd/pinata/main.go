package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
		return
	}

	apiKey := os.Getenv("PINATA_API_KEY")

	req, err := http.NewRequest(
		"GET", 
		"https://api.pinata.cloud/data/testAuthentication", 
		nil,
	)
	if err != nil {
		log.Fatal(err)
		return
	}

	req.Header.Add("accept", "application/json")
	req.Header.Add("authorization", fmt.Sprintf("Bearer %s", apiKey))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
		return
	}
	
	log.Println(string(body))
}