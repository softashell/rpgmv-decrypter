package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/Jeffail/gabs"
)

func getEncryptionKey(dataDir string) (string, error) {
	systemFile := filepath.Join(dataDir, "System.json")
	file, err := os.Open(systemFile)
	if err != nil {
		return "", fmt.Errorf("can't open %s", systemFile)
	}
	defer file.Close()

	jsonBytes, err := ioutil.ReadFile(systemFile)
	if err != nil {
		return "", fmt.Errorf("failed to read %s", systemFile)
	}

	jsonParsed, err := gabs.ParseJSON(jsonBytes)
	key, ok := jsonParsed.Path("encryptionKey").Data().(string)
	if !ok {
		return "", fmt.Errorf("failed to get encryption key")
	}

	return key, nil
}

func removeEncryptionKey(dir string) error {
	return nil
}
