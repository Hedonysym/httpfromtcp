package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/Hedonysym/httpfromtcp/internal/request"
	"github.com/Hedonysym/httpfromtcp/internal/response"
	"github.com/Hedonysym/httpfromtcp/internal/server"
)

const port = 42069

func main() {
	server, err := server.Serve(port, HTMLHandler)
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
