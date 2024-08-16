package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/nats-io/nats.go"
)

var nc *nats.Conn

func init() {
	log.Println("Initializing connection to NATS server...")
	var err error
	nc, err = nats.Connect(nats.DefaultURL)
	if err != nil {
		log.Fatalf("Failed to connect to NATS server: %v", err)
	}
	log.Println("Connected to NATS server.")

	log.Println("Subscribing to 'order.created' subject...")
	_, err = nc.QueueSubscribe("order.created", "order-group", func(msg *nats.Msg) {
		log.Printf("Received order: %s", string(msg.Data))

		log.Println("Publishing to 'order.completed' subject...")
		err := nc.Publish("order.completed", msg.Data)
		if err != nil {
			log.Printf("Failed to publish message to 'order.completed': %v", err)
		} else {
			log.Println("Message successfully published to 'order.completed'.")
		}
	})
	if err != nil {
		log.Fatalf("Failed to subscribe to 'order.created': %v", err)
	}
	log.Println("Subscribed to 'order.created' subject.")
}

func main() {
	log.Println("Service is running. Waiting for messages...")

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	<-stop
	log.Println("Shutting down service...")
	nc.Close()
	log.Println("NATS connection closed. Service stopped.")
}
