package request

import "github.com/Hedonysym/httpfromtcp/internal/headers"

type errString string

func (e errString) Error() string {
	return string(e)
}

const InvalidHTTPVersionError = errString("invalid http version")

const InvalidHTTPMethodError = errString("invalid http method")

const InvalidRequestLineError = errString("invalid request line")

const IncompleteFileError = errString("file is incomplete")

const NoHeadersError = errString("no headers")

const BodyTooLongError = errString("body too long")

type Request struct {
	RequestLine RequestLine
	Headers     headers.Headers
	State       int
	Body        []byte
}

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

type chunkReader struct {
	data            string
	numBytesPerRead int
	pos             int
}
