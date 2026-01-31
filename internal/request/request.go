package request

import (
	"strconv"
	"fmt"
	"github.com/evok02/httpfromtcp/internal/headers"
	"bytes"
	"strings"
	"io"
	"errors"
)

type Request struct {
	RequestLine RequestLine
	Headers headers.Headers
	Body []byte
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
	reqStateParsingBody
	reqStateDone
)

const buffSize = 8
var contentLength int


var	ERROR_INVALID_FORMAT = errors.New("malformed request-line")
var ERROR_INVALID_METHOD_FORMATING = errors.New("method should consist of capital errors")
var ERROR_INVALID_PROTOCOL_VERSION = errors.New("poorly formatted protocol version")
var ERROR_INVALID_REQUEST_STATE = errors.New("unknown state of request")
var ERROR_MALFORMED_BODY = errors.New("request with malformed body")
var ERROR_MALFORMED_CONTENT_LENGTH = errors.New("request with malformed content length") 


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
		fmt.Printf("request state: %d\n", req.State)
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
			val, err := r.Headers.Get("Content-Length")
			if err != nil {
				r.State = reqStateDone
				return 0, nil
			}
			contentLength, err = strconv.Atoi(val)
			if err != nil {
				return 0, ERROR_MALFORMED_CONTENT_LENGTH
			}
			r.State = reqStateParsingBody 
			return len(crlf), nil
		}
		return numBytes , nil

	case reqStateParsingBody:
		r.Body = append(r.Body, data...)
		if len(r.Body) > contentLength {
			return 0, ERROR_MALFORMED_BODY
		} else if len(r.Body) == contentLength {
			r.State = reqStateDone
		} else {
			if string(data) == string(crlf) {
				return 0, ERROR_MALFORMED_BODY
			}
			return len(data), nil
		}
	default:
		return 0, ERROR_INVALID_REQUEST_STATE
	}
	return 0, nil
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
