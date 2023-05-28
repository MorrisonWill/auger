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

// TODO: fix port selection, currently it'll empty then choose available
// TODO: add graceful shutdown on both sides
// currently one problem with this is detecting when client disconnects and freeing ports

var (
	serverAddress       string
	minPort, maxPort    int
	commaSeparatedPorts string
)

const serverPort = 49152

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
	serverCmd.Flags().IntVar(&minPort, "min-port", getEnvAsInt("TUNNEL_MIN_PORT", 0), "Minimum port range")
	serverCmd.Flags().IntVar(&maxPort, "max-port", getEnvAsInt("TUNNEL_MAX_PORT", 0), "Maximum port range")
	serverCmd.Flags().StringVar(&commaSeparatedPorts, "ports", getEnvAsString("TUNNEL_PORTS", ""), "Comma-separated ports")

	localCmd.Flags().StringVarP(&serverAddress, "remote-address", "r", getEnvAsString("TUNNEL_REMOTE_ADDRESS", ""), "address of the server to connect to")
	err := localCmd.MarkFlagRequired("remote-address")
	if err != nil {
		log.Fatalf("Failed to mark remote-address as required: %v", err)
	}

	rootCmd.AddCommand(localCmd)
	rootCmd.AddCommand(serverCmd)
}

func getEnvAsString(name string, defaultValue string) string {
	value, exists := os.LookupEnv(name)
	if !exists {
		return defaultValue
	}
	return value
}

func getEnvAsInt(name string, defaultValue int) int {
	valueStr, exists := os.LookupEnv(name)
	if !exists {
		return defaultValue
	}

	value, err := strconv.Atoi(valueStr)
	if err != nil {
		log.Printf("Failed to convert %s to integer. Using default value: %v", name, defaultValue)
		return defaultValue
	}

	return value
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

	// Set Port Range if defined
	if minPort != 0 && maxPort != 0 {
		server.SetPortRange(minPort, maxPort)
	}

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
