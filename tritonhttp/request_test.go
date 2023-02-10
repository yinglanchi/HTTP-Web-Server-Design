package tritonhttp

import (
	"bufio"
	"reflect"
	"strings"
	"testing"
)

func checkGoodRequest(t *testing.T, readErr error, reqGot, reqWant *Request) {
	if readErr != nil {
		t.Fatal(readErr)
	}
	if !reflect.DeepEqual(*reqGot, *reqWant) {
		t.Fatalf("\ngot: %v\nwant: %v", reqGot, reqWant)
	}
}

func checkBadRequest(t *testing.T, readErr error, reqGot *Request) {
	if readErr == nil {
		t.Errorf("\ngot unexpected request: %v\nwant: error", reqGot)
	}
}

func TestIncorrectFormat(t *testing.T) {
	var tests = []struct {
		name string
		req  string
	}{
		{
			"incorrect number of fields",
			"GET /testFiles/index.html HTTP/1.1 test \r\n" +
				"Host: host\r\n" +
				"\r\n",
		},
		{
			"incorrect method",
			"TEST /testFiles/index.html HTTP/1.1 \r\n",
		},
		{
			"non existing filepath",
			"GET /index.html HTTP/1.1 \r\n",
		},
		{
			"incorrect proto",
			"GET /testFiles/index.tml TEST \r\n",
		},
		{
			"incorrect header",
			"GET /testFiles/index.tml TEST \r\n" +
				"Host test\r\n" +
				"\r\n",
		},
		{
			"host not exist",
			"GET /testFiles/index.tml TEST \r\n" +
				"Connection: close\r\n" +
				"\r\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reqGot, _, err := ReadRequest(bufio.NewReader(strings.NewReader(tt.req)))
			checkBadRequest(t, err, reqGot)
		})
	}
}

func TestCombinedCalls(t *testing.T) {
	var tests = []struct {
		name     string
		reqText  string
		reqsWant []*Request
	}{
		{
			"Test4GoodGood",
			"GET /index.html HTTP/1.1\r\nHost: test\r\n\r\n" +
				"GET /index.html HTTP/1.1\r\nHost: test\r\n\r\n",
			[]*Request{
				{
					Method:  "GET",
					URL:     "/index.html",
					Proto:   "HTTP/1.1",
					Headers: map[string]string{},
					Host:    "test",
					Close:   false,
				},
				{
					Method:  "GET",
					URL:     "/index.html",
					Proto:   "HTTP/1.1",
					Headers: map[string]string{},
					Host:    "test",
					Close:   false,
				},
			},
		},
		{
			"Test4GoodBad",
			"GET /index.html HTTP/1.1\r\nHost: test\r\n\r\n" +
				"GETT /index.html HTTP/1.1\r\nHost: test\r\n\r\n",
			[]*Request{
				{
					Method:  "GET",
					URL:     "/index.html",
					Proto:   "HTTP/1.1",
					Headers: map[string]string{},
					Host:    "test",
					Close:   false,
				},
				nil,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			br := bufio.NewReader(strings.NewReader(tt.reqText))
			for _, reqWant := range tt.reqsWant {
				reqGot, _, err := ReadRequest(br)
				if reqWant != nil {
					checkGoodRequest(t, err, reqGot, reqWant)
				} else {
					checkBadRequest(t, err, reqGot)
				}
			}
		})
	}
}
