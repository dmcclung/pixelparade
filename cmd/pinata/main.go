package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
)

func UnPinFile(cid, jwt string) {
	url := fmt.Sprintf("https://api.pinata.cloud/pinning/unpin/%s", cid)

	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		log.Println(err)
		return
	}

	req.Header.Add("accept", "application/json")
	req.Header.Add("Authorization", "Bearer "+jwt)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Println(err)
		return
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)

	fmt.Println(string(body))
}

func PinFileToIPFS(filePath, jwt string) {
	targetURL := "https://api.pinata.cloud/pinning/pinFileToIPFS"

	file, err := os.Open(filePath)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("file", filepath.Base(filePath))
	if err != nil {
		fmt.Println("Error creating form file:", err)
		return
	}
	_, err = io.Copy(part, file)
	if err != nil {
		fmt.Println("Error copying file to form:", err)
		return
	}

	pinataMetadata := map[string]string{
		"name": filepath.Base(filePath),
	}
	metaJSON, _ := json.Marshal(pinataMetadata)
	_ = writer.WriteField("pinataMetadata", string(metaJSON))

	pinataOptions := map[string]int{
		"cidVersion": 0,
	}
	optionsJSON, _ := json.Marshal(pinataOptions)
	_ = writer.WriteField("pinataOptions", string(optionsJSON))

	err = writer.Close()
	if err != nil {
		fmt.Println("Error closing writer:", err)
		return
	}

	req, err := http.NewRequest("POST", targetURL, body)
	if err != nil {
		fmt.Println("Error creating request:", err)
		return
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Set("Authorization", "Bearer "+jwt)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error making request:", err)
		return
	}
	defer resp.Body.Close()

	respBody := &bytes.Buffer{}
	_, err = respBody.ReadFrom(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return
	}

	fmt.Println(respBody.String())
}

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

	wd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
		return
	}

	const galleryID = "391e851a-6821-4762-9995-3f3e133d06a8"
	const imageID = "IMG_7867.jpeg"

	filePath := fmt.Sprintf("%s/images/gallery-%s/%s", wd, galleryID, imageID)

	PinFileToIPFS(filePath, apiKey)
}