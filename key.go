package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strconv"

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

func calculateKey(decryptionKey string) ([]byte, error) {
	var i int
	var chunk string
	var key []byte

	if len(decryptionKey)/2 != 16 {
		return nil, fmt.Errorf("invalid key provided")
	}

	for _, n := range decryptionKey {
		if i == 2 {
			num, err := strconv.ParseInt(chunk, 16, 32)
			if err != nil {
				log.Fatal(err)
			}

			key = append(key, byte(num))

			i = 0
			chunk = ""
		}

		i++
		chunk += string(n)
	}

	if i == 2 {
		num, err := strconv.ParseInt(chunk, 16, 32)
		if err != nil {
			log.Fatal(err)
		}

		key = append(key, byte(num))
	}

	return key, nil
}
