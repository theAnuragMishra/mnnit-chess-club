package game

import "crypto/rand"

const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func GenerateGameID(length int) (string, error) {
	code := make([]byte, length)
	charsetLen := byte(len(charset))
	randomBytes := make([]byte, length)

	_, err := rand.Read(randomBytes)
	if err != nil {
		return "", err
	}

	for i, b := range randomBytes {
		code[i] = charset[b%charsetLen]
	}

	return string(code), nil
}
