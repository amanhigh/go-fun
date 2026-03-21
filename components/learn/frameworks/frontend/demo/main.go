package main

import (
	"log"
)

func main() {
	// Create and start the UI demo server
	server := NewUIServer("9090")
	log.Println("Starting Templ UI Demo Server on http://localhost:9090")
	log.Println("Visit http://localhost:9090 to see the component demos")

	if err := server.Start(); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
