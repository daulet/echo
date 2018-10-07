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

	var channelID int
	for {
		conn, err := l.Accept()
		if err != nil {
			log.Fatalf("Accept(): %v", err)
		}
		go handle(channelID, conn)
		channelID++
	}
}

func handle(id int, conn net.Conn) {
	defer conn.Close()

	for {
		var buf [128]byte
		n, err := conn.Read(buf[:])
		if err != nil {
			log.Printf("channel %v: Read() failed with: %v\n", id, err)
			return
		}
		log.Printf("channel %v: %v", id, string(buf[:n]))
		_, err = conn.Write([]byte(fmt.Sprintf("read %v bytes\n", n)))
		if err != nil {
			log.Printf("channel %v: Write() failed with: %v\n", id, err)
		}
	}
}
