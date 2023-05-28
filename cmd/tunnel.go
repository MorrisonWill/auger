package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/morrisonwill/tunnel/client"
	"github.com/morrisonwill/tunnel/server"
	"github.com/spf13/cobra"
)

var (
	serverPort          int
	localAddress        string
	serverAddress       string
	minPort, maxPort    int
	commaSeparatedPorts string
)

var rootCmd = &cobra.Command{
	Use:   "tunnel",
	Short: "Tunnel is an open-source alternative to ngrok",
	Long: `A Fast and Flexible tunneling tool written in Go.
Complete documentation is available at https://github.com/morrisonwill/tunnel`,
	Run: func(cmd *cobra.Command, args []string) {
		_ = cmd.Help()
	},
}

var localCmd = &cobra.Command{
	Use:   "local [localPort]",
	Short: "Starts a local proxy to the remote server",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		localPort, err := strconv.Atoi(args[0])
		if err != nil {
			log.Fatalf("Invalid local port: %v", err)
		}
		runClient(serverAddress, localPort)
	},
}

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Runs the remote proxy server",
	Run: func(cmd *cobra.Command, args []string) {
		runServer(serverPort)
	},
}

func init() {
	localCmd.Flags().StringVarP(&serverAddress, "remote-address", "r", "localhost", "address of the server to connect to")
	localCmd.Flags().IntVarP(&serverPort, "remote-port", "p", 49152, "server's control port")

	serverCmd.Flags().IntVar(&serverPort, "port", 49152, "control port on the server")
	serverCmd.Flags().IntVar(&minPort, "min-port", 0, "Minimum port range")
	serverCmd.Flags().IntVar(&maxPort, "max-port", 0, "Maximum port range")
	serverCmd.Flags().StringVar(&commaSeparatedPorts, "ports", "", "Comma-separated ports")

	rootCmd.AddCommand(localCmd)
	rootCmd.AddCommand(serverCmd)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

// runServer runs the server mode.
func runServer(port int) {
	server, err := server.NewServer(fmt.Sprintf("0.0.0.0:%d", port))

	if err != nil {
		log.Fatalf("Failed to create server: %v", err)
	}

	// Set Port Range
	server.SetPortRange(minPort, maxPort)

	// Set Specific Ports
	if commaSeparatedPorts != "" {
		ports := strings.Split(commaSeparatedPorts, ",")
		server.SetPorts(ports)
	}

	// Start the server
	server.Start()

	fmt.Println("Server is running. Press Ctrl+C to exit.")

	// Block indefinitely to keep the server running
	select {}
}

// runClient runs the client mode.
func runClient(serverAddress string, localPort int) {
	client := client.NewClient(fmt.Sprintf("%s:%d", serverAddress, serverPort), fmt.Sprintf("localhost:%d", localPort))

	err := client.Connect()
	if err != nil {
		log.Fatalf("Failed to connect to server: %v", err)
	}

	err = client.Start()
	if err != nil {
		log.Fatalf("Failed to start client: %v", err)
	}

	fmt.Println("Client is running. Press Ctrl+C to exit.")
	// Block indefinitely to keep the client running
	select {}
}
