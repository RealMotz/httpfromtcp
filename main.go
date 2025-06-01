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

	var line string
	for {
		var chunck [8]byte
		_, err := f.Read(chunck[:])

		if err != nil {
			if err == io.EOF {
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
			// Print completed line
			printLine(line)
			line = parts[i]
		}
	}

	if len(strings.TrimSpace(line)) > 0 {
		printLine(line)
	}
}

func printLine(line string) {
	fmt.Printf("read: %s\n", line)
}
