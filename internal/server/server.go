package server

import (
	"io"
	"net"
	"log"
	"strconv"
	"sync/atomic"
	"github.com/evok02/httpfromtcp/internal/request"
	"github.com/evok02/httpfromtcp/internal/response"
	"bytes"
)

type Server struct {
	addr string
	connection net.Listener
	handler Handler
	closed atomic.Bool
}

func Serve(port int, h Handler) (*Server, error){
	p := strconv.Itoa(port)
	addr := "localhost:" + p

	l, err := net.Listen("tcp", addr)
	if err != nil {
		return nil, err
	}

	s := &Server{
		addr: addr,
		connection: l,
		handler: h,
	}

	go s.listen()
	return s, nil
}

func (s *Server) Close() error {
	s.closed.Store(true)
	if s.connection != nil {
		return s.connection.Close()
	}
	return nil
}

func (s *Server) handle(c net.Conn) {
	req, err := request.RequestFromReader(c)
	if err != nil {
		WriteError(c, NewError(err.Error(), response.StatusBadRequest))
		return
	}
	defer c.Close()

	var buf bytes.Buffer
	handlerErr := s.handler(&buf, req)
	if handlerErr != nil {
		WriteError(c, handlerErr)
		return
	}

	response.WriteStatusLine(c, response.StatusOK)
	headers := response.GetDefaultHeaders(buf.Len())
	response.WriteHeaders(c, headers)
	c.Write(buf.Bytes())
}

func (s *Server) listen() {
	for {
		conn, err := s.connection.Accept()
		if err != nil {
			log.Printf("Connection(addr: %s) lost: %s", conn.RemoteAddr().String(), err.Error())
		}
		s.handle(conn)
	}
}

type HandlerError struct {
	Msg string
	StatusCode response.StatusCode
}

func (e *HandlerError) Error() string {
	return e.Msg
}

func NewError(msg string, status response.StatusCode) *HandlerError {
	return &HandlerError{
		Msg: msg,
		StatusCode: status,
	}
}


func  WriteError(w io.Writer, e *HandlerError) {
	response.WriteStatusLine(w, e.StatusCode)
	headers := response.GetDefaultHeaders(len(e.Msg))
	response.WriteHeaders(w, headers)
	w.Write([]byte("\r\n" + e.Msg))
}

type Handler func(io.Writer, *request.Request) *HandlerError




