# Microservices Trading System

## Project Overview

This project is a microservices-based trading system composed of several services that communicate through NATS. The architecture includes:

- **Gateway Service**: An API gateway that accepts client order requests and interacts with NATS and WebSocket clients.
- **Order Service**: Publishes order creation events to NATS.
- **Trade Stream Service**: Listens for order completion events and provides real-time updates via WebSocket.
- **Broker Service**: Processes, logs, and acknowledges order completion events.

## Prerequisites

Before getting started, ensure you have the following installed:

- **Go 1.20 or higher**
- **Docker** (to run the NATS server)

Make sure a NATS server instance is running before you start the services.

## Getting Started

### 1. Clone the Repository

Clone the repository to your local machine:

```sh
git clone <repository-url>
cd <repository-directory>
```

### 2. Set Up Go Modules

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

### 3. Run the NATS Server

Run a NATS server using Docker. Use the following command:

```sh
docker run -d -p 4222:4222 -p 6222:6222 -p 8222:8222 --name nats-server nats:latest
```

This command starts the NATS server and maps the standard ports.

### 4. Build and Run Each Service

#### Gateway Service

The Gateway Service serves as the entry point, accepting HTTP requests and interacting with NATS and WebSocket clients.

To run the Gateway Service:

```sh
cd gateway
go run main.go
```

Endpoints:

- `POST /order`: Accepts order requests, which are then published to NATS.

#### Order Service

The Order Service listens for order creation events, processes them, and publishes order completion events to NATS.

To run the Order Service:

```sh
cd order-service
go run main.go
```

NATS Subscription:

- **Subscribes to:** `order.created`
- **Publishes to:** `order.completed`

#### Trade Stream Service

The Trade Stream Service listens for order completion events and broadcasts them to connected WebSocket clients.

To run the Trade Stream Service:

```sh
cd trade-stream-service
go run main.go
```

WebSocket URL:

- **WebSocket URL:** `ws://localhost:8081/ws`

#### Broker Service

The Broker Service processes and logs order completion events received from NATS.

To run the Broker Service:

```sh
cd broker-service
go run main.go
```

NATS Subscription:

- **Subscribes to:** `order.completed`

## Docker Compose Setup

You can also use Docker Compose to manage and run all the services together with the NATS server. To do so:

1. Ensure Docker is running.
2. Navigate to the root of the project directory where the `docker-compose.yml` file is located.
3. Run the following command:

```sh
docker-compose up --build
```

This will build and start all the services along with the NATS server.

### Accessing the Services

- **Gateway Service API**: `http://localhost:8080/order`
- **Trade Stream WebSocket**: `ws://localhost:8081/ws`
- **NATS Dashboard** (if enabled): `http://localhost:8222`

## Conclusion

This project demonstrates a microservices architecture using Go, NATS, and WebSockets. Each service operates independently and communicates via NATS, enabling a scalable and resilient system for handling trading operations.

