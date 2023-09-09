# Variables
variable "api_endpoint" {
  description = "API end point of the Edison provider"
  type        = string
}

variable "token" {
  description = "Security token to be passed to the Edison provider"
  type        = string
}

provider "edison" {
  api_endpoint = var.api_endpoint
  token = var.token
}

terraform {
  required_providers {
    edison = {
      version = "0.6.0"
      source  = "hashicorp.com/edu/edison"
    }
  }
}
