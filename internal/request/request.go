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
	pState      parserState
}

type parserState int

const (
	initialized parserState = iota
	done
)

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

const crlf = "\r\n"

func RequestFromReader(reader io.Reader) (*Request, error) {
	theRequest := Request{}
	theRequest.pState = initialized
	bufferSize := 8
	buf := make([]byte, bufferSize)
	readBytes := 0
	parsedBytes := 0
	for theRequest.pState != done {
		// Handle out of buffer space
		if len(buf) == bufferSize {
			bufferSize *= 2
			bufCopy := buf
			buf = make([]byte, bufferSize)
			copy(buf, bufCopy)
		}
		//
		n, err := reader.Read(buf[readBytes:])
		readBytes += n
		if err != nil {
			if errors.Is(err, io.EOF) {
				theRequest.pState = done
				break
			}
			return nil, fmt.Errorf("error reading theRequest: %v", err)
		}
		//
		nb, err := theRequest.parse(buf[:readBytes])
		if err != nil {
			return nil, fmt.Errorf("error parsing theRequest: %v", err)
		}
		parsedBytes += nb
	}
	//
	return &theRequest, nil
}

func parseRequestLine(data []byte) (int, error) {
	idx := bytes.Index(data, []byte(crlf))
	if idx == -1 {
		return 0, nil
	}
	//
	return idx + len(crlf), nil
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
		HttpVersion:   versionParts[1],
		RequestTarget: requestTarget,
		Method:        method,
	}, nil
}

func (r *Request) parse(data []byte) (int, error) {
	// Accepts the next slice of bytes that needs to be parsed into the Request struct. But where? the Request struct doesn't have raw data so how
	if r.pState == initialized {
		consumededBytes, err := parseRequestLine(data)
		// Handle errors from parseRequestLine
		if err != nil {
			return 0, err
		}
		// parseRequestLine needs more data
		if consumededBytes == 0 {
			return 0, nil
		}
		// to get here consumedBytes is > 0, meaning we have a full request line.
		// Extract that part of the data. -don't include the crlf in the text
		requestLineText := string(data[:consumededBytes-len(crlf)])

		// call requestLineFromString to get the actual RequestLine struct
		requestLineStruct, err := requestLineFromString(requestLineText)
		if err != nil {
			return 0, err
		}
		// store the requestLineStruct in r.RequestLine
		r.RequestLine = *requestLineStruct
		// update r.pState
		r.pState = done
		// return num successfully consumed and parsed bytes
		return consumededBytes, nil
	} else if r.pState == done {
		return -1, fmt.Errorf("error: trying to read data in a done state %v", r.pState)
	} else {
		return -1, fmt.Errorf("error: unknown state %v", r.pState)
	}
}
