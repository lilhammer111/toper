package main

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
)

func GenerateRandomSecretKey(length int) string {
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		panic(err)
	}
	return hex.EncodeToString(bytes)
}

func main() {
	secretKey := GenerateRandomSecretKey(32)
	fmt.Println(secretKey)
}
