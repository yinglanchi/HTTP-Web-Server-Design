package tritonhttp

import (
	"bytes"
	"testing"
)

func Test400And404cases(t *testing.T) {
	var tests = []struct {
		name string
		res  *Response
		want string
	}{
		{
			"400",
			&Response{
				StatusCode: 400,
				Proto:      "HTTP/1.1",
				StatusText: "Bad Request",
				Headers: map[string]string{
					"Date": "testWriteDate",
				},
			},
			"HTTP/1.1 400 Bad Request\r\n" +
				"Date: testWriteDate\r\n" +
				"\r\n",
		},
		{
			"404",
			&Response{
				StatusCode: 404,
				Proto:      "HTTP/1.1",
				StatusText: "Not Found",
				Headers: map[string]string{
					"Date": "testWriteDate",
				},
			},
			"HTTP/1.1 404 Not Found\r\n" +
				"Date: testWriteDate\r\n" +
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

func Test200Cases(t *testing.T) {
	var tests = []struct {
		name string
		res  *Response
		want string
	}{
		{
			"without body file - normal case",
			&Response{
				StatusCode: 200,
				Proto:      "HTTP/1.1",
				StatusText: "OK",
				Headers: map[string]string{
					"Connection":    "close",
					"Date":          "testWriteDate",
					"Last-Modified": "testWriteLastModified",
				},
			},
			"HTTP/1.1 200 OK\r\n" +
				"Connection: close\r\n" +
				"Date: testWriteDate\r\n" +
				"Last-Modified: testWriteLastModified\r\n" +
				"\r\n",
		},
		{
			"without body file - canonical string func test",
			&Response{
				StatusCode: 200,
				Proto:      "HTTP/1.1",
				StatusText: "OK",
				Headers: map[string]string{
					"Connection":    "close",
					"Date":          "testWriteDate",
					"last-modified": "testWriteLastModified",
				},
			},
			"HTTP/1.1 200 OK\r\n" +
				"Connection: close\r\n" +
				"Date: testWriteDate\r\n" +
				"Last-Modified: testWriteLastModified\r\n" +
				"\r\n",
		},
		{
			"with body file",
			&Response{
				StatusCode: 200,
				Proto:      "HTTP/1.1",
				StatusText: "OK",
				Headers: map[string]string{
					"Connection":    "close",
					"Date":          "testWriteDate",
					"last-modified": "testWriteLastModified",
				},
				FilePath: "testFiles/index.html",
			},
			"HTTP/1.1 200 OK\r\n" +
				"Connection: close\r\n" +
				"Date: testWriteDate\r\n" +
				"Last-Modified: testWriteLastModified\r\n" +
				"\r\n" +
				"test\n",
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
