package main

import (
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/gorilla/websocket"
	"github.com/nats-io/nats.go"
)

var (
	nc        *nats.Conn
	upgrader  = websocket.Upgrader{}
	clients   = make(map[*websocket.Conn]bool)
	broadcast = make(chan []byte)
	mutex     sync.Mutex // Protects the clients map
)

func init() {
	log.Println("Initializing connection to NATS server...")
	var err error
	nc, err = nats.Connect(nats.DefaultURL)
	if err != nil {
		log.Fatalf("Failed to connect to NATS server: %v", err)
	}
	log.Println("Connected to NATS server.")

	log.Println("Subscribing to 'order.completed' subject...")
	_, err = nc.QueueSubscribe("order.completed", "trade-group", func(msg *nats.Msg) {
		log.Printf("Order completed: %s", string(msg.Data))
		log.Println("Broadcasting message to WebSocket clients...")
		broadcast <- msg.Data
	})
	if err != nil {
		log.Fatalf("Failed to subscribe to 'order.completed': %v", err)
	}
	log.Println("Subscribed to 'order.completed' subject.")
}

func wsHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Upgrading HTTP connection to WebSocket...")
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Failed to upgrade connection: %v", err)
		return
	}
	defer conn.Close()

	mutex.Lock()
	clients[conn] = true
	mutex.Unlock()
	log.Println("WebSocket connection established")

	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			log.Printf("Error reading message: %v", err)
			mutex.Lock()
			delete(clients, conn)
			mutex.Unlock()
			break
		}
		log.Printf("Received message from WebSocket client: %s", msg)
	}
}

func handleMessages() {
	for {
		msg := <-broadcast
		log.Printf("Broadcasting message to %d WebSocket clients", len(clients))

		mutex.Lock()
		for client := range clients {
			err := client.WriteMessage(websocket.TextMessage, msg)
			if err != nil {
				log.Printf("Error writing message to WebSocket client: %v", err)
				client.Close()
				delete(clients, client)
			}
		}
		mutex.Unlock()
	}
}

func main() {
	log.Println("Starting WebSocket server on port 8081")
	go handleMessages()

	http.HandleFunc("/ws", wsHandler)

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-stop
		log.Println("Shutting down server...")
		nc.Close()
		os.Exit(0)
	}()

	log.Fatal(http.ListenAndServe(":8081", nil))
}
