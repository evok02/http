package main

import (
	"github.com/evok02/httpfromtcp/internal/server"
	"github.com/evok02/httpfromtcp/internal/request"
	"os/signal"
	"syscall"
	"os"
	"log"
	"io"
)
const port = 42069

func HandlerFunction(w io.Writer, r *request.Request) *server.HandlerError {
	if r.RequestLine.RequestTarget == "/yourproblem" {
		return server.NewError("Your rpoblem is not my problem\n", 400)
	}

	if r.RequestLine.RequestTarget == "/myproblem" {
		return server.NewError("Woopsie, my bad\n", 500)
	}

	w.Write([]byte("All good, frfr\n"))
	return nil
}

func main() {
	server, err := server.Serve(port, HandlerFunction)
	if err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
	defer server.Close()
	log.Println("Server started on port", port)

	sigChan := make(chan os.Signal)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan
	log.Println("Server gracefully stopped")
}
