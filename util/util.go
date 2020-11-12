package util

import (
	"crypto/sha256"
	"encoding/hex"
	"math/rand"
	"os"
	"path/filepath"
	"strings"
)

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

var (
    // Version  version
    Version = "No Version Id"
    // BuildTime build time
    BuildTime = "No Build Time"
);

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
    return func (callback func(key int))  {
        for i:=0;i<steps;i++ {
            callback(i)
        }
    }
}


// FileCtime file ctime
// func FileCtime(filename string) (sec int64,err error) {
//     st := &syscall.Stat_t{}
//     if err = syscall.Stat(filename,st) ;err != nil {
//         return 0,err
//     }
//     return st.Ctimespec.Sec,nil
// }

var downloadCodeMap = map[rune]([]byte) {
    '0':{0,0,0,0},
    '1':{0,0,0,1},
    '2':{0,0,1,0},
    '3':{0,0,1,1},
    '4':{0,1,0,0},
    '5':{0,1,0,1},
    '6':{0,1,1,0},
    '7':{0,1,1,1},
    '8':{1,0,0,0},
    '9':{1,0,0,1},
    'A':{1,0,1,0},
    'B':{1,0,1,1},
    'C':{1,1,0,0},
    'D':{1,1,0,1},
    'E':{1,1,1,0},
    'F':{1,1,1,1},
}

// ParseDownloadCode parse download code
func ParseDownloadCode(code string) *[]byte {
    result := []byte{}
    for _,char := range(strings.ToUpper(code)) {
        current := downloadCodeMap[char]
        result = append(result, current...)
    }
    length := len(result)
    for i:=0;i<length/2;i++ {
        tmp := result[i]
        result[i] = result[length - 1 - i]
        result[length - 1 -i] = tmp
    }
    return &result
}