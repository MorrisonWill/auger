package server

import (
	"fmt"
	"log"
	"math/rand"
	"net"
	"strconv"
	"sync"
	"time"

	"github.com/hashicorp/yamux"

	"github.com/morrisonwill/tunnel/pkg"
)

type Server struct {
	listener net.Listener
	rand     *rand.Rand
	ports    ports
}

type ports struct {
	sync.Mutex
	list []int
}

func NewServer(address string) (*Server, error) {
	listener, err := net.Listen("tcp", address)
	if err != nil {
		return nil, fmt.Errorf("failed to listen: %w", err)
	}

	return &Server{
		listener: listener,
		rand:     rand.New(rand.NewSource(time.Now().UnixNano())),
		ports: ports{
			list: nil,
		},
	}, nil
}

func (s *Server) SetPortRange(start int, end int) {
	fmt.Println("ports before", s.ports.list)

	for port := start; port <= end; port++ {
		s.ports.Lock()
		s.ports.list = append(s.ports.list, port)
		s.ports.Unlock()
	}
}

func (s *Server) SetPorts(ports []string) {
	for _, port := range ports {
		p, err := strconv.Atoi(port)
		if err != nil {
			log.Printf("Invalid port: %v\n", port)
			continue
		}
		s.ports.Lock()
		s.ports.list = append(s.ports.list, p)
		s.ports.Unlock()
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
	doneChan := make(chan bool)

	defer clientConn.Close()

	var endUserListener net.Listener
	var err error

	var endUserPort int

	s.ports.Lock()
	if s.ports.list == nil {
		endUserListener, err = net.Listen("tcp", ":0") // 0 lets the system pick an available port
		endUserPort = endUserListener.Addr().(*net.TCPAddr).Port
	} else {
		randPortIdx := s.rand.Intn(len(s.ports.list))
		endUserPort = s.ports.list[randPortIdx]
		endUserListener, err = net.Listen("tcp", fmt.Sprintf(":%d", endUserPort))
		if err == nil {
			// remove port from list
			s.ports.list = append(s.ports.list[:randPortIdx], s.ports.list[randPortIdx+1:]...)
		}
	}

	s.ports.Unlock()

	if err != nil {
		log.Printf("Failed to listen for end users: %v\n", err)
		return
	}

	defer endUserListener.Close()

	fmt.Fprintf(clientConn, "%d\n", endUserListener.Addr().(*net.TCPAddr).Port)

	session, err := yamux.Server(clientConn, nil)
	if err != nil {
		log.Printf("Failed to create yamux session: %v\n", err)
		return
	}

	// check if client is still alive
	go func() {
		for {
			_, err := session.Ping()
			if err != nil {
				log.Printf("Client disconnected")
				endUserListener.Close()
				s.ports.Lock()
				s.ports.list = append(s.ports.list, endUserPort)
				s.ports.Unlock()
				doneChan <- true
				return
			}
			time.Sleep(time.Second * 10)
		}
	}()

	for {
		if <-doneChan {
			log.Printf("CLI disconnected, killing proxy")
			return
		}

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
