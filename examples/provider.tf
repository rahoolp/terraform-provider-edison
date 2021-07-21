provider "hashitalks" {
  api_endpoint = "http://localhost:12345/"
}

terraform {
  required_providers {
    hashitalks = {
      version = "0.1.0"
      source  = "registry.terraform.io/hashicorp/hashitalks"
    }
  }
}
