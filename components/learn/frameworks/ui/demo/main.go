package main

import (
	"log"
)

func main() {
	// Create and start the UI demo server
	server := NewUIServer("8080")
	log.Println("Starting Templ UI Demo Server on http://localhost:8080")
	log.Println("Visit http://localhost:8080 to see the component demos")

	if err := server.Start(); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
