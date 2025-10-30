package response

import (
	"io"
)

type StatusCode int

const StatusOk = StatusCode(200)

const StatusBadRequest = StatusCode(400)

const StatusInternalServerError = StatusCode(500)

type Writer struct {
	Writer io.Writer
	State  int
}

func NewWriter(w io.Writer) *Writer {
	return &Writer{
		Writer: w,
		State:  1,
	}
}
