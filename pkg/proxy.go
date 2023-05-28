package pkg

import (
	"io"
	"log"
	"net"
	"sync"
)

// Proxy starts the bidirectional proxy between sourceConn and destinationConn
func Proxy(sourceConn net.Conn, destinationConn net.Conn) {
	var wg sync.WaitGroup
	wg.Add(2) // We will wait for two goroutines

	go proxyData(&wg, sourceConn, destinationConn)
	go proxyData(&wg, destinationConn, sourceConn)

	// Wait for both goroutines to finish then close connections
	go func() {
		wg.Wait()
		sourceConn.Close()
		destinationConn.Close()
	}()
}

// proxyData proxies data between the source and destination connections.
func proxyData(wg *sync.WaitGroup, destinationConn net.Conn, sourceConn net.Conn) {
	defer wg.Done()

	_, err := io.Copy(destinationConn, sourceConn)
	if err != nil {
		log.Println("Error copying data from source to destination:", err.Error())
	}
	log.Println("Finished proxying data")
}
