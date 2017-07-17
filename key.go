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

func openJSON(filePath string) (*gabs.Container, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("can't open %s", filePath)
	}
	defer file.Close()

	jsonBytes, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read %s", filePath)
	}

	jsonParsed, err := gabs.ParseJSON(jsonBytes)
	if err != nil {
		return nil, err
	}

	return jsonParsed, nil
}

func getEncryptionKey(dataDir string) (string, error) {
	systemFile := filepath.Join(dataDir, "System.json")
	json, err := openJSON(systemFile)
	if err != nil {
		return "", err
	}

	key, ok := json.Path("encryptionKey").Data().(string)
	if !ok {
		return "", fmt.Errorf("failed to get encryption key")
	}

	return key, nil
}

func removeEncryptionKey(dataDir string) error {
	systemFile := filepath.Join(dataDir, "System.json")
	json, err := openJSON(systemFile)
	if err != nil {
		return err
	}

	err = json.Delete("hasEncryptedImages")
	if err != nil {
		log.Println(err)
	}

	err = json.Delete("hasEncryptedAudio")
	if err != nil {
		log.Println(err)
	}

	err = json.Delete("encryptionKey")
	if err != nil {
		log.Println(err)
	}

	jsonOutput := json.Bytes()
	err = writeContents(systemFile, &jsonOutput)
	if err != nil {
		log.Println(err)
	}

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
