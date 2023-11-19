package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
)

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
	part, err := writer.CreateFormFile("file", filePath)
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
		"name": "File name",
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
