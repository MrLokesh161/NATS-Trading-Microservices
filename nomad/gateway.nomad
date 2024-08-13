job "gateway" {
  datacenters = ["dc1"]
  type = "service"

  group "gateway" {
    network {
      port "http" {
        static = 8080
      }
    }

    task "gateway" {
      driver = "docker"

      config {
        image = "your-docker-image/gateway"
        port_map {
          http = 8080
        }
      }

      resources {
        cpu    = 500
        memory = 256
      }

      service {
        name = "gateway"
        tags = ["api"]
        port = "http"
      }
    }
  }
}
