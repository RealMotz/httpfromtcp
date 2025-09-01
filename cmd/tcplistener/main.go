package main

import (
	"fmt"
	"log"
	"net"
	"os"

	"github.com/RealMotz/httpfromtcp/internal/request"
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

		req, err := request.RequestFromReader(conn)
		if err != nil {
			fmt.Printf("%v\n", err)
			os.Exit(1)
		}

		line := req.RequestLine
		fmt.Println("Request line:")
		fmt.Printf("- Method: %s\n", line.Method)
		fmt.Printf("- Target: %s\n", line.RequestTarget)
		fmt.Printf("- Version: %s\n", line.HttpVersion)

		fmt.Println("Connection closed")
	}
}
