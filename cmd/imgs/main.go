package main

import (
	"fmt"
	"github.com/dmcclung/pixelparade/models"
)

func main() {
	gs := models.GalleryService{}
	images, err := gs.Images("d04ddfd3-0774-4a08-bc5e-b3d3afd3a653")
	if err != nil {
		panic(err)
	}
	fmt.Println(images)
}
