package main

import (
	"github.com/evok02/httpfromtcp/internal/server"
	"github.com/evok02/httpfromtcp/internal/response"
	"github.com/evok02/httpfromtcp/internal/request"
	"os/signal"
	"syscall"
	"os"
	"log"
)
const port = 42069

func HandlerFunction(w response.Writer, r *request.Request) {
	w.Headers.Set("Content-Type", "text/html")

	if r.RequestLine.RequestTarget == "/yourproblem" {

		body := []byte(`
		<html>
		  <head>
			<title>400 Bad Request</title>
		  </head>
		  <body>
			<h1>Bad Request</h1>
			<p>Your request honestly kinda sucked.</p>
		  </body>
		</html>
			`)
		w.WriteHeaders(400)
		w.WriteBody(body)
		return
	}

	if r.RequestLine.RequestTarget == "/myproblem" {
		body := []byte(`
		<html>
		  <head>
			<title>500 Bad Request</title>
		  </head>
		  <body>
			<h1>Internal Server Error</h1>
			<p>Okay, you know what? This one is on me.</p>
		  </body>
		</html>
			`)
		w.WriteHeaders(500)
		w.WriteBody(body)
		return
	}
	body := []byte(`
	<html>
	  <head>
		<title>200 OK</title>
	  </head>
	  <body>
		<h1>Success!</h1>
		<p>Your request was an absolute banger.</p>
	  </body>
	</html>
		`)
	w.WriteHeaders(200)
	w.WriteBody(body)
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
