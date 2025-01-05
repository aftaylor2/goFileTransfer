
# Go File Transfer Project

This project is a work in progress that I'm using to practice the Go programming language. It is not intended for production use at this time. I plan to write the same project in Go, C, and Zig. The project implements a simple file transfer system using Go. It consists of three components:

1. **Server**: Acts as an intermediary to relay data from the sender to the receiver.
2. **Sender**: Sends a file to a receiver via the server.
3. **Receiver**: Receives a file from a sender via the server.

Allows for file transfers where direct connections are not possible due to NAT or firewall issues. It also allows for one party to send files to another without knowing the IP address or identity of the other party. All that is required is the session ID.

## Features

- **Intermediary Server**: Relays data between the sender and receiver without storing it.
- **Session Management**: Uses unique session IDs to pair senders and receivers.
- **Synchronization**: Ensures the sender waits until the receiver is connected before starting the transfer.
- **Streaming**: Streams data directly from the sender to the receiver.

## Project Structure

```plaintext
.
├── server.go    # Intermediary server implementation
├── sender.go    # Sender implementation
├── receiver.go  # Receiver implementation
├── Makefile     # Compile / Build the project
├── README.md    # Project documentation
└── bin/         # Compiled Binaries
```

## Setup

1. Build the binaries:

   ```bash
   make
   ```

## Usage

### Start the Server

The server listens on port `8080` by default. Start it first:

```bash
./server
```

### Start the Receiver

The receiver connects to the server and waits for a file transfer:

```bash
./receiver
```

After running the receiver, it will display a session ID (e.g., `Session ID: 0`).

### Start the Sender

The sender connects to the server and sends a file to the receiver using the session ID:

```bash
./sender <server_ip> <session_id> <file_path>
```

Example:

```bash
./sender 127.0.0.1 0 /path/to/file
```

### File Transfer

- The server relays data from the sender to the receiver.
- The receiver writes the data to a file named `received_file` in the current directory.

## Troubleshooting

- Ensure the server is running prior to starting the receiver.
- Ensure the receiver is running prior to starting the sender.
- Ensure the session ID provided to the sender matches the one displayed by the receiver.
- Use `netstat`, `tcpdump` or a similar tool to verify that connections are established between components.

## Contributing

Contributions are welcome! Please fork the repository, make changes, and submit a pull request.

### Prerequisites

- Go 1.17 or later installed on your system.
- Basic understanding of networking concepts and Go programming.

## License

This project is licensed under the MIT License. See the `LICENSE` file for details.

## TODO

1. Use CLI arguments to specify which IP address to bind to.
2. Save the received file, retaining the original file name.
3. Add TLS encryption.
4. Add an optional session password.
5. Consider building the go binaries in a Go specific fashion using distinct directories for each binary.  
   This will prevent the "main redeclared" LSP errors.
