package main

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"os"
)

func main() {
	serverIP := "127.0.0.1"

	conn, err := net.Dial("tcp", serverIP+":8080")
	if err != nil {
		fmt.Printf("Error connecting to server: %v\n", err)
		os.Exit(1)
	}
	defer conn.Close()

	// Send client type
	_, err = fmt.Fprint(conn, "RECEIVER\n")
	if err != nil {
		fmt.Printf("Error sending client type: %v\n", err)
		os.Exit(1)
	}

	// Receive session ID
	scanner := bufio.NewScanner(conn)
	if scanner.Scan() {
		sessionID := scanner.Text()
		fmt.Printf("Session ID: %s\n", sessionID)
	} else {
		fmt.Println("Failed to receive session ID.")
		os.Exit(1)
	}

	// Wait for file transfer
	fmt.Println("Waiting for file transfer to start...")
	var outputFile *os.File
	totalBytesReceived := 0
	buffer := make([]byte, 4096)

	for {
		n, err := conn.Read(buffer)
		if n > 0 {
			// Create the file only when data is received
			if outputFile == nil {
				outputFile, err = os.Create("received_file")
				if err != nil {
					fmt.Printf("Error creating file: %v\n", err)
					os.Exit(1)
				}
				defer outputFile.Close()
				fmt.Println("Transfer started. Writing to 'received_file'...")
			}

			_, writeErr := outputFile.Write(buffer[:n])
			if writeErr != nil {
				fmt.Printf("Error writing to file: %v\n", writeErr)
				break
			}
			totalBytesReceived += n
			fmt.Printf("Received %d bytes.\n", n)
		}
		if err == io.EOF {
			fmt.Println("Connection closed by server.")
			break
		}
		if err != nil {
			fmt.Printf("Error reading from connection: %v\n", err)
			break
		}
	}

	if totalBytesReceived > 0 {
		fmt.Printf("Transfer finished. Total bytes received: %d\n", totalBytesReceived)
	} else {
		fmt.Println("No data received. Transfer aborted.")
	}
}
