package request

type errString string

func (e errString) Error() string {
	return string(e)
}

const InvalidHTTPVersionError = errString("invalid http version")

const InvalidHTTPMethodError = errString("invalid http method")

const InvalidRequestLineError = errString("invalid request line")

type Request struct {
	RequestLine RequestLine
	State       int
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
