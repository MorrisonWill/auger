package server

import (
	"fmt"
	"log"
	"net"
	"strconv"
	"sync"

	"github.com/hashicorp/yamux"

	"github.com/morrisonwill/tunnel/pkg"
)

type Server struct {
	listener  net.Listener
	listeners sync.Map
	ports     sync.Map
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

func (s *Server) SetPorts(ports []string) {
	for _, port := range ports {
		p, err := strconv.Atoi(port)
		if err != nil {
			log.Printf("Invalid port: %v\n", port)
			continue
		}
		s.ports.Store(p, true)
	}
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

	var endUserListener net.Listener
	var err error

	s.ports.Range(func(key, value interface{}) bool {
		port, _ := key.(int)
		endUserListener, err = net.Listen("tcp", fmt.Sprintf(":%d", port))
		if err == nil {
			s.ports.Delete(key)
			return false
		}
		return true
	})

	if err != nil {
		endUserListener, err = net.Listen("tcp", ":0") // 0 lets the system pick an available port
		if err != nil {
			log.Printf("Failed to listen for end users: %v\n", err)
			return
		}
	}

	defer endUserListener.Close()

	s.listeners.Store(endUserListener.Addr().(*net.TCPAddr).Port, endUserListener)

	fmt.Fprintf(clientConn, "%d\n", endUserListener.Addr().(*net.TCPAddr).Port)

	session, err := yamux.Server(clientConn, nil)
	if err != nil {
		log.Printf("Failed to create yamux session: %v\n", err)
		return
	}

	for {
		// Accept an end user connection
		endUserConn, err := endUserListener.Accept()
		if err != nil {
			log.Printf("Failed to accept end user connection: %v\n", err)
			continue
		}

		go func() {
			stream, err := session.Open()

			if err != nil {
				log.Printf("Failed to accept end user connection: %v\n", err)
				return
			}

			log.Println("Accepted end user connection:", endUserConn.RemoteAddr())

			// Start a proxy between the client and the end user
			pkg.Proxy(stream, endUserConn)
		}()
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
