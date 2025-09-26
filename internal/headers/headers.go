package headers

import (
	"bytes"
	"fmt"
	"strings"
)

const crlf = "\r\n"

type Headers map[string]string

func NewHeaders() Headers {
	return map[string]string{}
}

func (h Headers) Parse(data []byte) (n int, done bool, err error) {
	// Look for CRLF.
	idx := bytes.Index(data, []byte(crlf))
	// If there is not a CRLF then we need more data.
	if idx == -1 {
		return 0, false, nil
	}
	// if found a CRLF at the beginning of the data then we have the end of the headers. Return the proper values ASAP.
	if idx == 0 {
		// the empty line.
		// headers are done, consume the CRLF
		return 2, true, nil
	}
	//
	parts := bytes.SplitN(data[:idx], []byte(":"), 2)
	key := string(parts[0])
	// If the length of sliceKeyValue is greater than 2 then we know we have more than a key-value pair
	if key != strings.TrimRight(key, " ") {
		return 0, false, fmt.Errorf("invalid header name: %s", key)
	}
	//
	value := bytes.TrimSpace(parts[1])
	key = strings.TrimSpace(key)
	//
	h.Set(key, string(value))
	//
	return idx + 2, false, nil
}

func (h Headers) Set(key, value string) {
	h[key] = value
}
