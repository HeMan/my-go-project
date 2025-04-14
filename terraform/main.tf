terraform {
  required_version = ">= 1.0.0"

  required_providers {
    digitalocean = {
      source  = "digitalocean/digitalocean"
      version = ">= 2.0.0"
    }
  }
}

provider "digitalocean" {
  token = var.digitalocean_token
}

resource "digitalocean_app" "my_go_project_app" {
  spec {
    name   = "my-go-project-app"
    region = "ams3"

    service {
      name        = "my-go-project"
      http_port   = 8080
      run_command = "./my-go-project"

      image {
        registry      = "ghcr.io"
        repository    = "heman/my-go-project"
        tag           = "latest"
        registry_type = "GHCR"
      }
    }
  }
}

output "app_url" {
  value = digitalocean_app.my_go_project_app.live_url
}
