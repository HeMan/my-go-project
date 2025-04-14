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

resource "digitalocean_database_cluster" "my_go_project_db" {
  name       = "my-go-project-db"
  engine     = "PG"
  version    = "17"
  size       = "db-s-1vcpu-1gb"
  region     = "ams3"
  node_count = 1

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

      env {
        key   = "POSTGRES_HOSTNAME"
        value = "$${my-go-project-db.HOSTNAME}"
      }

      env {
        key   = "POSTGRES_PORT"
        value = "$${my-go-project-db.PORT}"
      }

      env {
        key   = "POSTGRES_USER"
        value = "$${my-go-project-db.USERNAME}"
      }

      env {
        key   = "POSTGRES_PASSWORD"
        value = "$${my-go-project-db.PASSWORD}"
      }

      env {
        key   = "POSTGRES_DB"
        value = "$${my-go-project-db.DATABASE}"
      }
      
      env {
        key = "POSTGRES_SSLMODE"
        value = "prefer"
      }

      env {
        key = "POSTGRES_TIMEZONE"
        value = "Europe/Stockholm"
      }
    }

    database {
      name    = "my-go-project-db"
    }
  }
}

output "app_url" {
  value = digitalocean_app.my_go_project_app.live_url
}