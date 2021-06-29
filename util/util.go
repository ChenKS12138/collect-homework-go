package util

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/afocus/captcha"
)

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

var (
	// Version  version
	Version = "No Version Id"
	// BuildTime build time
	BuildTime = "No Build Time"
)

var (
	// captcha
	CaptchaCap *captcha.Captcha
	// secret
	CapachaSecret string
)

func GenerateCapachaSecret() *string {
	nonce := time.Now().Unix() / int64(time.Minute.Seconds()*0.75)
	secret := CapachaSecret + fmt.Sprintf("%d", nonce)
	return &secret
}

// RandString random string
func RandString(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}

// HashWithSolt hash with solt
func HashWithSolt(password string, solt string) string {
	h := sha256.New()
	h.Write([]byte(password + solt))
	return hex.EncodeToString(h.Sum(nil))
}

//DirSizeB getFileSize get file size by path(B)
func DirSizeB(path string) (int64, error) {
	var size int64
	err := filepath.Walk(path, func(_ string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			size += info.Size()
		}
		return err
	})
	return size, err
}

// Iterator iterator
func Iterator(steps int) func(callback func(key int)) {
	return func(callback func(key int)) {
		for i := 0; i < steps; i++ {
			callback(i)
		}
	}
}

// generateKey generate key
func generateKey(secret *string) []byte {
	h := sha256.New()
	h.Write([]byte(*secret))
	return h.Sum(nil)
}

// Encrypt encrypt
func Encrypt(secret *string, message *string) (*string, error) {
	key := generateKey(secret)
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	ciphertext := make([]byte, aes.BlockSize+len(*message))
	iv := ciphertext[:aes.BlockSize]
	cfb := cipher.NewCFBEncrypter(block, iv)
	cfb.XORKeyStream(ciphertext[aes.BlockSize:], []byte(*message))
	enc := base64.StdEncoding.EncodeToString(ciphertext)
	return &enc, nil
}

// Decrypt decrypt
func Decrypt(secret *string, encrypted *string) (*string, error) {
	key := generateKey(secret)
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	text, err := base64.StdEncoding.DecodeString(*encrypted)
	if err != nil {
		return nil, err
	}
	if len(text) < aes.BlockSize {
		return nil, errors.New("ciphertext too short")
	}
	iv := text[:aes.BlockSize]
	text = text[aes.BlockSize:]
	cfb := cipher.NewCFBDecrypter(block, iv)
	cfb.XORKeyStream(text, text)
	msg := string(text)
	return &msg, nil
}

// FileCtime file ctime
// func FileCtime(filename string) (sec int64,err error) {
//     st := &syscall.Stat_t{}
//     if err = syscall.Stat(filename,st) ;err != nil {
//         return 0,err
//     }
//     return st.Ctimespec.Sec,nil
// }

var downloadCodeMap = map[rune]([]byte){
	'0': {0, 0, 0, 0},
	'1': {0, 0, 0, 1},
	'2': {0, 0, 1, 0},
	'3': {0, 0, 1, 1},
	'4': {0, 1, 0, 0},
	'5': {0, 1, 0, 1},
	'6': {0, 1, 1, 0},
	'7': {0, 1, 1, 1},
	'8': {1, 0, 0, 0},
	'9': {1, 0, 0, 1},
	'A': {1, 0, 1, 0},
	'B': {1, 0, 1, 1},
	'C': {1, 1, 0, 0},
	'D': {1, 1, 0, 1},
	'E': {1, 1, 1, 0},
	'F': {1, 1, 1, 1},
}

// ParseDownloadCode parse download code
func ParseDownloadCode(code string) *[]byte {
	result := []byte{}
	for _, char := range strings.ToUpper(code) {
		current := downloadCodeMap[char]
		result = append(result, current...)
	}
	length := len(result)
	for i := 0; i < length/2; i++ {
		tmp := result[i]
		result[i] = result[length-1-i]
		result[length-1-i] = tmp
	}
	return &result
}
