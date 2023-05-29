package main

import (
	"fmt"
	"net"
	"os"

	"github.com/hashicorp/yamux"
)

func main() {
	if os.Args[1] == "server" {
		fmt.Println("server")
		server()
	} else {
		fmt.Println("client")
		client()
	}
}

func server() {
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		panic(err)
	}

	// Accept a TCP connection
	conn, err := listener.Accept()
	if err != nil {
		panic(err)
	}

	// Setup server side of yamux
	session, err := yamux.Server(conn, nil)
	if err != nil {
		panic(err)
	}

	// Accept a stream
	stream, err := session.Open()
	if err != nil {
		panic(err)
	}

	// Stream implements net.Conn
	stream.Write([]byte("ping"))
}

func client() {
	// Get a TCP connection
	conn, err := net.Dial("tcp", "tnl.pub:8080")
	if err != nil {
		panic(err)
	}

	// Setup client side of yamux
	session, err := yamux.Client(conn, nil)
	if err != nil {
		panic(err)
	}

	// Open a new stream
	stream, err := session.Accept()
	if err != nil {
		panic(err)
	}

	// Listen for a message
	buf := make([]byte, 4)
	stream.Read(buf)
	fmt.Println(string(buf))
}
