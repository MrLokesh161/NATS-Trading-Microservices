package main

import (
    "encoding/json"
    "github.com/nats-io/nats.go"
    "log"
    "net/http"
)

var nc *nats.Conn

func init() {
    var err error
    nc, err = nats.Connect(nats.DefaultURL)
    if err != nil {
        log.Fatal(err)
    }
}

func orderHandler(w http.ResponseWriter, r *http.Request) {
    var order map[string]interface{}
    if err := json.NewDecoder(r.Body).Decode(&order); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    data, err := json.Marshal(order)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    if err := nc.Publish("order.created", data); err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusAccepted)
}

func main() {
    http.HandleFunc("/order", orderHandler)
    log.Fatal(http.ListenAndServe(":8080", nil))
}
