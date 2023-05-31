package server

import (
	"errors"
	"fmt"

	"math/rand"
	"net"
	"strconv"
	"sync"
	"time"

	"github.com/charmbracelet/log"

	"github.com/hashicorp/yamux"

	"github.com/morrisonwill/auger/proxy"
)

type Server struct {
	listener net.Listener
	rand     *rand.Rand
	ports    ports
	address  string
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

	// TODO consider moving things here or in Start
	return &Server{
		listener: listener,
		rand:     rand.New(rand.NewSource(time.Now().UnixNano())),
		ports: ports{
			list: nil,
		},
		address: address,
	}, nil
}

func (s *Server) SetPortRange(start int, end int) {
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
			log.Infof("Invalid port: %v\n", port)
			continue
		}
		s.ports.Lock()
		s.ports.list = append(s.ports.list, p)
		s.ports.Unlock()
	}
}

func (s *Server) Start() {
	log.Infof("Server is listening on %s", s.address)

	for {
		clientConn, err := s.listener.Accept()
		if err != nil {
			log.Errorf("Failed to accept client connection: %v\n", err)
			continue
		}
		log.Infof("New client: %s", clientConn.RemoteAddr())

		go s.handleClient(clientConn)
	}
}

func (s *Server) handleClient(clientConn net.Conn) {

	defer clientConn.Close()

	var endUserListener net.Listener

	var err error

	var endUserPort int

	s.ports.Lock()
	if s.ports.list == nil {
		endUserListener, err = net.Listen("tcp", ":0") // 0 lets the system pick an available port
		if err != nil {
			log.Errorf("Failed to listen for end users %v", err)
			return
		}
		endUserPort = endUserListener.Addr().(*net.TCPAddr).Port
	} else {
		randPortIdx := s.rand.Intn(len(s.ports.list))
		endUserPort = s.ports.list[randPortIdx]
		endUserListener, err = net.Listen("tcp", fmt.Sprintf(":%d", endUserPort))
		if err != nil {
			log.Errorf("Failed to listen for end users %v", err)
			return
		}
		s.ports.list = append(s.ports.list[:randPortIdx], s.ports.list[randPortIdx+1:]...)
	}
	s.ports.Unlock()

	defer endUserListener.Close()

	fmt.Fprintf(clientConn, "%d\n", endUserListener.Addr().(*net.TCPAddr).Port)

	// Define Yamux config
	config := yamux.DefaultConfig()
	// Enable keepalives
	config.KeepAliveInterval = 30 * time.Second

	session, err := yamux.Server(clientConn, config)

	if err != nil {
		log.Errorf("Failed to create session with client: %v\n", err)
		return
	}

	// check if client is still alive
	go func() {
		for {
			_, err := session.Ping()
			if err != nil {
				endUserListener.Close()
				session.Close()
				s.ports.Lock()
				s.ports.list = append(s.ports.list, endUserPort)
				s.ports.Unlock()
				return
			}
			time.Sleep(time.Second * 30)
		}
	}()

	for {
		endUserConn, err := endUserListener.Accept()
		if err != nil {

			if errors.Is(err, net.ErrClosed) {
				log.Infof("End user connection closed: %v\n", endUserConn.RemoteAddr())
				break
			}
			log.Errorf("Failed to accept end user connection: %v\n", err)
			continue
		}

		stream, err := session.Open()
		if err != nil {
			log.Info("Client disconnected:" + clientConn.RemoteAddr().String())
			endUserConn.Close()

			// add port back to list
			s.ports.Lock()
			s.ports.list = append(s.ports.list, endUserPort)
			s.ports.Unlock()
			return
		}

		// Accept an end user connection
		go func() {
			// Start a proxy between the client and the end user
			proxy.Proxy(stream, endUserConn)
		}()
	}

}
