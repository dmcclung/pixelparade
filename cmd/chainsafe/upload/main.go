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

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
		return
	}

	apiURL := os.Getenv("CHAINSAFE_API_URL")
	apiKey := os.Getenv("CHAINSAFE_API_KEY")
	bucketID := os.Getenv("CHAINSAFE_IMAGES_BUCKET_ID")

	const galleryID = "f61dc186-c1c1-48b4-b65d-6656101bd04f"
	const imageID = "4c7cda78-4bb1-44bd-a246-d16c850a886e"

	wd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
		return
	}

	filePath := fmt.Sprintf("%s/images/gallery-%s/%s.jpg", wd, galleryID, imageID)

	// Create a buffer to write our multipart form data
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	// Add the file field
	file, err := os.Open(filePath)
	if err != nil {
		log.Fatal(err)
		return
	}
	defer file.Close()

	part, err := writer.CreateFormFile("file", filepath.Base(file.Name()))
	if err != nil {
		log.Fatal(err)
		return
	}
	_, err = io.Copy(part, file)
	if err != nil {
		log.Fatal(err)
		return
	}

	// Add the path field
	err = writer.WriteField("path", "/")
	if err != nil {
		log.Fatal(err)
		return
	}

	// Close the writer before making the request
	err = writer.Close()
	if err != nil {
		log.Fatal(err)
		return
	}

	// Create the request
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/bucket/%s/upload", apiURL, bucketID), body)
	if err != nil {
		log.Fatal(err)
		return
	}

	// Set the content type, this must be done after writer.Close()
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Set("Authorization", "Bearer "+apiKey)

	// Make the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
		return
	}
	defer resp.Body.Close()

	// Read the response
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
		return
	}

	log.Println("Response:", string(respBody))
}
