package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatal("no directory provided")
	}

	start := time.Now()

	rootDir := os.Args[1]
	wwwDir := filepath.Join(rootDir, "www") // MV
	if _, err := os.Stat(wwwDir); os.IsNotExist(err) {
		wwwDir = rootDir // MZ
	}

	dataDir := filepath.Join(wwwDir, "data")

	encryptionKey, err := getEncryptionKey(dataDir)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("encryption key:", encryptionKey)

	files, err := getEncryptedFiles(wwwDir)
	if err != nil {
		log.Fatal(err)
	}

	key, err := calculateKey(encryptionKey)
	if err != nil {
		log.Fatal(err)
	}

	count := len(files)

	for i, file := range files {
		fmt.Printf("%d/%d %s - ", i+1, count, filepath.Base(file))
		err := processFile(file, key)
		if err != nil {
			log.Fatal(err)
		}
	}

	err = removeEncryptionKey(dataDir)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("decrypted %d files in %s\n", count, time.Since(start))
}

func getEncryptedFiles(wwwDir string) ([]string, error) {
	var fileList []string

	err := filepath.Walk(wwwDir, func(path string, f os.FileInfo, err error) error {
		if f.IsDir() || !isEncryptedFile(filepath.Ext(path)) {
			return nil
		}

		fileList = append(fileList, path)

		return nil
	})

	if err != nil {
		log.Fatal(err)
	}

	return fileList, nil
}
