package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func processFile(filePath string, decryptionKey []byte) error {
	outFile := getOutputFilePath(filePath)

	err := decryptFile(filePath, outFile, decryptionKey)
	if err != nil {
		log.Fatal(err)
	}

	// Delete old encrypted file
	err = os.Remove(filePath)
	if err != nil {
		log.Fatal(err)
	}

	return nil
}

func isEncryptedFile(ext string) bool {
	switch ext {
	case ".rpgmvp", ".rpgmvm", ".rpgmvo":
		return true
	}

	return false
}

func getRealExt(oldExt string) (string, error) {
	switch strings.ToLower(oldExt) {
	case ".rpgmvp":
		return ".png", nil
	case ".rpgmvm":
		return ".m4a", nil
	case ".rpgmvo":
		return ".ogg", nil
	}

	return "", fmt.Errorf("unknown extension")
}

func getOutputFilePath(filePath string) string {
	oldExt := filepath.Ext(filePath)
	newExt, err := getRealExt(oldExt)
	if err != nil {
		log.Fatal(err)
	}

	fileName := filepath.Base(filePath)
	fileName = fileName[0 : len(fileName)-len(oldExt)]
	fileName = fileName + newExt

	filePath = filepath.Join(filepath.Dir(filePath), fileName)

	return filePath
}

func readFileContents(filePath string) ([]byte, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("can't open %s: %s", filePath, err)
	}
	defer file.Close()

	bytes, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read %s: %s", filePath, err)
	}

	return bytes, nil
}

func writeFileContents(filePath string, content *[]byte) error {
	err := ioutil.WriteFile(filePath, *content, 0644)
	if err != nil {
		return err
	}

	return nil
}
