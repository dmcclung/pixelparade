package main

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
)

func main() {
	n := 32
	b := make([]byte, n)
	nRead, err := rand.Read(b)
	if nRead < n {
		panic("didn't read enough bytes")
	}

	if err != nil {
		panic(err)
	}

	fmt.Println(base64.URLEncoding.EncodeToString(b))
}
