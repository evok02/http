package headers

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHeaders(t *testing.T) {
	// Test: Valid single header
	headers := NewHeaders()
	data := []byte("Host: localhost:42069\r\n\r\n")
	n, done, err := headers.Parse(data)
	require.NoError(t, err)
	require.NotNil(t, headers)
	assert.Equal(t, 23, n)
	assert.False(t, done)

	// Test: Invalid spacing header
	headers = NewHeaders()
	data = []byte("     Host : localhost:42069           \r\n\r\n")
	n, done, err = headers.Parse(data)
	require.Error(t, err)
	assert.Equal(t, n, 0)
	assert.False(t, done)

	// Test: Valid Done
	headers = NewHeaders()
	data = []byte("\r\n\r\n")
	n, done, err = headers.Parse(data)
	require.True(t, done)
	require.Equal(t, 2, n)
	require.NoError(t, err)

	headers = NewHeaders()
	data = []byte("\r\n")
	n, done, err = headers.Parse(data)
	require.True(t, done)
	require.Equal(t, 2, n)
	require.NoError(t, err)

	//Test: 2 Headers
	headers = NewHeaders()
	data = []byte("Host: localhost:42069\r\nUser-Agent: curl/7.81.0\r\nAccept: */*\r\n\r\n")
	n, done, err = headers.Parse(data)
	require.False(t, done)
	require.Equal(t, n, 23)
	require.NoError(t, err)

	
	//Test: Valid header with extra white space
	headers = NewHeaders()
	data = []byte("Host : localhost:42069\r\n\r\n")
	n, done, err = headers.Parse(data)
	require.False(t, done)
	assert.Equal(t, n, 0)
	assert.Error(t, err)

	//Test: Case-insesitive GET
	headers = NewHeaders()
	data = []byte("Host: localhost:42069\r\n\r\n")
	n, done, err = headers.Parse(data)
	require.False(t, done)
	assert.Equal(t, n, 23)
	assert.NoError(t, err)
	val, ok := headers.Get("Host")
	assert.True(t, ok)
	assert.Equal(t, val, "localhost:42069")

	//Test: Case-sensitive GET
	headers = NewHeaders()
	data = []byte("Host: localhost:42069\r\n\r\n")
	n, done, err = headers.Parse(data)
	require.False(t, done)
	assert.NoError(t, err)
	val, ok = headers.Get("HOST")
	assert.True(t, ok)
	assert.Equal(t, val, "localhost:42069")

	//Test: Invalid token character
	headers = NewHeaders()
	data = []byte("H©st: localhost:42069\r\n\r\n")
	n, done, err = headers.Parse(data)
	require.False(t, done)
	assert.Error(t, err)
}
