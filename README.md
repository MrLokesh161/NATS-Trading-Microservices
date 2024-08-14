Microservices Trading System

# Project Overview

This project consists of multiple microservices that work together to handle trading operations. The architecture includes:

- **Gateway Service**: Provides an API for order creation and serves as the entry point for client requests.
- **Order Service**: Publishes order creation events to NATS.
- **Trade Stream Service**: Subscribes to order completion events and provides real-time updates via WebSocket.
- **Broker Service**: Processes and logs order completion events.

## Prerequisites

Make sure you have the following installed:

- Go 1.20 or higher
- Docker (for running NATS server)

Also, ensure that a NATS server instance is running.

## Getting Started

1. Clone the Repository

```sh
git clone <repository-url>
cd <repository-directory>
```

2. Set Up the Go Modules

Navigate to each service directory and initialize the Go modules:

```sh
cd gateway
go mod tidy

cd ../order-service
go mod tidy

cd ../trade-stream-service
go mod tidy

cd ../broker-service
go mod tidy
```

3. Run NATS Server

You can run a NATS server using Docker. Use the following command to start it:

```sh
docker run -d -p 4222:4222 nats:latest
```

4. Build and Run Each Service

### Gateway Service

The Gateway Service listens for HTTP requests and interacts with NATS and WebSocket clients.

To run the Gateway Service:

```sh
cd gateway
go run main.go
```

Endpoints:

- `POST /order`: Accepts order requests and publishes them to NATS.

### Order Service

The Order Service publishes order creation events to NATS.

To run the Order Service:

```sh
cd order-service
go run main.go
```

NATS Subscription:

- Subscribes to `order.created` and publishes to `order.completed`.

### Trade Stream Service

The Trade Stream Service listens for order completion events and broadcasts them to WebSocket clients.

To run the Trade Stream Service:

```sh
cd trade-stream-service
go run main.go
```

WebSocket URL:

- WebSocket URL: `ws://localhost:8080/ws`

### Broker Service

The Broker Service logs order completion events received from NATS.

To run the Broker Service:

```sh
cd broker-service
go run main.go
```
