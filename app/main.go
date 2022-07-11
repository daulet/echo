package main

import (
	"flag"

	"github.com/daulet/echo"
)

func main() {
	port := flag.Int("port", 9000, "port to listen on")
	flag.Parse()

	echo.Run(*port)
}
