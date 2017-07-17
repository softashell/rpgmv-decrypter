package main

import (
	"fmt"
	"log"
	"strconv"
	"time"
)

const (
	headerLen = 16
	signature = "5250474d56000000"
	version   = "000301"
	remain    = "0000000000"
)

func decryptFile(filePath, outPath string, key []byte) error {
	start := time.Now()

	content, err := readFileContents(filePath)
	if err != nil {
		log.Fatal(err)
	}

	if len(content) < headerLen*2 {
		return fmt.Errorf("file is too small")
	}

	if !checkFakeHeader(content[:headerLen]) {
		return fmt.Errorf("invalid header")
	}

	content = content[headerLen:]

	if len(content) < 1 {
		return fmt.Errorf("file without header is too small")
	}

	for i := 0; i < headerLen; i++ {
		content[i] = content[i] ^ key[i]
	}

	err = writeFileContents(outPath, &content)
	if err != nil {
		return err
	}

	fmt.Printf("decrypted in %s\n", time.Since(start))

	return nil
}

func checkFakeHeader(header []byte) bool {
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
