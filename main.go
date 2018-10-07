package main

import (
	"flag"

	"./punchbag"
)

func main() {
	port := flag.Int("port", 9000, "port to listen on")
	flag.Parse()

	punchbag.Run(*port)
}
