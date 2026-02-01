package response

import (
	"net"
	"fmt"
	"io"
	"strconv"
	"errors"
	"github.com/evok02/httpfromtcp/internal/headers"
)

type StatusCode int

const (
	StatusOK StatusCode = 200
	StatusBadRequest StatusCode = 400
	StatusInternalError StatusCode = 500
)

const CRLF = "\r\n"

var ERROR_UNSUPPORTED_STATUS_CODE = errors.New("unsupported status code")
var ERROR_MISSING_DEFAULT_HEADER = errors.New("error missing default header")

func (w *Writer)WriteStatusLine(statusCode StatusCode) error {
	switch statusCode {
	case StatusOK: 
		w.StatusCode = StatusOK
		w.StatusLine = "HTTP/1.1 200 OK" + CRLF
	case StatusBadRequest: 
		w.StatusCode = StatusBadRequest
		w.StatusLine = "HTTP/1.1 400 OK" + CRLF
	case StatusInternalError: 
		w.StatusCode = StatusInternalError
		w.StatusLine = "HTTP/1.1 500 OK" + CRLF
	default:
		return ERROR_UNSUPPORTED_STATUS_CODE
	}

	_, err := w.destBuffer.Write([]byte(w.StatusLine))
	if err != nil {
		return err
	}
	return nil
}

func GetDefaultHeaders(contentLen int) headers.Headers {
	h := headers.NewHeaders()
	h.Set("Content-Length", strconv.Itoa(contentLen))
	h.Set("Connection", "close")
	h.Set("Content-Type", "text/plain")
	return h
}

func (w *Writer) WriteHeaders(status StatusCode) error {
	w.WriteStatusLine(status)
	for k, v := range w.Headers {
		headerLine := fmt.Sprintf("%s: %s\r\n", k, v)
		_, err := w.destBuffer.Write([]byte(headerLine))
		if err != nil {
			return err
		}
	}

	return nil
}

func (w *Writer) WriteBody(body []byte) (int, error) {
	w.Body = append(w.Body, body...)
	n, err := w.destBuffer.Write(w.Body)
	return n, err
}

type Writer struct {
	StatusCode StatusCode
	StatusLine StatusLine
	Headers headers.Headers
	Body []byte
	destBuffer io.Writer
}

func NewWriter(c net.Conn) *Writer {
	return &Writer{
		StatusCode: StatusOK,
		StatusLine: "HTTP/1.1 200 OK",
		Headers: headers.Headers{},
		Body: make([]byte, 0),
		destBuffer: c,
	}
}

type StatusLine string
