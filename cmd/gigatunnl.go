package main

import (
	"flag"
	"fmt"
	"log"

	"gigatunnl/client"
	"gigatunnl/server"
)

func main() {
	// Parse command-line arguments
	mode := flag.String("mode", "", "The mode to run: 'server' or 'client'")
	flag.Parse()

	// Run the appropriate mode
	switch *mode {
	case "server":
		runServer()
	case "client":
		runClient()
	default:
		log.Fatal("Invalid mode specified. Use 'server' or 'client'.")
	}
}

// runServer runs the server mode.
func runServer() {
	address := "localhost"
	port := 3000

	server, err := server.NewServer(fmt.Sprintf("%s:%d", address, port))

	if err != nil {
		log.Fatalf("Failed to create server: %v", err)
	}

	// Start the server
	server.Start()

	fmt.Println("Server is running. Press Ctrl+C to exit.")
	// Block indefinitely to keep the server running
	select {}
}

// runClient runs the client mode.
func runClient() {

	client := client.NewClient("localhost", "localhost", 3000, 8000)

	err := client.Connect()
	if err != nil {
		log.Fatalf("Failed to connect to server: %v", err)
	}

	err = client.Start()
	if err != nil {
		log.Fatalf("Failed to start client: %v", err)
	}

	fmt.Println("Client is running. Press Ctrl+C to exit.")
	// Block indefinitely to keep the client running
	select {}
}
