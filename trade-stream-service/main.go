package main

import (
    "github.com/nats-io/nats.go"
    "github.com/gorilla/websocket"
    "log"
    "net/http"
)

var (
    nc       *nats.Conn
    upgrader = websocket.Upgrader{}
)

func init() {
    var err error
    nc, err = nats.Connect(nats.DefaultURL)
    if err != nil {
        log.Fatal(err)
    }

    _, err = nc.QueueSubscribe("order.completed", "trade-group", func(msg *nats.Msg) {
        log.Printf("Order completed: %s", string(msg.Data))
    })
    if err != nil {
        log.Fatal(err)
    }
}

func wsHandler(w http.ResponseWriter, r *http.Request) {
    conn, err := upgrader.Upgrade(w, r, nil)
    if err != nil {
        log.Println("Error while upgrading connection:", err)
        return
    }
    defer conn.Close()

    for {
        _, msg, err := conn.ReadMessage()
        if err != nil {
            log.Println("Error reading message:", err)
            break
        }
        log.Printf("Received message: %s", msg)
    }
}

func main() {
    http.HandleFunc("/ws", wsHandler)
    log.Fatal(http.ListenAndServe(":8081", nil))
}
