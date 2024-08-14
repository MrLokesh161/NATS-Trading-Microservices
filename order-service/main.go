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

    log.Println("Subscribing to 'order.created' subject...")
    _, err = nc.QueueSubscribe("order.created", "order-group", func(msg *nats.Msg) {
        log.Printf("Received order: %s", string(msg.Data))
        log.Println("Publishing to 'order.completed' subject...")
        err := nc.Publish("order.completed", msg.Data)
        if err != nil {
            log.Printf("Failed to publish message: %v", err)
        } else {
            log.Println("Message published to 'order.completed'.")
        }
    })
    if err != nil {
        log.Fatal("Failed to subscribe to 'order.created':", err)
    }
    log.Println("Subscribed to 'order.created' subject.")
}

func main() {
    log.Println("Service is running. Waiting for messages...")
    select {}
}