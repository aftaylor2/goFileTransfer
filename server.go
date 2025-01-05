package main

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"strconv"
	"sync"
)

const maxSessions = 100

type session struct {
	receiver        net.Conn
	sender          net.Conn
	mux             sync.Mutex
	readyChan       chan struct{} // Channel to signal readiness
	isReceiverReady bool          // Tracks receiver readiness
}

var sessions [maxSessions]session

func main() {
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		fmt.Printf("Error starting server: %v\n", err)
		return
	}
	defer listener.Close()

	fmt.Println("Server listening on port 8080...")
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Printf("Error accepting connection: %v\n", err)
			continue
		}
		go handleClient(conn)
	}
}

func handleClient(conn net.Conn) {
	defer conn.Close()
	reader := bufio.NewReader(conn)

	// Identify client type
	clientType, err := reader.ReadString('\n')
	if err != nil {
		fmt.Printf("Error reading client type: %v\n", err)
		return
	}
	clientType = clientType[:len(clientType)-1]

	switch clientType {
	case "RECEIVER":
		handleReceiver(conn)
	case "SENDER":
		handleSender(conn, reader)
	default:
		fmt.Printf("Unknown client type: %s\n", clientType)
	}
}

func handleReceiver(conn net.Conn) {
	for i := 0; i < maxSessions; i++ {
		session := &sessions[i]
		session.mux.Lock()
		if session.receiver == nil {
			session.receiver = conn
			session.readyChan = make(chan struct{}, 1)
			session.isReceiverReady = true
			session.mux.Unlock()

			sessionID := strconv.Itoa(i)
			fmt.Fprintf(conn, "%s\n", sessionID)
			fmt.Printf("Receiver connected. Session ID: %s\n", sessionID)
			select {} // Block until sender finishes
		}
		session.mux.Unlock()
	}
	fmt.Fprintln(conn, "Server full")
}

func handleSender(conn net.Conn, reader *bufio.Reader) {
	sessionIDStr, err := reader.ReadString('\n')
	if err != nil {
		fmt.Printf("Error reading session ID: %v\n", err)
		return
	}
	sessionID, err := strconv.Atoi(sessionIDStr[:len(sessionIDStr)-1])
	if err != nil || sessionID < 0 || sessionID >= maxSessions {
		fmt.Fprintln(conn, "Invalid session ID")
		fmt.Printf("Invalid session ID: %s\n", sessionIDStr)
		return
	}

	session := &sessions[sessionID]
	session.mux.Lock()
	if session.receiver == nil || !session.isReceiverReady {
		session.mux.Unlock()
		fmt.Fprintln(conn, "No receiver connected")
		fmt.Printf("No receiver connected for Session ID: %d\n", sessionID)
		return
	}
	session.sender = conn
	session.mux.Unlock()

	// Notify the sender that the receiver is ready
	session.readyChan <- struct{}{}
	close(session.readyChan) // Ensure the channel is closed after sending
	fmt.Fprintln(conn, "READY")
	fmt.Printf("Sender connected. Session ID: %d. Starting relay...\n", sessionID)

	buffer := make([]byte, 4096)
	for {
		n, err := conn.Read(buffer)
		if n > 0 {
			session.mux.Lock()
			_, writeErr := session.receiver.Write(buffer[:n])
			session.mux.Unlock()
			if writeErr != nil {
				fmt.Printf("Error relaying data to receiver (ID: %d): %v\n", sessionID, writeErr)
				break
			}
			fmt.Printf("Relayed %d bytes to receiver (ID: %d).\n", n, sessionID)
		}
		if err == io.EOF {
			fmt.Printf("Sender finished sending data.  (ID: %d)\n", sessionID)
			break
		}
		if err != nil {
			fmt.Printf("Error reading from sender (ID: %d): %v\n", sessionID, err)
			break
		}
	}

	// Clean up session
	session.mux.Lock()
	session.receiver.Close()
	session.receiver = nil
	session.sender = nil
	session.isReceiverReady = false
	session.mux.Unlock()
	fmt.Printf("Relay complete for Session ID: %d.\n", sessionID)
}
