package tritonhttp

import (
	"io"
	"log"
	"os"
	"sort"
	"strconv"
	"strings"
)

type Response struct {
	Proto      string // e.g. "HTTP/1.1"
	StatusCode int    // e.g. 200
	StatusText string // e.g. "OK"

	// Headers stores all headers to write to the response.
	Headers map[string]string

	// Request is the valid request that leads to this response.
	// It could be nil for responses not resulting from a valid request.
	// Hint: you might need this to handle the "Connection: Close" requirement
	Request *Request

	// FilePath is the local path to the file to serve.
	// It could be "", which means there is no file to serve.
	FilePath string
}

// Write writes the res to the w.
func (res *Response) WriteResponse(w io.Writer) error {
	// write first line (i.e: request line)
	requestLine := res.Proto + " " + strconv.Itoa(res.StatusCode) + " " + res.StatusText + "\r\n"
	_, err := w.Write([]byte(requestLine))
	if err != nil {
		log.Println("write request line error: ", err)
		return err
	}

	// write headers
	// sort the keys when writing for the convenience when testing
	SortedKeys := make([]string, 0, len(res.Headers))

	for key := range res.Headers {
		SortedKeys = append(SortedKeys, key)
	}
	sort.Strings(SortedKeys)
	for i := 0; i < len(SortedKeys); i++ {
		key := SortedKeys[i]
		// convert key into canonical format
		keyArray := strings.Split(key, "-")
		CanonicalKey := ""
		for i := 0; i < len(keyArray); i++ {
			keyParts := keyArray[i]
			CanonicalKey += strings.ToUpper(string(keyParts[0])) + keyParts[1:]
			if i != len(keyArray)-1 {
				CanonicalKey += "-"
			}
		}
		header := CanonicalKey + ": " + res.Headers[key] + "\r\n"
		_, err := w.Write([]byte(header))
		if err != nil {
			log.Println("write header error: ", err)
			return err
		}
	}
	_, err = w.Write([]byte("\r\n"))
	if err != nil {
		log.Println("write CRLF error: ", err)
		return err
	}

	// write body (might not exist)
	if res.FilePath != "" {
		file, err := os.ReadFile(res.FilePath)
		if err != nil {
			log.Println("read file error: ", err)
		}
		_, err = w.Write(file)
		if err != nil {
			log.Println("write body file error: ", err)
			return err
		}
	}
	return nil
}
