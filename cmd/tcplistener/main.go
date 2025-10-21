package main

import (
	"fmt"
	"log"
	"net"

	"github.com/Hedonysym/httpfromtcp/internal/request"
)

func main() {
	listener, err := net.Listen("tcp", ":42069")
	if err != nil {
		log.Fatal(err)
	}
	defer listener.Close()

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("Accepted connection from", conn.RemoteAddr())

		req, err := request.RequestFromReader(conn)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf(
			"Request line:\n- Method: %s\n- Target: %s\n- Version: %s\n",
			req.RequestLine.Method,
			req.RequestLine.RequestTarget,
			req.RequestLine.HttpVersion,
		)
		fmt.Printf("Headers:\n")
		for key, value := range req.Headers {
			fmt.Printf("- %s: %s\n", key, value)
		}

		fmt.Println("Connection closed")
		conn.Close()
	}
}

/*
func getLinesChannel(f io.ReadCloser) <-chan string {
	ch := make(chan string)
	go func() {
		completeLine := ""
		for {
			bytes := make([]byte, 8)
			_, err := f.Read(bytes)
			if err == io.EOF {
				ch <- completeLine
				break
			}
			if err != nil {
				log.Fatal(err)
			}
			split := strings.Split(string(bytes), "\n")
			if completeLine == "" {
				completeLine = split[0]
			} else {
				completeLine += split[0]
			}
			if len(split) > 1 {
				ch <- completeLine
				completeLine = split[1]
			}

		}
		close(ch)
	}()
	return ch
}
*/
