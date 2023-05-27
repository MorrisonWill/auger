package pkg

import (
	"io"
	"log"
	"net"
)

// Proxy handles bidirectional proxying of TCP connections.
type Proxy struct {
	sourceConn      net.Conn
	destinationConn net.Conn
}

// NewProxy creates a new instance of the Proxy.
func NewProxy(sourceConn net.Conn, destinationConn net.Conn) *Proxy {
	return &Proxy{
		sourceConn:      sourceConn,
		destinationConn: destinationConn,
	}
}

// StartProxy starts the bidirectional proxy.
func (p *Proxy) StartProxy() {
	go p.proxyData(p.sourceConn, p.destinationConn)
	go p.proxyData(p.destinationConn, p.sourceConn)
}

// proxyData proxies data between the source and destination connections.
func (p *Proxy) proxyData(sourceConn net.Conn, destinationConn net.Conn) {
	defer sourceConn.Close()
	defer destinationConn.Close()

	_, err := io.Copy(destinationConn, sourceConn)
	if err != nil {
		log.Println("Error copying data from source to destination:", err.Error())
	}
	log.Println("Finished proxying data")
}
