package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/dmcclung/pixelparade/pinata"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
		return
	}

	jwt := os.Getenv("PINATA_API_KEY")

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
	req.Header.Add("authorization", fmt.Sprintf("Bearer %s", jwt))

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

	pinataClient := &pinata.Client{
		Jwt: jwt,
	}
	pinataClient.PinFile(filePath)
}
