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

	const galleryID = "f61dc186-c1c1-48b4-b65d-6656101bd04f"
	const imageID = "4c7cda78-4bb1-44bd-a246-d16c850a886e.jpg"

	filePath := fmt.Sprintf("%s/images/gallery-%s/%s", wd, galleryID, imageID)

	pinFileResponse, err := pinataClient.PinFile(filePath)
	if err != nil {
		log.Fatal(err)
		return
	}

	log.Printf("CID: %s\n", pinFileResponse.IpfsHash)

	err = pinataClient.UnPinFile(pinFileResponse.IpfsHash)
	if err != nil {
		log.Fatal(err)
		return
	}
}
