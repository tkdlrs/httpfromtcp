package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"strings"
)

func main() {
	ln, err := net.Listen("tcp", ":42069")
	if err != nil {
		log.Fatalf("could not listen to connection %s \n", err)
	}
	defer ln.Close()
	//
	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Printf("could not accept connection %s \n", err)
			continue
		}
		fmt.Println("accepted connection")
		//
		linesChan := getLinesChannel(conn)
		//
		for line := range linesChan {
			line = strings.TrimSuffix(line, "\r")
			fmt.Println(line)
		}
		fmt.Println("connection closed")
	}
	//
}

func getLinesChannel(f io.ReadCloser) <-chan string {
	lines := make(chan string)
	go func() {
		defer f.Close()
		defer close(lines)
		currentLineContents := ""
		for {
			buffer := make([]byte, 8)
			n, err := f.Read(buffer)
			if err != nil {
				if currentLineContents != "" {
					lines <- currentLineContents
				}
				if errors.Is(err, io.EOF) {
					break
				}
				fmt.Printf("error: %s\n", err.Error())
				break
			}
			//
			str := string(buffer[:n])
			parts := strings.Split(str, "\n")

			for i := 0; i < len(parts)-1; i++ {
				lines <- fmt.Sprintf("%s%s", currentLineContents, parts[i])
				currentLineContents = ""
			}
			currentLineContents += parts[len(parts)-1]
		}
	}()
	return lines
}
