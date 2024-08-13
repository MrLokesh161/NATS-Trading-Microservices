job "trade-stream-service" {
  datacenters = ["dc1"]
  type = "service"

  task "trade-stream-service" {
    driver = "docker"

    config {
      image = "your-docker-image/trade-stream-service"
      port_map {
        http = 8081
      }
    }

    resources {
      cpu    = 500
      memory = 256
      network {
        mbits = 10
        port "http" {}
      }
    }

    service {
      name = "trade-stream-service"
      tags = ["websocket"]
      port = "http"
    }
  }
}
