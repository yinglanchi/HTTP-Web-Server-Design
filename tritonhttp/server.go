package tritonhttp

import (
	"bufio"
	"fmt"
	"log"
	"mime"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

type Server struct {
	// Addr specifies the TCP address for the server to listen on,
	// in the form "host:port". It shall be passed to net.Listen()
	// during ListenAndServe().
	Addr string // e.g. ":0"

	// VirtualHosts contains a mapping from host name to the docRoot path
	// (i.e. the path to the directory to serve static files from) for
	// all virtual hosts that this server supports
	VirtualHosts map[string]string
}

// ListenAndServe listens on the TCP network address s.Addr and then
// handles requests on incoming connections.
func (s *Server) ListenAndServe() error {
	// Hint: Validate all docRoots
	log.Println("Validating all docRoots...")
	for _, docRoot := range s.VirtualHosts {
		f, err := os.Stat(docRoot)
		if os.IsNotExist(err) {
			log.Println("docRoot does not exist")
			return err
		}
		if !f.IsDir() {
			log.Println("docRoot is not a directory")
			return err
		}
	}
	log.Println("Finish validating all docRoots.")
	// Hint: create your listen socket and spawn off goroutines per incoming client
	log.Println("Start listening...")
	l, err := net.Listen("tcp", s.Addr)
	if err != nil {
		log.Println("listen error: ", err)
	}
	log.Println("finish listening.")
	for {
		conn, err := l.Accept()
		if err != nil {
			log.Println("accept error: ", err)
		}
		go s.handleConn(conn)
	}
}

func (s *Server) handleConn(conn net.Conn) {
	reader := bufio.NewReader(conn)
	for {
		// timeout 5 seconds
		conn.SetReadDeadline(time.Now().Add(5 * time.Second))

		req, readIn, err := ReadRequest(reader)

		if err, ok := err.(net.Error); ok && err.Timeout() {
			// if nothing read in, close the connection and return
			if !readIn {
				conn.Close()
				return
			}
			// else if only partial request is being processed, return 400 error
			//res := &Response{}
			res := s.handle400Requests(req)
			res.WriteResponse(conn)
			conn.Close()
			return
		}

		// if error exists, 400 error
		if err != nil {
			res := s.handle400Requests(req)
			res.WriteResponse(conn)
			conn.Close()
			return
		}

		// when no error exists, process the request
		// if host not in virtualHosts or if escape document root, 404 error
		hostName := req.Host
		docRoot := s.VirtualHosts[hostName]
		if docRoot == "" {
			res := s.handle404Requests(req)
			res.WriteResponse(conn)
			conn.Close()
			return
		}
		absolutePath := filepath.Join(docRoot, filepath.Clean(req.URL))
		if absolutePath[:len(docRoot)] != docRoot {
			res := s.handle404Requests(req)
			res.WriteResponse(conn)
			conn.Close()
			return
		}

		_, err = os.Stat(absolutePath)
		if err != nil {
			res := s.handle404Requests(req)
			res.WriteResponse(conn)
			conn.Close()
			return
		}

		// else, 200 ok
		res := s.handle200Requests(req)
		res.WriteResponse(conn)

		if req.Close {
			conn.Close()
			return
		}
	}
}

func (s *Server) handle400Requests(req *Request) (res *Response) {
	res = &Response{}
	res.Proto = "HTTP/1.1"
	res.StatusCode = 400
	res.StatusText = "Bad Request"
	res.Headers = make(map[string]string)
	res.Headers["Date"] = time.Now().Format("Tue, 19 Oct 2021 18:12:55 GMT")
	res.Headers["Connection"] = "close"
	res.Request = nil
	res.FilePath = ""
	return res
}

func (s *Server) handle200Requests(req *Request) (res *Response) {
	res = &Response{}
	res.Proto = "HTTP/1.1"
	res.StatusCode = 200
	res.StatusText = "OK"
	res.Headers = make(map[string]string)
	res.Headers["Date"] = time.Now().Format("Tue, 19 Oct 2021 18:12:55 GMT")
	absolutePath := filepath.Join(s.VirtualHosts[req.Host], filepath.Clean(req.URL))
	f, _ := os.Stat(absolutePath)
	res.Headers["Last-Modified"] = f.ModTime().Format("Tue, 19 Oct 2021 18:12:55 GMT")
	file, err := os.Open(absolutePath)
	if err != nil {
		log.Println(err)
	}
	res.Headers["Content-Type"] = mime.TypeByExtension(filepath.Ext(absolutePath))
	fmt.Println(res.Headers["Content-Type"])
	fmt.Println(res.Headers["Content-Type"])
	fi, _ := file.Stat()
	res.Headers["Content-Length"] = strconv.FormatInt(fi.Size(), 10)
	file.Close()
	if req.Close {
		res.Headers["Connection"] = "close"
	}
	res.Request = req
	res.FilePath = absolutePath
	return res
}

func (s *Server) handle404Requests(req *Request) (res *Response) {
	res = &Response{}
	res.Proto = "HTTP/1.1"
	res.StatusCode = 404
	res.StatusText = "Not Found"
	res.Headers = make(map[string]string)
	res.Headers["Date"] = time.Now().Format("Tue, 19 Oct 2021 18:12:55 GMT")
	res.Headers["Connection"] = "close"
	res.Request = req
	res.FilePath = ""
	return res
}

func GetFileContentType(ouput *os.File) (string, error) {

	// to sniff the content type only the first
	// 512 bytes are used.

	buf := make([]byte, 512)

	_, err := ouput.Read(buf)

	if err != nil {
		return "", err
	}

	// the function that actually does the trick
	contentType := http.DetectContentType(buf)

	return contentType, nil
}
