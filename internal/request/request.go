package request

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"strings"
)

type Request struct {
	RequestLine RequestLine
	state       requestState
}

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}
type requestState int

const (
	requestStateInitialized requestState = iota
	requestStateDone
)
const crlf = "\r\n"
const bufferSize = 8

func RequestFromReader(reader io.Reader) (*Request, error) {
	buf := make([]byte, bufferSize)
	readToIndex := 0
	req := &Request{
		state: requestStateInitialized,
	}
	//
	for req.state != requestStateDone {
		// Handle out of buffer space
		if readToIndex >= len(buf) {
			newBuf := make([]byte, len(buf)*2)
			copy(newBuf, buf)
			buf = newBuf
		}
		//
		numBytesRead, err := reader.Read(buf[readToIndex:])
		if err != nil {
			if errors.Is(err, io.EOF) {
				req.state = requestStateDone
				break
			}
			return nil, err
		}
		readToIndex += numBytesRead
		//
		numBytesParsed, err := req.parse(buf[:readToIndex])
		if err != nil {
			return nil, err
		}
		//
		copy(buf, buf[numBytesParsed:])
		readToIndex -= numBytesParsed
	}
	//
	return req, nil
}

func parseRequestLine(data []byte) (*RequestLine, int, error) {
	idx := bytes.Index(data, []byte(crlf))
	if idx == -1 {
		return nil, 0, nil
	}
	requestLineText := string(data[:idx])
	requestLine, err := requestLineFromString(requestLineText)
	if err != nil {
		return nil, 0, nil
	}
	//
	return requestLine, idx + 2, nil
}

func requestLineFromString(str string) (*RequestLine, error) {
	parts := strings.Split(str, " ")
	if len(parts) != 3 {
		return nil, fmt.Errorf("poorly formatted request-line: %s", str)
	}
	//
	method := parts[0]
	// Loop over each chacter and ensure that it's alphabetic
	for _, c := range method {
		if c < 'A' || c > 'Z' {
			return nil, fmt.Errorf("invalid method: %s", method)
		}
	}
	//
	requestTarget := parts[1]
	//
	versionParts := strings.Split(parts[2], "/")
	if len(versionParts) != 2 {
		return nil, fmt.Errorf("malformed start-line %s", str)
	}
	//
	httpPart := versionParts[0]
	if httpPart != "HTTP" {
		return nil, fmt.Errorf("unrecognized HTTP-part: %s", httpPart)
	}
	//
	version := versionParts[1]
	if version != "1.1" {
		return nil, fmt.Errorf("unrecognixed HTTP-version: %s", version)
	}
	//
	return &RequestLine{
		Method:        method,
		RequestTarget: requestTarget,
		HttpVersion:   versionParts[1],
	}, nil
}

func (r *Request) parse(data []byte) (int, error) {
	// Accepts the next slice of bytes that needs to be parsed into the Request struct. But where? the Request struct doesn't have raw data so how
	switch r.state {
	case requestStateInitialized:
		requestLine, n, err := parseRequestLine(data)
		// Handle errors from parseRequestLine
		if err != nil {
			return 0, err
		}
		// parseRequestLine needs more data
		if n == 0 {
			return 0, nil
		}
		// store the requestLineStruct in r.RequestLine
		r.RequestLine = *requestLine
		// update r.state
		r.state = requestStateDone
		// return num successfully consumed and parsed bytes
		return n, nil
	case requestStateDone:
		return 0, fmt.Errorf("error: trying to read data in a done state")
	default:
		return 0, fmt.Errorf("error: unknown state %v", r.state)
	}
}
