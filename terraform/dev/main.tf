terraform {
  required_version = ">= 1.0.0"

  required_providers {
    portainer = {
      source  = "grulicht/portainer"
      version = "1.3.0"
    }
  }

  backend "local" {
    path = "terraform-dev.tfstate"
  }
}

provider "portainer" {
  endpoint = var.portainer_url
  api_key  = var.portainer_api_key
  skip_ssl_verify = true
}

variable "postgres_password" {
  description = "The password for the PostgreSQL database"
  type        = string
  sensitive   = true
}

resource "portainer_stack" "my_go_project_stack" {
  name        = "my-go-project-app"
  endpoint_id = var.portainer_endpoint_id
  deployment_type = "standalone"
  method = "string"
  tlsskip_verify = true

  stack_file_content = <<-EOT
version: '3.8'
services:
  my-go-project-app:
    image: ghcr.io/heman/my-go-project:latest
    ports:
      - "8080:8080"
    environment:
      POSTGRES_USER: local_pgdb
      POSTGRES_PASSWORD: ${var.postgres_password}
      POSTGRES_DB: local_pgdb
      POSTGRES_HOSTNAME: 172.20.0.1
      POSTGRES_PORT: "5432"
      POSTGRES_SSLMODE: prefer
      POSTGRES_TIMEZONE: Europe/Stockholm
EOT
}
