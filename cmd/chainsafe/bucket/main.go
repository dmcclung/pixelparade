package main

import (
	"bytes"
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

	apiURL := os.Getenv("CHAINSAFE_API_URL")
	apiKey := os.Getenv("CHAINSAFE_API_KEY")

	bucketName := "images"

	jsonData := []byte(fmt.Sprintf(`{
		"type": "fps",
		"name": "%s"
	}`, bucketName))

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/buckets", apiURL), bytes.NewBuffer(jsonData))
	if err != nil {
		log.Fatal(err)
		return
	}

	req.Header.Set("Authorization", "Bearer "+apiKey)
	req.Header.Set("Content-Type", "application/json")

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

	log.Println("Response:", string(body))
}
