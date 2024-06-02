package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

// usage ./unifi-restore dump.json
// AVDANCED: set the buffer size (default 5 GiB) and max token (default 1 GiB) size which is needed for large numbers of devices
// BUFFER_SIZE_GB=5 TOKEN_SIZE_MB=1 ./unifi-restore dump.json

func main() {
	args := os.Args[1:]
	path := args[0]
	file, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	bufferSizeGB := 5
	maxTokenSizeMB := 1

	bufferSizeEnv := os.Getenv("BUFFER_SIZE_GB")
	if bufferSizeEnv != "" {
		bufferSizeGB, _ = strconv.Atoi(bufferSizeEnv)
	}
	maxTokenSizeEnv := os.Getenv("TOKEN_SIZE_MB")
	if maxTokenSizeEnv != "" {
		maxTokenSizeMB, _ = strconv.Atoi(maxTokenSizeEnv)
	}

	scanner := bufio.NewScanner(file)
	buf := make([]byte, 0, bufferSizeGB*1024*1024*1024)
	scanner.Buffer(buf, maxTokenSizeMB*1024*1024)

	collection := ""
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, `{"__cmd":"select","collection":"`) {
			s := strings.TrimPrefix(line, `{"__cmd":"select","collection":"`)
			s = strings.TrimSuffix(s, `"}`)
			collection = s
			fmt.Printf("running updates for %s\n", collection)
		}
		fmt.Printf("echo '%s' | mongoimport -h localhost:27117 --db ace --collection %s\n", line, collection)

		cmd := exec.Command("bash", "-c", fmt.Sprintf("echo '%s' | mongoimport -h localhost:27117 --db ace --collection %s", line, collection))
		out, _ := cmd.CombinedOutput()
		fmt.Printf("%s", out)
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}
