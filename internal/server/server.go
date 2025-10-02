package server

import (
	"fmt"
	"log"
	"net"
	"sync/atomic"
)

// Contains the state of the server. Server is an HTTP 1.1 server
type Server struct {
	listener net.Listener
	closed   atomic.Bool
}

// Creates a net.Listener and returns a new Server instance. Starts listening for requests inside a goroutine.
func Serve(port int) (*Server, error) {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return nil, err
	}
	s := &Server{
		listener: listener,
	}
	go s.listen()
	return s, nil
}

// Closes the listener and the server
func (s *Server) Close() error {
	s.closed.Store(true)
	if s.listener != nil {
		return s.listener.Close()
	}
	return nil
}

// Uses a loop to .Accept new connection as they come in, and handles each one in a new goroutine.
func (s *Server) listen() {
	for {
		conn, err := s.listener.Accept()
		if err != nil {
			if s.closed.Load() {
				return
			}
			log.Printf("Error accepting connection: %v", err)
			continue
		}
		go s.handle(conn)
	}
}

// Handles a single connection by writing the following response and then closing the connection.
func (s *Server) handle(conn net.Conn) {
	defer conn.Close()
	response := "HTTP/1.1 200 OK\r\n" + // Status line
		"Content-Type: text/plain\r\n" + // Example header
		// "Content-Length: 13\r\n" + // Content length header
		"\r\n" + // Blank line to separate headers from the body
		"Hello World!\n" //Body
	conn.Write([]byte(response))
	return
}
