terraform {
  required_providers {
    alicloudsecurity = {
      source = "registry.terraform.io/trendmicro/alicloudsecurity"
    }
  }
}

provider "alicloudsecurity" {
  # Configuration-based authentication
  visionone_api_key = "your_api_key_here"
  visionone_region = "us"
}

data "alicloudsecurity_connected_account" "connected_account" {}
