provider "hashitalks" {
  api_endpoint = "http://localhost:12345/"
}

terraform {
  required_providers {
    hashitalks = {
      version = "0.6.0"
      source  = "staging-registry.terraform.io/paddycarver/hashitalks"
    }
  }
}
