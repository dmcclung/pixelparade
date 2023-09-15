package rand

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
)

func Bytes(n int) ([]byte, error) {
	b := make([]byte, n)
	nRead, err := rand.Read(b)
	if err != nil {
		return nil, fmt.Errorf("reading random bytes: %w", err)
	}
	if nRead < n {
		return nil, fmt.Errorf("not enough bytes read")
	}
	return b, nil
}

func String(n int) (string, error) {
	b, err := Bytes(n)
	if err != nil {
		return "", fmt.Errorf("creating random string: %w", err)
	}

	str := base64.URLEncoding.EncodeToString(b)
	return str, nil
}
