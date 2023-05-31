package client

import (
	"bufio"
	"net"
	"os"
	"os/signal"
	"strconv"
	"strings"
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

	// Get end user port from the server
	// TODO switch to byte[] and binary.Uint16
	// TODO consider doing this over the yamux session.
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

	session, err := yamux.Client(conn, config)
	if err != nil {
		return err
	}

	c.session = session

	addressParts := strings.Split(c.serverAddress, ":")
	hostname := strings.Join(addressParts[:len(addressParts)-1], ":")
	log.Infof("Listening on %s:%d\n", hostname, endUserPort)

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
