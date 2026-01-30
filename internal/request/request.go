package request

import (
	"github.com/evok02/httpfromtcp/internal/headers"
	"bytes"
	"strings"
	"io"
	"errors"
)

type Request struct {
	RequestLine RequestLine
	Headers headers.Headers
	State parseRequestStatus
	
}

type RequestLine struct {
	HttpVersion string
	RequestTarget string
	Method string
}

const crlf = "\r\n"

type parseRequestStatus int

const (
	reqStateInitialized parseRequestStatus = iota 
	reqStateParsingHeaders
	reqStateDone
)

const buffSize = 8


var	ERROR_INVALID_FORMAT = errors.New("malformed request-line")
var ERROR_INVALID_METHOD_FORMATING = errors.New("method should consist of capital errors")
var ERROR_INVALID_PROTOCOL_VERSION = errors.New("poorly formatted protocol version")
var ERROR_INVALID_REQUEST_STATE = errors.New("errors unknown state of request")



func RequestFromReader(r io.Reader) (*Request, error) {
	buf := make([]byte, buffSize)
	var readToIdx int
	req := NewRequest()
	for req.State != reqStateDone{
		if readToIdx >= len(buf) {
			extendo := make([]byte, len(buf) * 2)
			copy(extendo, buf)
			buf = extendo
		}
		n, err := r.Read(buf[readToIdx:])
		if err != nil {
			if err == io.EOF {
				req.State = reqStateDone
				break
			}
			return nil, err
		}
		readToIdx += n
		numParsedBytes, err := req.parse(buf[:readToIdx])
		if err != nil {
			return nil, err
		}
		copy(buf, buf[numParsedBytes:])
		readToIdx -= numParsedBytes
	}
	return req, nil
}

func (r *Request) parse(data []byte) (int, error) {
	switch r.State {
	case reqStateInitialized:
		numBytes, reqLine, err := parseRequestLine(data)
		if err != nil {
			return 0, err
		}
		if numBytes > 0 {
			r.State = reqStateParsingHeaders
			r.RequestLine = *reqLine
		}
		return numBytes, nil
	case reqStateParsingHeaders:
		numBytes, done, err := r.Headers.Parse(data)
		if err != nil {
			return 0, err
		}
		if done {
			r.State = reqStateDone
		}
		return numBytes, nil
	default:
		return 0, ERROR_INVALID_REQUEST_STATE
	}
}

func parseRequestLine(buf []byte) (int, *RequestLine, error) {
	idx := bytes.Index(buf , []byte(crlf))
	if idx == -1 {
		return 0, nil, nil
	}

	reqLine, err := requestLineFromString(string(buf[:idx]))
	if err != nil {
		return 0, nil, err
	}
	return idx + len(crlf), reqLine, nil
}

func requestLineFromString(s string) (*RequestLine, error) {
	parts := strings.Fields(s)
	
	if len(parts) != 3 {
		return nil, ERROR_INVALID_FORMAT 
	}

	method, path, proto := parts[0], parts[1], parts[2]

	for _, char := range method {
		if char < 'A' || char > 'Z' {
			return nil, ERROR_INVALID_METHOD_FORMATING
		}
	}

	if !strings.HasPrefix(proto, "HTTP/") {
		return nil, ERROR_INVALID_PROTOCOL_VERSION
	}

	version := strings.TrimPrefix(proto, "HTTP/")
	if  version != "1.1" {
		return nil, ERROR_INVALID_PROTOCOL_VERSION
	}

	return &RequestLine{
		Method: method,
		RequestTarget: path,
		HttpVersion: version,
	}, nil
}

func NewRequest() *Request {
	return &Request{
		State: reqStateInitialized,
		RequestLine: RequestLine{},
		Headers: headers.NewHeaders(),
	}
}
