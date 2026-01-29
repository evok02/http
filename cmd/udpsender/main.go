package main

import (
	"bufio"
	"log"
	"net"
	"os"
)

func main() {
	addr, err := net.ResolveUDPAddr("udp", ":42069")
	if err != nil {
		log.Fatal(err)
	}

	conn, err := net.DialUDP("udp", &net.UDPAddr{}, addr)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	r := bufio.NewReader(os.Stdin)
	for {
		print(">")
		input, err := r.ReadString('-')
		if err != nil {
			log.Println(err.Error())
		}
		conn.Write([]byte(input))
	}

}
