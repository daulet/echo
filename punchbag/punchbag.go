package punchbag

import (
	"fmt"
	"io"
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

	// TODO this might be suboptimal when there is nothing to write
	// add sleep to avoid busy loop
	_, err := io.Copy(conn, conn)
	if err != nil {
		log.Printf("channel %v: finished with: %v\n", id, err)
		return
	}
}
