package response

import (
	//"log"
	"strconv"
	"io"
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

func WriteStatusLine(w io.Writer, statusCode StatusCode) error {
	switch statusCode {
	case StatusOK: 
		w.Write([]byte("HTTP/1.1 200 OK" + CRLF))
	case StatusBadRequest: 
		w.Write([]byte("HTTP/1.1 400 BadRequest" + CRLF))
	case StatusInternalError: 
		w.Write([]byte("HTTP/1.1 500 InternalError" + CRLF))
	default:
		return ERROR_UNSUPPORTED_STATUS_CODE
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

func WriteHeaders(w io.Writer, headers headers.Headers) error {
	contentLength, ok := headers.Get("Content-Length")
	if !ok {
		WriteStatusLine(w, 400)
		return ERROR_MISSING_DEFAULT_HEADER
	}

	contentType, ok := headers.Get("Content-Type")
	if !ok {
		WriteStatusLine(w, 400)
		return ERROR_MISSING_DEFAULT_HEADER
	}

	connection, ok := headers.Get("Connection")
	if !ok {
		WriteStatusLine(w, 400)
		return ERROR_MISSING_DEFAULT_HEADER
	}


	
	WriteStatusLine(w, 200)
	w.Write([]byte("Content-Length: " + contentLength + CRLF))
	w.Write([]byte("Content-Type: " + contentType + CRLF))
	w.Write([]byte("Connection: " + connection + CRLF + CRLF))
	return nil
}


