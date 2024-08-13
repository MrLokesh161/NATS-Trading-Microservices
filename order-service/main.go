package main

import (
    "github.com/nats-io/nats.go"
    "log"
)

var nc *nats.Conn

func init() {
    var err error
    nc, err = nats.Connect(nats.DefaultURL)
    if err != nil {
        log.Fatal(err)
    }

    _, err = nc.QueueSubscribe("order.created", "order-group", func(msg *nats.Msg) {
        log.Printf("Received order: %s", string(msg.Data))
        nc.Publish("order.completed", msg.Data)
    })
    if err != nil {
        log.Fatal(err)
    }
}

func main() {
    select {}
}
