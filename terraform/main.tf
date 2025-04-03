provider "digitalocean" {
  token = var.digitalocean_token
}

resource "digitalocean_app" "my_go_project_app" {
  spec {
    name = "my-go-project-app"
    region = "ams3"

    static_site {
      name       = "my-go-project"
      git {
        repo_clone_url = var.repo_url
        branch         = "main"
      }
      build_command = "go build -o my-go-project"
      run_command   = "./my-go-project"
      http_port     = 8080
    }

    database {
      name       = "my-go-project-dev-db"
      engine     = "pg"
      version    = "17"
      size       = "db-s-1vcpu-1gb"
      production = false
    }

    env {
      key   = "POSTGRES_HOSTNAME"
      value = "${my-go-project-dev-db.HOSTNAME}"
      scope = "RUN_AND_BUILD_TIME"
    }

    env {
      key   = "POSTGRES_USER"
      value = "${my-go-project-dev-db.USERNAME}"
      scope = "RUN_AND_BUILD_TIME"
    }

    env {
      key   = "POSTGRES_PASSWORD"
      value = "${my-go-project-dev-db.PASSWORD}"
      scope = "RUN_AND_BUILD_TIME"
    }

    env {
      key   = "POSTGRES_DB"
      value = "${my-go-project-dev-db.DATABASE}"
      scope = "RUN_AND_BUILD_TIME"
    }

    env {
      key   = "POSTGRES_PORT"
      value = "${my-go-project-dev-db.PORT}"
      scope = "RUN_AND_BUILD_TIME"
    }
  }
}

output "app_url" {
  value = digitalocean_app.my_go_project_app.live_url
}
