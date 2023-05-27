package server

import (
	"fmt"
	"log"
	"net"
	"sync"

	"github.com/morrisonwill/gigatunnl/pkg"
)

type Server struct {
	listener  net.Listener
	listeners sync.Map
}

func NewServer(address string) (*Server, error) {
	listener, err := net.Listen("tcp", address)
	if err != nil {
		return nil, fmt.Errorf("failed to listen: %w", err)
	}
	return &Server{
		listener: listener,
	}, nil
}

func (s *Server) Start() {
	log.Printf("Server listening on %s\n", s.listener.Addr().String())
	for {
		clientConn, err := s.listener.Accept()
		if err != nil {
			log.Printf("Failed to accept client connection: %v\n", err)
			continue
		}
		log.Println("Accepted client connection:", clientConn.RemoteAddr())

		go s.handleClient(clientConn)
	}
}

func (s *Server) handleClient(clientConn net.Conn) {
	defer clientConn.Close()

	// Listen for end users
	endUserListener, err := net.Listen("tcp", ":0") // 0 lets the system pick an available port
	if err != nil {
		log.Printf("Failed to listen for end users: %v\n", err)
		return
	}
	defer endUserListener.Close()

	s.listeners.Store(endUserListener.Addr().(*net.TCPAddr).Port, endUserListener)

	fmt.Fprintf(clientConn, "%d\n", endUserListener.Addr().(*net.TCPAddr).Port)

	for {
		// Accept an end user connection
		endUserConn, err := endUserListener.Accept()
		if err != nil {
			log.Printf("Failed to accept end user connection: %v\n", err)
			continue
		}

		log.Println("Accepted end user connection:", endUserConn.RemoteAddr())

		// Start a proxy between the client and the end user
		proxy := pkg.NewProxy(clientConn, endUserConn)
		go proxy.StartProxy()
	}
}

func (s *Server) Close() {
	s.listener.Close()

	s.listeners.Range(func(key, value interface{}) bool {
		if listener, ok := value.(net.Listener); ok {
			listener.Close()
		}
		return true
	})
}
