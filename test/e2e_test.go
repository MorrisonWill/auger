package tests

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"testing"
	"time"

	"github.com/morrisonwill/tunnel/client"
	"github.com/morrisonwill/tunnel/server"
)

func TestEndToEnd(t *testing.T) {
	// Step 1: Start local server
	startLocalServer()

	// Step 2: Start tunnel server
	tunnelServer := startTunnelServer()

	// Step 3: Start tunnel client
	tunnelClient := startTunnelClient(tunnelServer)

	// Step 4: Make a request through the tunnel
	resp, err := http.Get(fmt.Sprintf("http://localhost:%d", tunnelClient.EndUserPort))
	if err != nil {
		t.Fatalf("Failed to make request through the tunnel: %v", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	if string(body) != "Hello, world!" {
		t.Fatalf("Unexpected response: %s", body)
	}
}

func startLocalServer() *http.Server {
	localServer := &http.Server{Addr: ":8080", Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "Hello, world!")
	})}
	go func() {
		if err := localServer.ListenAndServe(); err != nil {
			log.Fatalf("Failed to start local server: %v", err)
		}
	}()
	return localServer
}

func startTunnelServer() *server.Server {
	tunnelServer, err := server.NewServer("localhost:49152")
	if err != nil {
		log.Fatalf("Failed to create server: %v", err)
	}

	go tunnelServer.Start()
	time.Sleep(1 * time.Second)
	return tunnelServer
}

func startTunnelClient(tunnelServer *server.Server) *client.Client {
	tunnelClient := client.NewClient("localhost:49152", "localhost:8080")
	err := tunnelClient.Connect()
	if err != nil {
		log.Fatalf("Failed to connect to server: %v", err)
	}

	go func() {
		if err := tunnelClient.Start(); err != nil {
			log.Fatalf("Failed to start tunnel client: %v", err)
		}
	}()
	time.Sleep(1 * time.Second)
	return tunnelClient
}
