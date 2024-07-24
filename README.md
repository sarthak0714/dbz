# DBZ 

## Overview
DBZ is a simple, efficient, and robust database service built on Go and gRPC. It provides basic CRUD operations and transactional support to manage key-value pairs securely and reliably.

## Features
- **CRUD Operations**: Supports basic Create, Read, Update, and Delete operations.
- **Transactional Support**: Handles transactions with prepare and commit phases to ensure data integrity.
- **Concurrency Control**: Utilizes mutexes to manage data consistency across concurrent operations.
- **gRPC Interface**: Leverages gRPC for efficient and scalable communication.

## Components
- **Database Core**: The core logic for handling database operations, including transaction management.
- **gRPC Server**: Facilitates communication over gRPC, exposing methods for interacting with the database.
- **Client Example**: A sample client demonstrating how to interact with the database service using gRPC.

## Getting Started
To run the DBZ database service:
1. Clone the repository.
2. Build the server using Go:
   ```bash
   go build -o dbz_server cmd/main.go
   ```
3. Start the server:
   ```bash
   ./dbz_server
   ```
4. Run the example client to interact with the server:
   ```bash
   go run example/example.go
   ```

## Code Structure
- **Database Implementation**: `internal/database/database.go` 
- **Server Setup**: `cmd/main.go`
- **Protocol Buffers Definitions**: `pkg/api/database.proto`
- **Client Example**: `example/example.go` 

## Dependencies
- Go 1.22.2
- gRPC 1.65.0

For a full list of dependencies, refer to the `go.mod` file in the project root.

## TODO
* add auth
* add health check
* complete replication and sharding logic