job "order-service" {
  datacenters = ["dc1"]
  type = "service"

  task "order-service" {
    driver = "docker"

    config {
      image = "your-docker-image/order-service"
    }

    resources {
      cpu    = 500
      memory = 256
    }

    service {
      name = "order-service"
    }
  }
}
