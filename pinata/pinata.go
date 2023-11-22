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

type PinFileResponse struct {
	IpfsHash  string `json:"IpfsHash"`
	PinSize   string `json:"PinSize"`
	Timestamp string `json:"Timestamp"`
}

type TestAuthenticationResponse struct {
	Message string `json:"message"`
}

func (c *Client) TestAuthentication() (*TestAuthenticationResponse, error) {
	url := "https://api.pinata.cloud/data/testAuthentication"

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("new get request: %w", err)
	}

	req.Header.Add("accept", "application/json")
	req.Header.Add("authorization", fmt.Sprintf("Bearer %s", c.Jwt))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("get request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read response body: %w", err)
	}

	var testAuthenticationResponse TestAuthenticationResponse
	err = json.Unmarshal(body, &testAuthenticationResponse)
	if err != nil {
		return nil, fmt.Errorf("unmarshal test auth: %w", err)
	}
	
	return &testAuthenticationResponse, nil
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

	if resp.StatusCode != 200 {
		return fmt.Errorf("delete response: %d %s", resp.StatusCode, resp.Status)
	}

	log.Printf("delete CID %s: %s\n", cid, resp.Status)
	return nil
}

func (c *Client) PinFile(filePath string) (*PinFileResponse, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("error opening file: %w", err)
	}
	defer file.Close()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("file", filepath.Base(filePath))
	if err != nil {
		return nil, fmt.Errorf("creating form file: %w", err)
	}
	_, err = io.Copy(part, file)
	if err != nil {
		return nil, fmt.Errorf("copying file to form: %w", err)
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
		return nil, fmt.Errorf("closing writer: %w", err)
	}

	url := "https://api.pinata.cloud/pinning/pinFileToIPFS"

	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		return nil, fmt.Errorf("creating request: %w", err)
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Set("Authorization", "Bearer "+c.Jwt)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("making request: %w", err)
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("reading response body: %w", err)
	}

	var pinFileResponse PinFileResponse
	json.Unmarshal(bodyBytes, &pinFileResponse)

	return &pinFileResponse, nil
}
