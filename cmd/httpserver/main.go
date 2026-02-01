package main

import (
	"strconv"
	"fmt"
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
	"crypto/sha256"
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
func ProxyHandler(w response.Writer, r *request.Request) {
	if path, found := strings.CutPrefix(r.RequestLine.RequestTarget, "/httpbin"); found {
		url := "https://httpbin.org" + path

		buf := make([]byte, buffSize)
		//bufWriter := bytes.NewBuffer(buf)
		res, err := http.Get(url)
		if err != nil {
			server.WriteError(w, server.NewError(err.Error(), response.StatusInternalError))
		}

		w.Headers.Set("Transfer-Encoding", "chunked")
		w.Headers.Set("Trailer", "X-Content-SHA256")
		w.Headers.Set("Trailer", "X-Content-Length")
		w.WriteHeaders(response.StatusOK)

		for {
			rN, err := res.Body.Read(buf)
			if err != nil {
				if err == io.EOF {
					break 
				}
				server.WriteError(w, server.NewError(err.Error(), response.StatusInternalError))
			}

			buf = buf[:rN]
			wN, err := w.WriteChunkedBody(buf)
			if err != nil {
				server.WriteError(w, server.NewError(err.Error(), response.StatusInternalError))
			}
			fmt.Printf("%d\t%d\n", wN, rN)
			
		}

		_, err = w.WriteChunkedBodyDone()
		if err != nil {
			server.WriteError(w, server.NewError(err.Error(), response.StatusInternalError))
		}
		
		checksum := sha256.Sum256(w.Body)
		trailerHeaders := map[string]string{
			"X-Content-Sha256": fmt.Sprintf("%x", checksum),
			"X-Content-Length": strconv.Itoa(len(w.Body)),
		}

		err = w.WriteTrailers(trailerHeaders)
	}
}

func BinaryHandler(w response.Writer, r *request.Request) {
	if r.RequestLine.RequestTarget == "/video" {
		w.Headers.Set("Content-Type", "video/mp4")
		w.WriteHeaders(response.StatusOK)
	}
}

func main() {
	server, err := server.Serve(port, BinaryHandler)
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
