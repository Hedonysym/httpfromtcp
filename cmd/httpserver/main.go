package main

import (
	"crypto/sha256"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/Hedonysym/httpfromtcp/internal/headers"
	"github.com/Hedonysym/httpfromtcp/internal/request"
	"github.com/Hedonysym/httpfromtcp/internal/response"
	"github.com/Hedonysym/httpfromtcp/internal/server"
)

const port = 42069

func main() {
	server, err := server.Serve(port, videoHandler)
	if err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
	defer server.Close()
	log.Println("Server started on port", port)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan
	log.Println("Server gracefully stopped")
}

func HTMLHandler(w *response.Writer, req *request.Request) {
	status := response.StatusOk
	body := []byte{}

	switch req.RequestLine.RequestTarget {
	case "/yourproblem":
		status = response.StatusBadRequest
		body = []byte(`<html>
  <head>
    <title>400 Bad Request</title>
  </head>
  <body>
    <h1>Bad Request</h1>
    <p>Your request honestly kinda sucked.</p>
  </body>
</html>`)
	case "/myproblem":
		status = response.StatusInternalServerError
		body = []byte(`<html>
  <head>
    <title>500 Internal Server Error</title>
  </head>
  <body>
    <h1>Internal Server Error</h1>
    <p>Okay, you know what? This one is on me.</p>
  </body>
</html>`)
	default:
		body = []byte(`<html>
  <head>
    <title>200 OK</title>
  </head>
  <body>
    <h1>Success!</h1>
    <p>Your request was an absolute banger.</p>
  </body>
</html>`)
	}

	headers := response.GetDefaultHeaders(len(body))
	headers.Override("content-type", "text/html")
	err := w.WriteStatusLine(status)
	if err != nil {
		log.Println(err)
		return
	}
	err = w.WriteHeaders(headers)
	if err != nil {
		log.Println(err)
		return
	}
	bytesWritten, err := w.WriteBody(body)
	if err != nil || bytesWritten != len(body) {
		log.Println(err)
		return
	}
}

func chunkedHandler(w *response.Writer, req *request.Request) {
	if !strings.HasPrefix(req.RequestLine.RequestTarget, "/httpbin") {
		err := w.WriteStatusLine(response.StatusBadRequest)
		if err != nil {
			log.Println(err)
		}
		return
	}
	url := "https://httpbin.org" + strings.TrimPrefix(req.RequestLine.RequestTarget, "/httpbin")
	resp, err := http.Get(url)
	if err != nil {
		err = w.WriteStatusLine(response.StatusInternalServerError)
		if err != nil {
			log.Println(err)
		}
		return
	}
	defer resp.Body.Close()

	err = w.WriteStatusLine(response.StatusOk)
	if err != nil {
		log.Println(err)
		return
	}
	header := response.GetDefaultHeaders(0)
	header.Override("transfer-encoding", "chunked")
	header["trailer"] = "X-Content-SHA256, X-Content-Length"
	delete(header, "content-length")
	err = w.WriteHeaders(header)
	if err != nil {
		log.Println(err)
		return
	}
	buf := make([]byte, 1024)
	bodySize := 0
	body := []byte{}
	for {
		n, err := resp.Body.Read(buf)
		if err != nil {
			break
		}
		log.Printf("%v bytes read\n", n)
		_, err = w.Writer.Write([]byte(fmt.Sprintf("%x\r\n", n)))
		if err != nil {
			log.Println(err)
			return
		}
		body = append(body, copyBytes(buf[:n], n)...)
		i, err := w.WriteChunkedBody(buf[:n])
		if err != nil {
			log.Println(err)
			return
		}
		bodySize += i
		_, err = w.Writer.Write([]byte("\r\n"))
		if err != nil {
			log.Println(err)
			return
		}
	}
	_, err = w.WriteChunkedBodyDone()
	if err != nil {
		log.Println(err)
		return
	}
	sha := sha256.Sum256(body)
	trail := headers.NewHeaders()
	trail["X-Content-SHA256"] = fmt.Sprintf("%x", sha)
	trail["X-Content-Length"] = fmt.Sprintf("%d", bodySize)
	err = w.WriteTrailers(trail)
	if err != nil {
		log.Println(err)
		return
	}
}

func copyBytes(src []byte, n int) []byte {
	dst := make([]byte, n)
	copy(dst, src)
	return dst
}

func videoHandler(w *response.Writer, req *request.Request) {
	if !strings.HasPrefix(req.RequestLine.RequestTarget, "/video") {
		err := w.WriteStatusLine(response.StatusBadRequest)
		if err != nil {
			log.Println(err)
		}
		return
	}
	vid, err := os.ReadFile("/home/Big_Man/workspace/github.com/Hedonysym/httpfromtcp/assets/vim.mp4")
	if err != nil {
		err = w.WriteStatusLine(response.StatusInternalServerError)
		if err != nil {
			log.Println(err)
		}
		return
	}
	err = w.WriteStatusLine(response.StatusOk)
	if err != nil {
		log.Println(err)
		return
	}
	header := response.GetDefaultHeaders(len(vid))
	header.Override("content-type", "video/mp4")
	err = w.WriteHeaders(header)
	if err != nil {
		log.Println(err)
		return
	}
	bytesWritten, err := w.WriteBody(vid)
	if err != nil || bytesWritten != len(vid) {
		log.Println(err)
		return
	}
}
