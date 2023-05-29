package client

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/charmbracelet/log"
	"github.com/hashicorp/yamux"
	"github.com/morrisonwill/tunnel/proxy"
)

type Client struct {
	serverAddress string
	localAddress  string
	EndUserPort   int
	session       *yamux.Session
}

func NewClient(serverAddr string, localAddr string) *Client {
	return &Client{
		serverAddress: serverAddr,
		localAddress:  localAddr,
	}
}

func (c *Client) Connect() error {
	// Connect to the server
	conn, err := net.Dial("tcp", c.serverAddress)
	if err != nil {
		return err
	}

	// Define Yamux config
	config := yamux.DefaultConfig()
	// Enable keepalives
	config.KeepAliveInterval = 30 * time.Second

	session, err := yamux.Client(conn, config)
	if err != nil {
		return err
	}

	c.session = session

	// Get end user port from the server
	reader := bufio.NewReader(conn)
	line, _, err := reader.ReadLine()
	if err != nil {
		return err
	}
	endUserPort, err := strconv.Atoi(string(line))
	if err != nil {
		return err
	}
	c.EndUserPort = endUserPort
	log.Infof("Listening on %s:%d\n", c.serverAddress, endUserPort)

	return nil
}

func (c *Client) Start() error {
	// Intercept sigint (ctrl-c)
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-interrupt
		c.session.Close()
		os.Exit(0)
	}()

	// Accept new yamux streams and forward them to the local port
	for {
		fmt.Println("Waiting for new stream")
		newStream, err := c.session.Accept()
		if err != nil {
			return err
		}
		go func() {
			newLocalConnection, err := net.Dial("tcp", c.localAddress)
			if err != nil {
				log.Errorf("Failed to connect to local port: %v\n", err)
				return
			}
			proxy.Proxy(newStream, newLocalConnection)
		}()
	}
}
