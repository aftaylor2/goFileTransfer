package main

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"os"
)

func main() {
	if len(os.Args) != 4 {
		fmt.Println("Usage: sender <server_ip> <session_id> <file_path>")
		os.Exit(1)
	}

	serverIP := os.Args[1]
	sessionID := os.Args[2]
	filePath := os.Args[3]

	conn, err := net.Dial("tcp", serverIP+":8080")
	if err != nil {
		fmt.Printf("Error connecting to server: %v\n", err)
		os.Exit(1)
	}
	defer conn.Close()

	// Send client type and session ID
	_, err = fmt.Fprintf(conn, "SENDER\n%s\n", sessionID)
	if err != nil {
		fmt.Printf("Error sending session ID: %v\n", err)
		os.Exit(1)
	}

	// Wait for READY signal
	scanner := bufio.NewScanner(conn)
	if scanner.Scan() {
		response := scanner.Text()
		if response != "READY" {
			fmt.Printf("Error: %s\n", response)
			os.Exit(1)
		}
		fmt.Println("Receiver is ready. Starting file transfer...")
	} else {
		fmt.Println("No response from server.")
		os.Exit(1)
	}

	// Open the file to send
	file, err := os.Open(filePath)
	if err != nil {
		fmt.Printf("Error opening file: %v\n", err)
		os.Exit(1)
	}
	defer file.Close()

	fmt.Printf("Sending file: %s\n", filePath)
	buffer := make([]byte, 4096)
	totalBytesSent := 0

	for {
		n, err := file.Read(buffer)
		if n > 0 {
			_, writeErr := conn.Write(buffer[:n])
			if writeErr != nil {
				fmt.Printf("Error sending data: %v\n", writeErr)
				break
			}
			totalBytesSent += n
			fmt.Printf("Sent %d bytes to server.\n", n)
		}
		if err == io.EOF {
			break
		}
		if err != nil {
			fmt.Printf("Error reading file: %v\n", err)
			break
		}
	}

	fmt.Printf("File sent successfully (%d bytes).\n", totalBytesSent)
}
