package server

import (
	"net"
	"log"
	"strconv"
	"sync/atomic"
	"github.com/evok02/httpfromtcp/internal/request"
	"github.com/evok02/httpfromtcp/internal/response"
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
	responseWriter := response.NewWriter(c)
	if err != nil {
		WriteError(*responseWriter, NewError(err.Error(), response.StatusInternalError))
	}
	defer c.Close()

	s.handler(*responseWriter, req)
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


func  WriteError(w response.Writer, e *HandlerError) {
	w.WriteStatusLine(e.StatusCode)
	w.WriteHeaders(response.StatusInternalError)
	w.WriteBody([]byte(e.Msg))
}

type Handler func(response.Writer, *request.Request)




