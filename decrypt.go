package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"time"
)

const (
	headerLen = 16
	signature = "5250474d56000000"
	version   = "000301"
	remain    = "0000000000"
)

func calculateKey(decryptionKey string) ([]string, error) {
	var i int
	var chunk string
	var key []string

	if len(decryptionKey)/2 != 16 {
		return nil, fmt.Errorf("invalid key provided")
	}

	for _, n := range decryptionKey {
		if i == 2 {
			key = append(key, chunk)

			i = 0
			chunk = ""
		}

		i++
		chunk += string(n)
	}

	if i == 2 {
		key = append(key, chunk)
	}

	return key, nil
}

func decryptFile(filePath, outPath string, key []string) error {
	start := time.Now()

	content, err := getContents(filePath)
	if err != nil {
		log.Fatal(err)
	}

	if len(content) < headerLen*2 {
		return fmt.Errorf("file is too small")
	}

	if !checkFakeHeader(&content) {
		return fmt.Errorf("invalid header")
	}

	content = content[headerLen:]

	if len(content) < 1 {
		return fmt.Errorf("file without header is too small")
	}

	for i := 0; i < headerLen; i++ {
		num, err := strconv.ParseInt(key[i], 16, 32)
		if err != nil {
			log.Fatal(err)
		}

		content[i] = content[i] ^ byte(num)
	}

	err = writeContents(outPath, &content)
	if err != nil {
		return err
	}

	fmt.Printf("decrypted in %s\n", time.Since(start))

	return nil
}

func getContents(filePath string) ([]byte, error) {
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

func writeContents(filePath string, content *[]byte) error {
	f, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = f.Write(*content)
	if err != nil {
		return err
	}

	f.Sync()

	return nil
}

func getByteArray(byteArray []byte, startPos, length int) []byte {
	// Don't allow start-values below 0
	if startPos < 0 {
		startPos = 0
	}

	// Check if length is to below 0 (to end of array)
	if length < 0 {
		length = len(byteArray) - startPos
	}

	newByteArray := make([]byte, length)
	n := 0

	for i := startPos; i < (startPos + length); i++ {
		// Check if byte array is on the last pos and return shorter byte array if
		if len(byteArray) <= i {
			return getByteArray(newByteArray, 0, n)
		}

		newByteArray[n] = byteArray[i]
		n++
	}

	return newByteArray
}

func checkFakeHeader(content *[]byte) bool {
	header := getByteArray(*content, 0, headerLen)
	refBytes := make([]byte, headerLen)
	refStr := signature + version + remain

	// Generate reference bytes
	for i := 0; i < headerLen; i++ {
		subStrStart := i * 2
		num, err := strconv.ParseInt(refStr[subStrStart:subStrStart+2], 16, 32)
		if err != nil {
			log.Fatal(err)
		}

		refBytes[i] = byte(num)
	}

	// Verify header (Check if its an encrypted file)
	for i := 0; i < headerLen; i++ {
		if refBytes[i] != header[i] {
			return false
		}
	}

	return true
}
