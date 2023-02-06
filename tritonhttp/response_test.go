package tritonhttp

import (
	"bytes"
	"testing"
)

// normal case: check whether the function works in general
func TestWriteResponse1(t *testing.T) {
	var tests = []struct {
		name string
		res  *Response
		want string
	}{
		{
			"OK",
			&Response{
				StatusCode: 200,
				Proto:      "HTTP/1.1",
				StatusText: "OK",
				Headers: map[string]string{
					"Connection": "close",
					"Date":       "foobar",
					"Misc":       "hello world",
				},
			},
			"HTTP/1.1 200 OK\r\n" +
				"Connection: close\r\n" +
				"Date: foobar\r\n" +
				"Misc: hello world\r\n" +
				"\r\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buffer bytes.Buffer
			if err := tt.res.WriteResponse(&buffer); err != nil {
				t.Fatal(err)
			}
			got := buffer.String()
			if got != tt.want {
				t.Fatalf("got: %q, want: %q", got, tt.want)
			}
		})
	}
}

// test for canonical conversion of headers
func TestWriteResponse2(t *testing.T) {
	var tests = []struct {
		name string
		res  *Response
		want string
	}{
		{
			"OK",
			&Response{
				StatusCode: 200,
				Proto:      "HTTP/1.1",
				StatusText: "OK",
				Headers: map[string]string{
					"Connection":    "close",
					"date-modified": "foobar",
					"misc":          "hello world",
				},
			},
			"HTTP/1.1 200 OK\r\n" +
				"Connection: close\r\n" +
				"Date-Modified: foobar\r\n" +
				"Misc: hello world\r\n" +
				"\r\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buffer bytes.Buffer
			if err := tt.res.WriteResponse(&buffer); err != nil {
				t.Fatal(err)
			}
			got := buffer.String()
			if got != tt.want {
				t.Fatalf("got: %q, want: %q", got, tt.want)
			}
		})
	}
}

// test with body file
func TestWriteResponse3(t *testing.T) {
	var tests = []struct {
		name string
		res  *Response
		want string
	}{
		{
			"OK",
			&Response{
				StatusCode: 200,
				Proto:      "HTTP/1.1",
				StatusText: "OK",
				Headers: map[string]string{
					"Connection":    "close",
					"date-modified": "foobar",
					"misc":          "hello world",
				},
				FilePath: "index.html",
			},
			"HTTP/1.1 200 OK\r\n" +
				"Connection: close\r\n" +
				"Date-Modified: foobar\r\n" +
				"Misc: hello world\r\n" +
				"\r\n" +
				"<html>\n\n<head>\n    <title>Basic index file for website 2</title>\n</head>\n\n<body>\n    <h1>This is a basic index file for website 2</h1>\n    You can use this for testing.\n</body>\n\n</html>\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buffer bytes.Buffer
			if err := tt.res.WriteResponse(&buffer); err != nil {
				t.Fatal(err)
			}
			got := buffer.String()
			if got != tt.want {
				t.Fatalf("got: %q, want: %q", got, tt.want)
			}
		})
	}
}

// test for empty value
func TestWriteResponse4(t *testing.T) {
	var tests = []struct {
		name string
		res  *Response
		want string
	}{
		{
			"OK",
			&Response{
				StatusCode: 200,
				Proto:      "HTTP/1.1",
				StatusText: "OK",
				Headers: map[string]string{
					"Connection": "close",
					"Date":       "foobar",
					"Misc":       "",
				},
			},
			"HTTP/1.1 200 OK\r\n" +
				"Connection: close\r\n" +
				"Date: foobar\r\n" +
				"Misc: \r\n" +
				"\r\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buffer bytes.Buffer
			if err := tt.res.WriteResponse(&buffer); err != nil {
				t.Fatal(err)
			}
			got := buffer.String()
			if got != tt.want {
				t.Fatalf("got: %q, want: %q", got, tt.want)
			}
		})
	}
}
