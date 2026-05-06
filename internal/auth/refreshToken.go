package auth

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
)

func MakeRefreshToken() string {
	key := make([]byte, 32)
	_, err := rand.Read(key)
	if err != nil {
		fmt.Printf("Failed to make refresh token: %v\n", err)
		return ""
	}
	hexString := hex.EncodeToString(key)

	return hexString
}
