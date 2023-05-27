package client

import (
	"bufio"
	"fmt"
	"gigatunnl/pkg"
	"net"
	"strconv"
)

type Client struct {
	serverAddr  string
	localAddr   string
	serverPort  int
	localPort   int
	conn        net.Conn
	localConn   net.Conn
	endUserPort int
}

func NewClient(serverAddr string, localAddr string, serverPort int, localPort int) *Client {
	return &Client{
		serverAddr: serverAddr,
		localAddr:  localAddr,
		serverPort: serverPort,
		localPort:  localPort,
	}
}

func (c *Client) Connect() error {
	// Connect to the server
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", c.serverAddr, c.serverPort))
	if err != nil {
		return fmt.Errorf("failed to connect to server: %w", err)
	}
	c.conn = conn

	// Get end user port from the server
	reader := bufio.NewReader(c.conn)
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

	// Connect to the local service
	localConn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", c.localAddr, c.localPort))
	if err != nil {
		return fmt.Errorf("failed to connect to local service: %w", err)
	}
	c.localConn = localConn

	return nil
}

func (c *Client) Start() error {
	// Start the bidirectional proxy
	proxy := pkg.NewProxy(c.localConn, c.conn)
	proxy.StartProxy()

	return nil
}

func (c *Client) Close() {
	c.conn.Close()
	c.localConn.Close()
}
