package main

import (
	"fmt"
	"log"
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

	apiKey := os.Getenv("PINATA_API_KEY")

	pinataClient := &pinata.Client{
		Jwt: apiKey,
	}	

	testAuthenticationResponse, err := pinataClient.TestAuthentication()
	if err != nil {
		log.Fatal(err)
		return
	}

	log.Println(testAuthenticationResponse.Message)

	wd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
		return
	}

	const galleryID = "391e851a-6821-4762-9995-3f3e133d06a8"
	const imageID = "IMG_7867.jpeg"

	filePath := fmt.Sprintf("%s/images/gallery-%s/%s", wd, galleryID, imageID)

	pinFileResponse, err := pinataClient.PinFile(filePath)
	if err != nil {
		log.Fatal(err)
		return
	}

	log.Printf("CID: %s\n", pinFileResponse.IpfsHash)
}
