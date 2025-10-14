package request

import (
	"fmt"
	"io"
	"strings"
)

const bufferSize = 8

func RequestFromReader(r io.Reader) (*Request, error) {
	buf := make([]byte, bufferSize)
	readToIndex := 0
	req := &Request{
		State: 1,
	}
	for req.State != 0 {
		if len(buf)-readToIndex < bufferSize {
			new := make([]byte, len(buf)*2)
			copy(new, buf)
			buf = new
		}
		n, err := r.Read(buf[readToIndex:])
		if err != nil {
			return nil, err
		}
		readToIndex += n
		if n == 0 {
			break
		}
		bytesParsed, err := req.parse(buf[:readToIndex])
		if err != nil {
			if err == io.EOF {
				req.State = 0
				break
			} else {
				return nil, err
			}
		}
		if bytesParsed == readToIndex {
			break
		}
		if bytesParsed != 0 {
			new := make([]byte, len(buf)-bytesParsed)
			copy(new, buf[bytesParsed:])
			buf = new
			readToIndex -= bytesParsed
		}
	}
	if req.RequestLine.HttpVersion != "1.1" {
		return nil, InvalidHTTPVersionError
	}
	if req.RequestLine.Method != "GET" && req.RequestLine.Method != "POST" {
		return nil, InvalidHTTPMethodError
	}
	return req, nil
}

func parseRequestLine(r string) (RequestLine, int, error) {
	lines := strings.Split(r, "\r\n")
	if len(lines) < 2 {
		return RequestLine{}, 0, nil
	}
	split := strings.Split(lines[0], " ")
	if len(split) != 3 {
		return RequestLine{}, len(r), InvalidRequestLineError
	}
	reqLine := RequestLine{
		Method:        strings.ToUpper(split[0]),
		RequestTarget: split[1],
		HttpVersion:   strings.Replace(split[2], "HTTP/", "", 1),
	}
	return reqLine, len(r), nil
}

func (r *Request) parse(data []byte) (int, error) {
	if r.State == 1 {
		reqLine, bytesParsed, err := parseRequestLine(string(data))
		if err != nil {
			return 0, err
		}
		if bytesParsed == 0 {
			return 0, nil
		}

		r.RequestLine = reqLine
		r.State = 0
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
