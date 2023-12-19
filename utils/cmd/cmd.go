package cmd

import (
	"github.com/google/uuid"
	"math/rand"
)

func RandUserId() string {

	charSet := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	userID := ""
	for i := 0; i < 8; i++ {
		randomIndex := rand.Intn(len(charSet))
		char := string(charSet[randomIndex])
		userID += char
	}
	return userID

}

func RandUuid() string {

	return uuid.NewString()

}
