job "broker-adapter" {
  datacenters = ["dc1"]
  type = "service"

  task "broker-adapter" {
    driver = "docker"

    config {
      image = "your-docker-image/broker-adapter"
    }

    resources {
      cpu    = 500
      memory = 256
    }

    service {
      name = "broker-adapter"
    }
  }
}
