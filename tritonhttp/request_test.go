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

// normal good request, checking the general funcationality of ReadRequest
// one test with connection not closed, the other one with connection closed
func TestReadRequest1(t *testing.T) {
	var tests = []struct {
		name    string
		reqText string
		reqWant *Request
	}{
		{
			"Test1NotCloseConnection",
			"GET /index.html HTTP/1.1\r\n" +
				"Host: test\r\n" +
				"\r\n",
			&Request{
				Method:  "GET",
				URL:     "/index.html",
				Proto:   "HTTP/1.1",
				Headers: map[string]string{},
				Host:    "test",
				Close:   false,
			},
		},
		{
			"Test1CloseConnection",
			"GET /index.html HTTP/1.1\r\n" +
				"Host: test\r\n" +
				"Connection: close\r\n" +
				"\r\n",
			&Request{
				Method:  "GET",
				URL:     "/index.html",
				Proto:   "HTTP/1.1",
				Headers: map[string]string{},
				Host:    "test",
				Close:   true,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reqGot, _, err := ReadRequest(bufio.NewReader(strings.NewReader(tt.reqText)))
			checkGoodRequest(t, err, reqGot, tt.reqWant)
		})
	}

}

// test the case with inccorect request line format
func TestReadRequest2(t *testing.T) {
	var tests = []struct {
		name string
		req  string
	}{
		{
			"Test2",
			"This is a bad request\r\n",
		},
		{
			"Empty",
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

// test the case without host specified
func TestReadRequest3(t *testing.T) {
	var tests = []struct {
		name    string
		reqText string
		reqWant *Request
	}{
		{
			"Test3",
			"GET /index.html HTTP/1.1\r\n" +
				"Connection: close\r\n" +
				"\r\n",
			&Request{
				Method:  "GET",
				URL:     "/index.html",
				Proto:   "HTTP/1.1",
				Headers: map[string]string{},
				Host:    "",
				Close:   true,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reqGot, _, err := ReadRequest(bufio.NewReader(strings.NewReader(tt.reqText)))
			checkBadRequest(t, err, reqGot)
		})
	}

}

func TestReadRequest4(t *testing.T) {
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
