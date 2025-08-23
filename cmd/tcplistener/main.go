package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"strings"
)

func main() {
	listener, err := net.Listen("tcp", ":42069")
	if err != nil {
		fmt.Printf("%v\n", err)
		os.Exit(1)
	}
	defer listener.Close()

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println("Connection accepted")

		ch := getLinesChannel(conn)

		for line := range ch {
			fmt.Printf("read: %s\n", line)
		}

		fmt.Println("Connection closed")
	}
}

func getLinesChannel(f io.ReadCloser) <-chan string {
	ch := make(chan string)
	go readChunks(f, ch)
	return ch
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
