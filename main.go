package main

import (
	"fmt"
	"io"
	"os"
	"strings"
)

func main() {
	f, err := os.Open("messages.txt")
	if err != nil {
		fmt.Printf("%v\n", err)
		os.Exit(1)
	}
	defer f.Close()

	ch := getLinesChannel(f)

	for line := range ch {
		fmt.Printf("read: %s\n", line)
	}
}

func readChunks(f io.ReadCloser, ch chan<- string) {
	var line string
	for {
		var chunck [8]byte
		_, err := f.Read(chunck[:])

		if err != nil {
			if err == io.EOF {
				if len(strings.TrimSpace(line)) > 0 {
					ch <- line
				}
				close(ch)
				break
			}
			fmt.Printf("%v\n", err)
			os.Exit(1)
		}

		str := string(chunck[:])
		parts := strings.Split(str, "\n")

		// Add first part
		line += parts[0]

		for i := 1; i < len(parts); i++ {
			ch <- line
			line = parts[i]
		}
	}
}

func getLinesChannel(f io.ReadCloser) <-chan string {
	ch := make(chan string)
	go readChunks(f, ch)
	return ch
}
