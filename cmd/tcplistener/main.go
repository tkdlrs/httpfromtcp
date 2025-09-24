package main

import (
	"fmt"
	"log"
	"net"

	"github.com/tkdlrs/httpfromtcp/internal/request"
)

const port = ":42069"

func main() {
	listener, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("error listening for TCP traffic: %s\n", err.Error())
	}
	defer listener.Close()
	//
	fmt.Println("Listening for TCP traffic on", port)
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatalf("error: %s\n", err.Error())
			continue
		}
		fmt.Println("Accepted connection from", conn.RemoteAddr())
		//
		linesChan, err := request.RequestFromReader(conn)
		if err != nil {
			fmt.Errorf("%v", err)
			return
		}
		printTemplate := fmt.Sprintf(`Request line:
- Method: %s
- Target: %s
- Version: %s
`, linesChan.RequestLine.Method, linesChan.RequestLine.RequestTarget, linesChan.RequestLine.HttpVersion)

		fmt.Printf(printTemplate)
	}
}
