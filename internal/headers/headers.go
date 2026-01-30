package headers

import (
	"bytes"
	"strings"
	"errors"
	"fmt"
)

type Headers map[string]string

const CRLF = "\r\n"

var ERROR_MALFORMED_HEADER = errors.New("malformed header block")
var ERROR_INVALID_HEADER_NAME = errors.New("invalid header name")

func NewHeaders() Headers {
	return make(Headers)
}

func (h Headers) Parse(data []byte) (n int, done bool, err error) {
	idx := bytes.Index(data, []byte(CRLF))
	switch idx {
	case -1:
		return 0, false, nil
	case 0:
		return 0, true, nil
	}

	key, value, err := parseFieldLine(data[:idx])
	if err != nil {
		return 0, false, err
	}

	h.Set(key, value)

	return idx + len(CRLF), false, nil
}

func parseFieldLine(line []byte) (string, string, error) {
	idx := bytes.Index(line, []byte(":"))
	if idx == -1 {
		return "", "", ERROR_MALFORMED_HEADER
	}
	sepLineStr := []string{string(line[:idx]), string(line[idx + 1:])}
	
	key, value := strings.ToLower(sepLineStr[0]), sepLineStr[1]
	if err := validateFieldName(key); err != nil {
		return "", "", err
	}

	return strings.TrimSpace(key), strings.TrimSpace(value), nil
}

func validateFieldName(token string) error {
	if len(token) < 1 {
		return ERROR_MALFORMED_HEADER
	}
	if token[len(token) - 1] == ' ' {
		return ERROR_MALFORMED_HEADER
	}
	if err := validateTokenChars(token); err != nil {
		return err
	}
	return nil
}

func validateTokenChars(token string) error {
	for _, char := range token {
		if  char == 32 || char == 34 || char == 40 || char == 41 || char == 44 || char == 47 ||
			(char > 57 && char < 65) || char > 126 {
			return ERROR_MALFORMED_HEADER
		}
	}
	return nil
}

func (h Headers) Get(token string) (string, error) {
	val, ok := h[strings.ToLower(token)]
	if !ok {
		return "", ERROR_INVALID_HEADER_NAME
	}
	return val, nil
}

func (h Headers) Set(token string, value string) {
	if b, _ := h.IsExist(strings.ToLower(token)); b {
		h[strings.ToLower(token)] += fmt.Sprintf(", %s", value)
	} else {
		h[strings.ToLower(token)] = value
	}
}

func (h Headers) IsExist(token string) (bool, string) {
	if val, err := h.Get(token); err == nil {
		return true, val
	}
	return false, ""
}

