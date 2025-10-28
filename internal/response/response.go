package response

import (
	"fmt"
	"io"

	"github.com/Hedonysym/httpfromtcp/internal/headers"
)

func WriteStatus(w io.Writer, statusCode StatusCode) error {
	switch statusCode {
	case StatusOk:
		_, err := w.Write([]byte("HTTP/1.1 200 OK\r\n"))
		return err
	case StatusBadRequest:
		_, err := w.Write([]byte("HTTP/1.1 400 Bad Request\r\n"))
		return err
	case StatusInternalServerError:
		_, err := w.Write([]byte("HTTP/1.1 500 Internal Server Error\r\n"))
		return err
	default:
		defaultMsg := fmt.Sprintf("HTTP/1.1 %d \r\n", statusCode)
		_, err := w.Write([]byte(defaultMsg))
		return err
	}
}

func GetDefaultHeaders(contentLen int) headers.Headers {
	h := headers.Headers{
		"Content-Length": fmt.Sprintf("%d", contentLen),
		"Connection":     "close",
		"Content-Type":   "text/plain",
	}
	return h
}

func WriteHeaders(w io.Writer, headers headers.Headers) error {
	for k, v := range headers {
		_, err := w.Write([]byte(fmt.Sprintf("%s: %s\r\n", k, v)))
		if err != nil {
			return err
		}
	}
	_, err := w.Write([]byte("\r\n"))
	return err
}
