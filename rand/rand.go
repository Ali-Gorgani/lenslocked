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
		return nil, fmt.Errorf("could not read random bytes: %w", err)
	}
	if n != nRead {
		return nil, fmt.Errorf("could not read enough random bytes: %d", n)
	}
	return b, nil
}

// String generates a random string of length n.
// The string is generated using the crypto/rand package.
func String(n int) (string, error) {
	b, err := Bytes(n)
	if err != nil {
		return "", fmt.Errorf("could not generate random string: %w", err)
	}
	return base64.URLEncoding.EncodeToString(b), nil
}
