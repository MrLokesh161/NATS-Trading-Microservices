package main

import (
    "github.com/nats-io/nats.go"
    "github.com/gorilla/websocket"
    "log"
    "net/http"
)

var (
    nc        *nats.Conn
    upgrader  = websocket.Upgrader{}
    clients   = make(map[*websocket.Conn]bool)
    broadcast = make(chan []byte)
)

func init() {
    log.Println("Initializing connection to NATS server...")
    var err error
    nc, err = nats.Connect(nats.DefaultURL)
    if err != nil {
        log.Fatal("Failed to connect to NATS server:", err)
    }
    log.Println("Connected to NATS server.")

    log.Println("Subscribing to 'order.completed' subject...")
    _, err = nc.QueueSubscribe("order.completed", "trade-group", func(msg *nats.Msg) {
        log.Printf("Order completed: %s", string(msg.Data))
        log.Println("Broadcasting message to WebSocket clients...")
        broadcast <- msg.Data
    })
    if err != nil {
        log.Fatal("Failed to subscribe to 'order.completed':", err)
    }
    log.Println("Subscribed to 'order.completed' subject.")
}

func wsHandler(w http.ResponseWriter, r *http.Request) {
    log.Println("Upgrading HTTP connection to WebSocket...")
    conn, err := upgrader.Upgrade(w, r, nil)
    if err != nil {
        log.Println("Error while upgrading connection:", err)
        return
    }
    defer conn.Close()

    clients[conn] = true
    log.Println("WebSocket connection established")

    for {
        _, msg, err := conn.ReadMessage()
        if err != nil {
            log.Println("Error reading message:", err)
            delete(clients, conn)
            break
        }
        log.Printf("Received message from WebSocket client: %s", msg)
    }
}

func handleMessages() {
    for {
        msg := <-broadcast
        log.Printf("Broadcasting message to %d WebSocket clients", len(clients))
        for client := range clients {
            err := client.WriteMessage(websocket.TextMessage, msg)
            if err != nil {
                log.Printf("Error writing message to WebSocket client: %v", err)
                client.Close()
                delete(clients, client)
            }
        }
    }
}

func main() {
    log.Println("Starting WebSocket server on port 8081")
    go handleMessages()

    http.HandleFunc("/ws", wsHandler)
    log.Fatal(http.ListenAndServe(":8081", nil))
}