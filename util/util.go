package util

import (
	"crypto/sha256"
	"encoding/hex"
	"math/rand"
)

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

// RandString random string
func RandString(n int) string {
    b := make([]byte, n)
    for i := range b {
        b[i] = letterBytes[rand.Intn(len(letterBytes))]
    }
    return string(b)
}

// HashWithSolt hash with solt
func HashWithSolt(password string,solt string) string {
    h := sha256.New()
    h.Write([]byte(password+solt))
    return hex.EncodeToString(h.Sum(nil))
}