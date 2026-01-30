package server

import (
	"io"
	"net"
	"log"
	"strconv"
	"github.com/evok02/httpfromtcp/internal/request"
	"github.com/evok02/httpfromtcp/internal/response"
	"bytes"
)

type Server struct {
	Addr string
	Connection net.Listener
	Handler Handler
}

func Serve(port int, h Handler) (*Server, error){
	p := strconv.Itoa(port)
	addr := "localhost:" + p

	l, err := net.Listen("tcp", addr)
	if err != nil {
		return nil, err
	}

	s := &Server{
		Addr: addr,
		Connection: l,
		Handler: h,
	}

	go func() {
		s.listen()
	}()
	return s, nil
}

func (s *Server) Close() error {
	return s.Connection.Close()
}

func (s *Server) handle(c net.Conn) {
	req, err := request.RequestFromReader(c)
	if err != nil {
		WriteError(c, NewError(err.Error(), response.StatusBadRequest))
	}
	defer c.Close()

	var buf bytes.Buffer
	handlerErr := s.Handler(&buf, req)
	if err != nil {
		WriteError(c, handlerErr)
	}


	response.WriteStatusLine(c, response.StatusOK)
	headers := response.GetDefaultHeaders(buf.Len())
	response.WriteHeaders(c, headers)
	c.Write(buf.Bytes())
}

func (s *Server) listen() {
	for {
		conn, err := s.Connection.Accept()
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




