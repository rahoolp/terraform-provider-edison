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
  api_endpoint = "http://localhost:12345/"
  token = "secrettoken"
}

terraform {
  required_providers {
    edison = {
      version = "0.6.0"
      source  = "hashicorp.com/edu/edison"
    }
  }
}
