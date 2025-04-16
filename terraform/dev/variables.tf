variable "portainer_api_key" {
  description = "Portainer API key"
  type        = string
}

variable "portainer_url" {
  description = "Portainer API endpoint URL"
  type        = string
}

variable "portainer_endpoint_id" {
  description = "Portainer endpoint ID"
  type        = number
}

variable "enable_portainer_stack" {
  description = "Set to true to enable the creation of the Portainer stack"
  type        = bool
  default     = false
}
