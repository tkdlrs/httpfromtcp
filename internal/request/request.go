package request

import (
	"fmt"
	"io"
	"log"
)

type Request struct {
	RequestLine RequestLine
}

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

// must end up like the commented out one below...
// func RequestFromReader(reader io.Reader) (*Request, error) {
// -
func RequestFromReader(reader io.Reader) error {
	b, err := io.ReadAll(reader)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%s", b)
	return nil
}

func parseRequestLine(b string) RequestLine {

}
