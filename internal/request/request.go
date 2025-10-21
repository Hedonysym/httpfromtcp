package request

import (
	"fmt"
	"io"
	"strings"

	"github.com/Hedonysym/httpfromtcp/internal/headers"
)

const bufferSize = 8

func RequestFromReader(r io.Reader) (*Request, error) {
	buf := make([]byte, bufferSize)
	readToIndex := 0
	req := &Request{
		State:   1,
		Headers: headers.NewHeaders(),
	}
	for req.State != 0 {
		if len(buf)-readToIndex < bufferSize {
			nb := make([]byte, len(buf)*2)
			copy(nb, buf[:readToIndex])
			buf = nb
		}

		if readToIndex == len(buf) {
			nb := make([]byte, len(buf)*2)
			copy(nb, buf[:readToIndex])
			buf = nb
		}

		if readToIndex == len(buf) {
			nb := make([]byte, len(buf)*2)
			copy(nb, buf[:readToIndex])
			buf = nb
		}
		n, readErr := r.Read(buf[readToIndex:])
		if n > 0 {
			readToIndex += n
		}
		bytesParsed, parseErr := req.parse(buf[:readToIndex])
		if parseErr != nil {
			return nil, parseErr
		}
		if bytesParsed > 0 {
			copy(buf, buf[bytesParsed:readToIndex])
			readToIndex -= bytesParsed
		}

		if req.State == 0 {
			break
		}

		if readErr == io.EOF {
			return nil, IncompleteFileError
		}
	}
	if req.RequestLine.HttpVersion != "1.1" {
		return nil, fmt.Errorf("HTTP version %s not supported", req.RequestLine.HttpVersion)
	}
	if req.RequestLine.Method != "GET" && req.RequestLine.Method != "POST" {
		return nil, InvalidHTTPMethodError
	}
	if len(req.Headers) == 0 {
		return nil, NoHeadersError
	}
	return req, nil
}

func parseRequestLine(s string) (RequestLine, int, error) {
	i := strings.Index(s, "\r\n")
	if i == -1 {
		return RequestLine{}, 0, nil
	}
	line := s[:i]
	fields := strings.Fields(line)
	if len(fields) != 3 {
		return RequestLine{}, 0, InvalidRequestLineError
	}
	rl := RequestLine{
		Method:        strings.ToUpper(fields[0]),
		RequestTarget: fields[1],
		HttpVersion:   strings.TrimSpace(strings.Replace(fields[2], "HTTP/", "", 1)),
	}
	return rl, i + 2, nil
}

func (r *Request) parse(data []byte) (int, error) {
	total := 0
	for r.State != 0 {
		n, err := r.parseSingle(data[total:])
		if err != nil {
			return total, err
		}
		if n == 0 {
			// no progress this call, let caller read more
			break
		}
		total += n
	}
	return total, nil
}

func (r *Request) parseSingle(data []byte) (int, error) {
	if r.State == 1 {
		reqLine, bytesParsed, err := parseRequestLine(string(data))
		if err != nil {
			return 0, err
		}
		if bytesParsed == 0 {
			return 0, nil
		}

		r.RequestLine = reqLine
		r.State = 2
		return bytesParsed, nil
	}
	if r.State == 2 {
		bytesParsed, done, err := r.Headers.Parse(data)
		if err != nil {
			return 0, err
		}
		if bytesParsed == 0 {
			return 0, nil
		}
		if done {
			r.State = 0
		}
		return bytesParsed, nil
	}
	if r.State == 0 {
		return 0, fmt.Errorf("failed to parse in done state")
	}
	return 0, fmt.Errorf("invalid state %d", r.State)

}

func (cr *chunkReader) Read(p []byte) (n int, err error) {
	if cr.pos >= len(cr.data) {
		return 0, io.EOF
	}
	endIndex := cr.pos + cr.numBytesPerRead
	if endIndex > len(cr.data) {
		endIndex = len(cr.data)
	}
	n = copy(p, cr.data[cr.pos:endIndex])
	cr.pos += n

	return n, nil
}
