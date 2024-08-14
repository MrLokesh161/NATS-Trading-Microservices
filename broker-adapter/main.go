package main

import (
    "github.com/nats-io/nats.go"
    "log"
)

var nc *nats.Conn

func init() {
    log.Println("Initializing connection to NATS server...")
    var err error
    nc, err = nats.Connect(nats.DefaultURL)
    if err != nil {
        log.Fatal("Failed to connect to NATS server:", err)
    }
    log.Println("Connected to NATS server.")

    log.Println("Subscribing to 'order.completed' subject...")
    _, err = nc.QueueSubscribe("order.completed", "broker-group", func(msg *nats.Msg) {
        log.Printf("Received order completion message: %s", string(msg.Data))
        log.Println("Processing order...")
        
        err := nc.Publish("order.completed", msg.Data)
        if err != nil {
            log.Printf("Failed to publish acknowledgment: %v", err)
        } else {
            log.Printf("Acknowledged order completion: %s", string(msg.Data))
        }
    })
    if err != nil {
        log.Fatal("Failed to subscribe to 'order.completed':", err)
    }
    log.Println("Subscribed to 'order.completed' subject.")
}

func main() {
    log.Println("Service is running. Waiting for messages...")
    select {}
}