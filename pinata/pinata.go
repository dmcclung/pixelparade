package pinata

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
)

type Client struct {
	Jwt string
}

func (c *Client) UnPinFile(cid string) error {
	url := fmt.Sprintf("https://api.pinata.cloud/pinning/unpin/%s", cid)

	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return fmt.Errorf("new delete request: %w", err)
	}

	req.Header.Add("accept", "application/json")
	req.Header.Add("Authorization", "Bearer "+c.Jwt)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("delete http: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("read delete response: %w", err)
	}

	log.Println(string(body))
	return nil
}

func (c *Client) PinFile(filePath string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("error opening file: %w", err)
	}
	defer file.Close()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("file", filepath.Base(filePath))
	if err != nil {
		return fmt.Errorf("creating form file: %w", err)
	}
	_, err = io.Copy(part, file)
	if err != nil {
		return fmt.Errorf("copying file to form: %w", err)
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
		return fmt.Errorf("closing writer: %w", err)
	}

	url := "https://api.pinata.cloud/pinning/pinFileToIPFS"

	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		return fmt.Errorf("creating request: %w", err)
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Set("Authorization", "Bearer "+c.Jwt)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("making request: %w", err)
	}
	defer resp.Body.Close()

	respBody := &bytes.Buffer{}
	_, err = respBody.ReadFrom(resp.Body)
	if err != nil {
		return fmt.Errorf("reading response body: %w", err)
	}

	// TODO: Need to create a struct to parse the json respBody
	log.Println(respBody.String())
	return nil
}
