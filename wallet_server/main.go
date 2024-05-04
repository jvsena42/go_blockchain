package main

import (
	"flag"
	"log"
)

func Init() {
	log.SetPrefix("Wallet Server: ")
}

func main() {
	port := flag.Uint("port", 8080, "TCP Port number for online wallet")
	gateway := flag.String("gateway", "http://127.0.0.1:3333", "TCP Port number for online wallet")
	flag.Parse()

	app := NewWalletServer(uint16(*port), *gateway)
	log.Println("Starting server on port:", *port, "Using blockchain node", *gateway, " as gateway")
	app.Run()
}
