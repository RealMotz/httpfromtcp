package request

import (
	"bytes"
	"fmt"
	"io"
	"strings"
)

type ReqStatus int

const (
	Initiliazed ReqStatus = iota
	Done
)

type Request struct {
	RequestLine RequestLine
	State       ReqStatus
}

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

func RequestFromReader(reader io.Reader) (*Request, error) {
	request := Request{
		State: Initiliazed,
	}

	var readToIndex int
	var bufferSize int = 8
	var buffer []byte = make([]byte, bufferSize)
	for request.State == Initiliazed {
		readBytes, err := reader.Read(buffer[readToIndex:])
		if err != nil {
			if err == io.EOF {
				request.State = Done
				break
			}
			return nil, err
		}

		readToIndex += readBytes

		parsedBytes, err := request.parse(buffer[:readToIndex])
		if err != nil {
			return nil, err
		}

		copy(buffer, buffer[parsedBytes:])
		readToIndex -= parsedBytes

		// If buffer is full
		if readToIndex >= len(buffer) {
			newBuffer := make([]byte, len(buffer)*2)
			copy(newBuffer, buffer)
			buffer = newBuffer
		}
	}

	return &request, nil
}

func (r *Request) parse(data []byte) (int, error) {
	switch r.State {
	case Initiliazed:
		line, bytesRead, err := parseRequestLine(data)
		if err != nil {
			return 0, err
		}

		if bytesRead == 0 {
			// need more data
			return 0, nil
		}

		r.RequestLine = *line
		r.State = Done
		return bytesRead, nil
	case Done:
		return 0, fmt.Errorf("error: trying to read data in Done state")
	default:
		return 0, fmt.Errorf("error: unknown state")
	}
}

func parseRequestLine(request []byte) (*RequestLine, int, error) {
	idx := bytes.Index(request, []byte("\r\n"))
	if idx == -1 {
		return nil, 0, nil
	}

	requestStr := string(request[:idx])
	// Example: GET /coffee HTTP/1.1 -> [GET /Coffee HTTP/1.1]
	parts := strings.Split(requestStr, " ")
	if len(parts) != 3 {
		return nil, 0, fmt.Errorf("poorly formatted request-line: %s", requestStr)
	}

	// Verify that the method only contains alphanumeric chars
	method := parts[0]
	for _, char := range method {
		if char < 65 || char > 122 || (char > 90 && char < 97) {
			return nil, len(request), fmt.Errorf("invalid method")
		}
	}

	target := parts[1]

	// Get the start line parts
	versionParts := strings.Split(parts[2], "/")
	if len(versionParts) != 2 {
		return nil, len(request), fmt.Errorf("malformed start-line: %s", requestStr)
	}

	// verify that http is present
	httpPart := versionParts[0]
	if httpPart != "HTTP" {
		return nil, len(request), fmt.Errorf("unrecognized HTTP-version: %s", httpPart)
	}

	// Verify that the http version is 1.1
	version := versionParts[1]
	if version != "1.1" {
		return nil, len(request), fmt.Errorf("unrecognized HTTP-version: %s", version)
	}

	r := RequestLine{
		HttpVersion:   version,
		RequestTarget: target,
		Method:        method,
	}

	return &r, len(request), nil
}
