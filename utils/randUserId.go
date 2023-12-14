package utils

import (
	"encoding/base64"
	"math/rand"
)

func RandUserId() string {

	userIdLen := 8
	randomBytes := make([]byte, userIdLen)

	_, _ = rand.Read(randomBytes)

	userID := base64.URLEncoding.EncodeToString(randomBytes)[:userIdLen]

	return userID
}
