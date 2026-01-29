package main

import (
	"github.com/evok02/httpfromtcp/internal/request"
	"fmt"
	"log"
	"net"
)


func main() {
	fmt.Println("Listening on :42069")
	l, err := net.Listen("tcp", ":42069")
	if err != nil {
		log.Fatal(err)
	}
	defer l.Close()

	for {
		conn, err := l.Accept()
		if err != nil {
			return
		}
		defer conn.Close()

		go func(c net.Conn) {
			req, err := request.RequestFromReader(conn)
			if err != nil {
				return
			}
			reqLine := req.RequestLine
			fmt.Printf("\nRequest line:\n- Method: %s\n- Target: %s\n- Version: %s",
				reqLine.Method, reqLine.RequestTarget, reqLine.HttpVersion)
		}(conn)
	}
}
