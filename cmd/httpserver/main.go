package main

import (
	"github.com/evok02/httpfromtcp/internal/server"
	"github.com/evok02/httpfromtcp/internal/response"
	"github.com/evok02/httpfromtcp/internal/request"
	"os/signal"
	"syscall"
	"os"
	"log"
	"strings"
	"net/http"
	"io"
)
const port = 42069
const buffSize = 32

func asyncRead(r io.Reader) <- chan []byte {
	stream := make(chan []byte)	
	buf := make([]byte, buffSize)
	go func() {
		defer close(stream)
		for {
			n, err := r.Read(buf)	
			if err != nil {
				if err == io.EOF {
					return
				}
				// TODO: create channel for errors 
				println(err.Error())
			}
			stream <- buf[:n]
			copy(buf, buf[n:])
		}
	}()
	return stream
}

func asyncWrite(w io.Writer, in <- chan []byte) <- chan []byte {
	buf := make([]byte, buffSize)
	stream := make(chan []byte)
	go func() {
		for rawBytes := range in {
			n, err := w.Write(rawBytes)
			if err != nil {
				//TODO: should i past the same errChan, or create different for each one?
				println(err.Error())
			}
			stream <- buf[:n]
			copy(buf, buf[n:])
		}
	}()
	return stream
}
func HandlerFunction(w response.Writer, r *request.Request) {
	if path, found := strings.CutPrefix(r.RequestLine.RequestTarget, "/httpbin"); found {
		url := "https://httpbin.org" + path

		buf := make([]byte, buffSize)
		//bufWriter := bytes.NewBuffer(buf)
		res, err := http.Get(url)
		if err != nil {
			server.WriteError(w, server.NewError(err.Error(), response.StatusInternalError))
		}

		w.Headers.Set("Transfer-Encoding", "chunked")
		w.WriteHeaders(response.StatusOK)

		for {
			nReadBytes, err := res.Body.Read(buf)
			if err != nil {
				if err == io.EOF {
					break
				}
				server.WriteError(w, server.NewError(err.Error(), response.StatusInternalError))
			}

			nWriteBytes, err := w.WriteChunkedBody(buf)
			if err != nil {
				server.WriteError(w, server.NewError(err.Error(), response.StatusInternalError))
			}
			copy(buf, buf[:nReadBytes - nWriteBytes])
		}

		_, err = w.WriteChunkedBodyDone()
		if err != nil {
			server.WriteError(w, server.NewError(err.Error(), response.StatusInternalError))
		}
	}
}

func main() {
	server, err := server.Serve(port, HandlerFunction)
	if err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
	defer server.Close()
	log.Println("Server started on port", port)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan
	log.Println("Server gracefully stopped")
}
