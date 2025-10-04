package server

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net"
	"sync/atomic"

	"github.com/tkdlrs/httpfromtcp/internal/request"
	"github.com/tkdlrs/httpfromtcp/internal/response"
)

type Handler func(w io.Writer, req *request.Request) *HandlerError

type HandlerError struct {
	StatusCode response.StatusCode
	Message    string
}

// logic that writes a HandlerError to an io.Writer.
func (he HandlerError) Write(w io.Writer) {
	response.WriteStatusLine(w, he.StatusCode)
	messageBytes := []byte(he.Message)
	headers := response.GetDefaultHeaders(len(messageBytes))
	response.WriteHeaders(w, headers)
	w.Write(messageBytes)
}

// Server is an HTTP 1.1 server
type Server struct {
	handler  Handler
	listener net.Listener
	closed   atomic.Bool
}

func Serve(port int, handler Handler) (*Server, error) {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return nil, err
	}
	s := &Server{
		handler:  handler,
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
	// parse the request from the connection
	req, err := request.RequestFromReader(conn)
	if err != nil {
		hErr := &HandlerError{
			StatusCode: response.StatusCodeBadRequest,
			Message:    err.Error(),
		}
		hErr.Write(conn)
		return
	}
	// New empty bytes.Buffer for handler to write to
	buf := bytes.NewBuffer([]byte{})
	// call the handler function
	hErr := s.handler(buf, req)
	// if handler errors, write the error to the connection
	if hErr != nil {
		hErr.Write(conn)
		return
	}
	// if handler succeeds
	b := buf.Bytes()
	// Create a new default response header
	// write the status line
	response.WriteStatusLine(conn, response.StatusCodeSuccess)
	// write headers
	headers := response.GetDefaultHeaders(len(b))
	// write response body from handler's buffer
	response.WriteHeaders(conn, headers)
	conn.Write(b)
	return
}
