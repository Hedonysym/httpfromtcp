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
		"content-length": fmt.Sprintf("%d", contentLen),
		"connection":     "close",
		"content-type":   "text/plain",
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

func (w *Writer) WriteStatusLine(statusCode StatusCode) error {
	if w.State != 1 {
		return fmt.Errorf("invalid state: %d", w.State)
	}
	w.State = 2
	return WriteStatus(w.Writer, statusCode)
}

func (w *Writer) WriteHeaders(headers headers.Headers) error {
	if w.State != 2 {
		return fmt.Errorf("invalid state: %d", w.State)
	}
	w.State = 3
	return WriteHeaders(w.Writer, headers)
}

func (w *Writer) WriteBody(body []byte) (int, error) {
	if w.State != 3 {
		return 0, fmt.Errorf("invalid state: %d", w.State)
	}
	w.State = 0
	numBytes, err := w.Writer.Write(body)
	return numBytes, err
}

func (w *Writer) WriteChunkedBody(body []byte) (int, error) {
	if w.State != 3 {
		return 0, fmt.Errorf("invalid state: %d", w.State)
	}
	numBytes, err := w.Writer.Write(body)
	return numBytes, err
}

func (w *Writer) WriteChunkedBodyDone() (int, error) {
	if w.State != 3 {
		return 0, fmt.Errorf("invalid state: %d", w.State)
	}
	w.State = 0
	numBytes, err := w.Writer.Write([]byte("0\r\n"))
	return numBytes, err
}

func (w *Writer) WriteTrailers(h headers.Headers) error {
	if w.State != 0 {
		return fmt.Errorf("invalid state: %d", w.State)
	}
	err := WriteHeaders(w.Writer, h)
	if err != nil {
		return err
	}
	return nil
}
