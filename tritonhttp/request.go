package tritonhttp

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"strings"
)

type Request struct {
	Method string // e.g. "GET"
	URL    string // e.g. "/path/to/a/file"
	Proto  string // e.g. "HTTP/1.1"

	// Headers stores the key-value HTTP headers
	Headers map[string]string

	Host  string // determine from the "Host" header
	Close bool   // determine from the "Connection" header
}

func ReadRequest(reader *bufio.Reader) (req *Request, readIn bool, err error) {
	req = &Request{}
	req.Headers = make(map[string]string)

	// read initial request line
	request, _, err := reader.ReadLine()
	if err != nil {
		log.Println("read request line error: ", err)
		return nil, false, err
	}

	requestFields := strings.Split(string(request), " ")
	// check for incorrect request line formats
	// if format incorrect, return 400 error
	if len(requestFields) != 3 {
		log.Println("incorrect request line format")
		req.Method = ""
		req.URL = ""
		req.Proto = ""
		return nil, true, fmt.Errorf("400")
	}

	req.Method = requestFields[0]
	req.URL = requestFields[1]
	req.Proto = requestFields[2]

	if req.Method != "GET" {
		return nil, true, fmt.Errorf("400")
	}

	if req.URL[0] != '/' {
		return nil, true, fmt.Errorf("400")
	}

	if req.Proto != "HTTP/1.1" {
		return nil, true, fmt.Errorf("400")
	}

	if req.URL == "/" {
		req.URL = "/index.html"
	}

	// start reading in body of the request file
	hostExist := false
	req.Close = false
	for {
		line, _, err := reader.ReadLine()
		if err != nil {
			log.Println("read line error: ", err)
			return nil, true, fmt.Errorf("400")
		}
		// reached the end
		if len(line) == 0 || err == io.EOF {
			break
		}
		lineFields := strings.Split(string(line), ": ")
		if len(lineFields) != 2 {
			return nil, true, fmt.Errorf("400")
		}
		key := lineFields[0]
		value := lineFields[1]

		if key == "Host" {
			hostExist = true
			req.Host = value
		} else if key == "Connection" && value == "close" {
			req.Close = true
		} else {
			req.Headers[key] = value
		}
	}

	if !hostExist {
		return nil, true, fmt.Errorf("400")
	}

	return req, true, nil
}
