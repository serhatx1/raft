<img width="1467" alt="raft1" src="https://github.com/user-attachments/assets/bfb7f22e-82e2-45d6-8c0d-4c4626f412ef" />

<img width="1436" alt="raft2" src="https://github.com/user-attachments/assets/b35e5d5d-3d05-4a7d-9760-d179f73947d0" />

This project is a simple implementation of the Raft consensus algorithm written in Go.

## Getting Started (with Docker)

You can easily start the project using Docker and docker-compose.

### 1. Build the Docker image and start the nodes

```sh
docker-compose up --build
```

This command will start three raft nodes:
- Node 1: http://localhost:8001
- Node 2: http://localhost:8002
- Node 3: http://localhost:8003

### 2. API Endpoints

Each node exposes the following HTTP endpoints:
- `/` : Simple welcome message
- `/request-vote` : Raft vote request (POST)
- `/heartbeat` : Raft heartbeat (POST)

### 3. Environment Variables

When started with docker-compose, each node is configured with the following environment variables:
- `NODE_ID` : Unique ID of the node
- `NODE_ADDR` : Address the node listens on (e.g. :8001)
- `NODE_PEERS` : Comma-separated list of other node addresses (e.g. http://node2:8002,http://node3:8003)

> Note: By default, the code expects to read these environment variables. If you run the code manually, you may need to provide them.

## Manual Start (for Developers)

If you have Go installed, you can run the project manually:

```sh
go run main.go
```

## License

MIT

## Directory Structure

- `cmd/node.go`: Node start, RPC server/client
- `internal/raft/state.go`: Node state and term management
- `internal/raft/rpc.go`: RequestVote, AppendEntries RPC
- `internal/raft/election.go`: Lider selection logic
- `internal/raft/log.go`: Log entry management
- `internal/statemachine.go`: Logic to handle committed commands
- `test/`: Unit and integration tests

## Setup

```sh
go mod tidy
go build ./cmd/node.go
```

## Explanation

Each file and directory serves a different responsibility in the Raft algorithm. Detailed explanations will be found in the relevant files.
