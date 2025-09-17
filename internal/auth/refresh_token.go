package auth

import (
	"encoding/hex"

	"crypto/rand"
)

func MakeRefreshToken() (string, error) {
	newRandomNumber := make([]byte, 32)
	_, err := rand.Read(newRandomNumber)
	if err != nil {
		return "", err
	}

	hexEncodedNewRandomNumber := hex.EncodeToString(newRandomNumber)
	return hexEncodedNewRandomNumber, nil
}
