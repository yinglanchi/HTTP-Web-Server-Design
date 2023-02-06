package tritonhttp

import (
	"testing"
)

const (
	contentTypeHTML = "text/html; charset=utf-8"
	contentTypeJPG  = "image/jpeg"
	contentTypePNG  = "image/png"
)

func TestHandleGoodRequest(t *testing.T) {
	var tests = []struct {
		name             string
		req              *Request
		statusWant       int
		headersWant      []string
		headerValuesWant map[string]string
		filePathWant     string // relative to doc root
	}{
		{
			"OKBasic",
			&Request{
				Method:  "GET",
				URL:     "/index.html",
				Proto:   "HTTP/1.1",
				Headers: map[string]string{},
				Host:    "website1",
				Close:   false,
			},
			200,
			[]string{
				"Date",
				"Last-Modified",
			},
			map[string]string{
				"Content-Type":   contentTypeHTML,
				"Content-Length": "377",
			},
			"../docroot_dirs/htdocs1/index.html",
		},
	}
	virtualHosts := ParseVHConfigFile("../virtual_hosts.yaml", "../docroot_dirs")

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Server{
				Addr:         ":0",
				VirtualHosts: virtualHosts,
			}
			res := s.handle200Requests(tt.req)
			if res.StatusCode != tt.statusWant {
				t.Fatalf("status code got: %v, want: %v", res.StatusCode, tt.statusWant)
			}
			for _, h := range tt.headersWant {
				if _, ok := res.Headers[h]; !ok {
					t.Fatalf("missing header %q", h)
				}
			}
			for h, vWant := range tt.headerValuesWant {
				v, ok := res.Headers[h]
				if !ok {
					t.Fatalf("missing header %q", h)
				}
				if v != vWant {
					t.Fatalf("header %q value got: %q, want %q", h, v, vWant)
				}
			}
			if tt.filePathWant != "" {
				// Case with file to serve
				if res.FilePath != tt.filePathWant {
					t.Fatalf("file path (relative to testdata/) got: %q, want: %q", res.FilePath, tt.filePathWant)
				}
			} else {
				// Case with no file to serve
				if res.FilePath != tt.filePathWant {
					t.Fatalf("file path got: %q, want: %q", res.FilePath, tt.filePathWant)
				}
			}
		})
	}
}
