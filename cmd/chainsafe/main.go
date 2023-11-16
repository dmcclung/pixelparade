package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
)

var apiURL, apiKey string

func createBucket() {
	bucketName := "images"

	jsonData := []byte(fmt.Sprintf(`{
		"type": "fps",
		"name": "%s"
	}`, bucketName))

	req, err := http.NewRequest("POST", apiURL, bytes.NewBuffer(jsonData))
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

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
		return
	}

	apiURL = os.Getenv("CHAINSAFE_API_URL")
	apiKey = os.Getenv("CHAINSAFE_API_KEY")
	
	createBucket()

	filePath := "/Desktop/image.png" // Replace with the path to your file

	// Create a buffer to write our multipart form data
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	// Add the file field
	file, err := os.Open(filePath)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	part, err := writer.CreateFormFile("file", filepath.Base(file.Name()))
	if err != nil {
		panic(err)
	}
	_, err = io.Copy(part, file)
	if err != nil {
		panic(err)
	}

	// Add the path field
	err = writer.WriteField("path", "/")
	if err != nil {
		panic(err)
	}

	// Close the writer before making the request
	err = writer.Close()
	if err != nil {
		panic(err)
	}

	// Create the request
	req, err := http.NewRequest("POST", apiURL, body)
	if err != nil {
		panic(err)
	}

	// Set the content type, this must be done after writer.Close()
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Set("Authorization", "Bearer "+apiKey)

	// Make the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	// Read the response
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	fmt.Println("Response:", string(respBody))
}