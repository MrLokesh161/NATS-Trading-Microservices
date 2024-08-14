package main

import (
    "encoding/json"
    "github.com/gorilla/websocket"
    "github.com/nats-io/nats.go"
    "log"
    "net/http"
)

var (
    nc         *nats.Conn
    upgrader   = websocket.Upgrader{}
    clients    = make(map[*websocket.Conn]bool)
    broadcast  = make(chan []byte)
)

func init() {
    var err error
    nc, err = nats.Connect(nats.DefaultURL)
    if err != nil {
        log.Fatalf("Failed to connect to NATS server: %v", err)
    }
    log.Println("Connected to NATS server.")

    go func() {
        _, err := nc.QueueSubscribe("order.completed", "trade-group", func(msg *nats.Msg) {
            log.Printf("Received order completion: %s", string(msg.Data))
            broadcast <- msg.Data
        })
        if err != nil {
            log.Fatalf("Failed to subscribe to NATS topic: %v", err)
        }
    }()
}

func orderHandler(w http.ResponseWriter, r *http.Request) {
    log.Println("Received order request")

    var order map[string]interface{}
    if err := json.NewDecoder(r.Body).Decode(&order); err != nil {
        log.Printf("Error decoding request body: %v", err)
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    data, err := json.Marshal(order)
    if err != nil {
        log.Printf("Error marshalling order data: %v", err)
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    if err := nc.Publish("order.created", data); err != nil {
        log.Printf("Error publishing message to NATS: %v", err)
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    log.Println("Order published to NATS.")
    w.WriteHeader(http.StatusAccepted)
}

func wsHandler(w http.ResponseWriter, r *http.Request) {
    conn, err := upgrader.Upgrade(w, r, nil)
    if err != nil {
        log.Printf("Error upgrading connection: %v", err)
        return
    }
    defer conn.Close()

    clients[conn] = true
    log.Println("WebSocket connection established")

    for {
        _, msg, err := conn.ReadMessage()
        if err != nil {
            log.Printf("Error reading message: %v", err)
            delete(clients, conn)
            break
        }
        log.Printf("Received message: %s", msg)
    }
}

func handleMessages() {
    for {
        msg := <-broadcast
        for client := range clients {
            err := client.WriteMessage(websocket.TextMessage, msg)
            if err != nil {
                log.Printf("Error writing message to WebSocket: %v", err)
                client.Close()
                delete(clients, client)
            }
        }
    }
}

func main() {
    go handleMessages()

    log.Println("Starting server on port 8080")
    http.HandleFunc("/order", orderHandler)
    http.HandleFunc("/ws", wsHandler)
    log.Fatal(http.ListenAndServe(":8080", nil))
}
