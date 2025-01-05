# Compiler
GOCMD = go

# Build Targets
SENDER = ./bin/sender
RECEIVER = ./bin/receiver
SERVER = ./bin/server

# Source Files
SENDER_SRC = sender.go
RECEIVER_SRC = receiver.go
SERVER_SRC = server.go

# Default Target
all: $(SENDER) $(RECEIVER) $(SERVER)

# Build Sender
$(SENDER): $(SENDER_SRC)
	$(GOCMD) build -o $(SENDER) $(SENDER_SRC)

# Build Receiver
$(RECEIVER): $(RECEIVER_SRC)
	$(GOCMD) build -o $(RECEIVER) $(RECEIVER_SRC)

# Build Server
$(SERVER): $(SERVER_SRC)
	$(GOCMD) build -o $(SERVER) $(SERVER_SRC)

# Run Server
run-server:
	./$(SERVER)

# Run Receiver
run-receiver:
	./$(RECEIVER)

# Run Sender
run-sender:
	./$(SENDER) 127.0.0.1 0 /path/to/file

# Clean Up Build Artifacts
clean:
	rm -f $(SENDER) $(RECEIVER) $(SERVER)
