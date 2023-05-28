package client

import (
	"bufio"
	"fmt"
	"net"
	"strconv"

	"github.com/hashicorp/yamux"
	"github.com/morrisonwill/gigatunnl/pkg"
)

type Client struct {
	serverAddress string
	localAddress  string
	serverPort    int
	localPort     int
	session       *yamux.Session
	endUserPort   int
}

func NewClient(serverAddr string, localAddr string, serverPort int, localPort int) *Client {
	return &Client{
		serverAddress: serverAddr,
		localAddress:  localAddr,
		serverPort:    serverPort,
		localPort:     localPort,
	}
}

func (c *Client) Connect() error {
	// Connect to the server
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", c.serverAddress, c.serverPort))
	if err != nil {
		return fmt.Errorf("failed to connect to server: %w", err)
	}
	session, err := yamux.Client(conn, nil)
	if err != nil {
		return fmt.Errorf("failed to create yamux session: %w", err)
	}

	c.session = session

	// Get end user port from the server
	reader := bufio.NewReader(conn)
	line, _, err := reader.ReadLine()
	if err != nil {
		return fmt.Errorf("failed to read end user port: %w", err)
	}
	endUserPort, err := strconv.Atoi(string(line))
	if err != nil {
		return fmt.Errorf("invalid end user port: %w", err)
	}
	c.endUserPort = endUserPort
	fmt.Printf("End user port on server: %d\n", endUserPort)

	return nil
}

func (c *Client) Start() error {
	// Accept new yamux streams and forward them to the local port
	for {
		newStream, err := c.session.Accept()
		if err != nil {
			return fmt.Errorf("failed to accept new stream: %w", err)
		}
		go func() {
			newLocalConnection, err := net.Dial("tcp", fmt.Sprintf("%s:%d", c.localAddress, c.localPort))
			if err != nil {
				fmt.Printf("failed to connect to local port: %v\n", err)
				return
			}
			pkg.Proxy(newStream, newLocalConnection)
		}()
	}
}

func (c *Client) Close() {
	c.session.Close()
}
