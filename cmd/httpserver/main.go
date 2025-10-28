package main

import (
	"io"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/Hedonysym/httpfromtcp/internal/request"
	"github.com/Hedonysym/httpfromtcp/internal/server"
)

const port = 42069

func main() {
	server, err := server.Serve(port, TestHandler)
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

func TestHandler(w io.Writer, req *request.Request) *server.HandlerError {
	switch req.RequestLine.RequestTarget {
	case "/yourproblem":
		return &server.HandlerError{
			StatusCode: 400,
			Err:        "Your problem is not my problem\n",
		}
	case "/myproblem":
		return &server.HandlerError{
			StatusCode: 500,
			Err:        "Woopsie, my bad\n",
		}
	default:
		_, err := w.Write([]byte("All good, frfr\n"))
		if err != nil {
			return &server.HandlerError{
				StatusCode: 500,
				Err:        err.Error(),
			}
		}
		return nil
	}
}
