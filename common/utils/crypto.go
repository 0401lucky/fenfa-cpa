package utils

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
)

func GenerateAPIKey() (plainKey string, keyHash string, keyPrefix string) {
	bytes := make([]byte, 24)
	rand.Read(bytes)
	randomHex := hex.EncodeToString(bytes)
	plainKey = "sk-cpa-" + randomHex
	keyHash = HashKey(plainKey)
	if len(randomHex) >= 8 {
		keyPrefix = fmt.Sprintf("sk-cpa-%s...%s", randomHex[:4], randomHex[len(randomHex)-4:])
	}
	return
}

func HashKey(key string) string {
	hash := sha256.Sum256([]byte(key))
	return hex.EncodeToString(hash[:])
}
