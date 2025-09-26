package headers

import (
	"bytes"
	"fmt"
)

type Headers map[string]string

const crlf = "\r\n"

func (h Headers) Parse(data []byte) (n int, done bool, err error) {
	//
	bytesConsumed := 0

	// Look for CRLF.
	idx := bytes.Index(data, []byte(crlf))
	// If there is not a CRLF then we need more data.
	if idx == -1 {
		return bytesConsumed, false, nil
	}
	// if found a CRLF at the beginning of the data then we have the end of the headers. Return the proper values ASAP.
	if idx == 0 {
		return bytesConsumed, true, nil
	}
	//
	dat := bytes.TrimSpace(data)
	sliceKeyValue := bytes.Fields(dat)
	// If the length of sliceKeyValue is greater than 2 then we know we have more than a key-value pair
	if len(sliceKeyValue) > 2 {
		return bytesConsumed, true, fmt.Errorf("error length of key value pair is longer than expected")
	}
	key := sliceKeyValue[0]
	val := sliceKeyValue[1]
	//
	h[string(key)[:(len(key)-1)]] = string(val)
	fmt.Println("string(key)", string(key))
	fmt.Println("h", h)
	//
	return bytesConsumed, true, nil
}
