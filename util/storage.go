package util

import (
	"archive/zip"
	"io"
	"log"
	"os"
	"path/filepath"
)

// Zip zip
func Zip(dir, zipFile string) (err error) {
	fz, err := os.Create(zipFile)
    if err != nil {
        log.Fatalf("Create zip file failed: %s\n", err.Error())
    }
    defer fz.Close()

    w := zip.NewWriter(fz)
    defer w.Close()

    err = filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
        if !info.IsDir() {
            fDest, err := w.Create(path[len(dir)+1:])
            if err != nil {
                log.Printf("Create failed: %s\n", err.Error())
                return nil
            }
            fSrc, err := os.Open(path)
            if err != nil {
                log.Printf("Open failed: %s\n", err.Error())
                return nil
            }
            defer fSrc.Close()
            _, err = io.Copy(fDest, fSrc)
            if err != nil {
                log.Printf("Copy failed: %s\n", err.Error())
                return nil
            }
        }
        return nil
		})
		return err
}