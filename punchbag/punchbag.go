package punchbag

import (
	"fmt"
	"log"
	"net"
)

// Run accepts all messages on specified port.
func Run(port int) {
	address := fmt.Sprintf("localhost:%v", port)
	l, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatalf("Listen(%v): %v", address, err)
	}

	for {
		conn, err := l.Accept()
		if err != nil {
			log.Fatalf("Accept(): %v", err)
		}
		go handle(conn)
	}
}

func handle(conn net.Conn) {
	defer conn.Close()

	for {
		var buf [128]byte
		n, err := conn.Read(buf[:])
		if err != nil {
			log.Printf("Read() failed with: %v\n", err)
		}
		_, err = conn.Write([]byte(fmt.Sprintf("read %v bytes\n", n)))
		if err != nil {
			log.Printf("Write() failed with: %v\n", err)
		}
	}
}
