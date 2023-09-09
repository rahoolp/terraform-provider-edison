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
