package request

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/evok02/httpfromtcp/internal/headers"
	"io"
	"testing"
)

type chunkReader struct {
	data            string
	numBytesPerRead int
	pos             int
}

// Read reads up to len(p) or numBytesPerRead bytes from the string per call
// its useful for simulating reading a variable number of bytes per chunk from a network connection
func (cr *chunkReader) Read(p []byte) (n int, err error) {
	if cr.pos >= len(cr.data) {
		return 0, io.EOF
	}
	endIndex := cr.pos + cr.numBytesPerRead
	if endIndex > len(cr.data) {
		endIndex = len(cr.data)
	}
	n = copy(p, cr.data[cr.pos:endIndex])
	cr.pos += n

	return n, nil
}

func TestRequestLineParse(t *testing.T) {
	// Test: Standard Headers
	reader := &chunkReader{
		data:            "GET / HTTP/1.1\r\nHost: localhost:42069\r\nUser-Agent: curl/7.81.0\r\nAccept: */*\r\n\r\n",
		numBytesPerRead: 3,
	}
	r, err := RequestFromReader(reader)
	require.NoError(t, err)
	require.NotNil(t, r)
	assert.Equal(t, "localhost:42069", r.Headers["host"])
	assert.Equal(t, "curl/7.81.0", r.Headers["user-agent"])
	assert.Equal(t, "*/*", r.Headers["accept"])

	// Test: Malformed Header
	reader = &chunkReader{
		data:            "GET / HTTP/1.1\r\nHost localhost:42069\r\n\r\n",
		numBytesPerRead: 3,
	}
	r, err = RequestFromReader(reader)
	require.Error(t, err)

	//Test: Empty Headers
	reader = &chunkReader{
		data:            "GET / HTTP/1.1\r\n\r\n",
		numBytesPerRead: 3,
	}
	r, err = RequestFromReader(reader)
	assert.Equal(t, r.RequestLine, RequestLine{
		Method: "GET",
		RequestTarget: "/",
		HttpVersion: "1.1",
	})
	assert.Equal(t, r.Headers, headers.Headers{})

	//Test: Duplicate Headers
	reader = &chunkReader{
		data:            "GET / HTTP/1.1\r\nAccept: */jpg\r\nAccept: */png\r\nBrowser: Firefox\r\n\r\n",
		numBytesPerRead: 3,
	}
	r, err = RequestFromReader(reader)
	assert.Equal(t, r.Headers["accept"], "*/jpg, */png")
	assert.Equal(t, r.Headers["browser"], "Firefox")

	//Test: Case sensitive headers
	reader = &chunkReader{
		data:            "GET / HTTP/1.1\r\nACCEPT: */png\r\nBROWSER: Firefox\r\n\r\n",
		numBytesPerRead: 3,
	}
	r, err = RequestFromReader(reader)
	assert.Equal(t, r.Headers["accept"], "*/png")
	assert.Equal(t, r.Headers["browser"], "Firefox")
}
