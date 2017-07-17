package main

import (
	"fmt"
	"log"
	"path/filepath"
	"strings"
)

func processFile(filePath string, decryptionKey []string) error {
	outFile := getOutputFilePath(filePath)

	err := decryptFile(filePath, outFile, decryptionKey)
	if err != nil {
		log.Fatal(err)
	}

	return nil
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
