terraform {
  required_providers {
    alicloudsecurity = {
      source = "registry.terraform.io/trendmicro/alicloudsecurity"
      # version = "0.0.1"
    }
  }
}

provider "alicloudsecurity" {
  visionone_endpoint = var.visionone_endpoint
  visionone_endpoint_type = var.visionone_endpoint_type
  visionone_business_id = var.visionone_business_id
  visionone_api_key = var.visionone_api_key
  visionone_region = var.visionone_region
}

locals {
  alicloud_account_id = "5245320496665124" # obtained from the output of CAM module
  alicloud_role_arn = "acs:ram::5245320496665124:role/visiononerole" # obtained from the output of CAM module
  alicloud_oidc_provider_id = "vision_one_oidc" # obtained from the output of CAM module 
  alicloud_name = var.visionone_account_name
  alicloud_description = var.visionone_account_description
}

resource "alicloudsecurity_connected_account" "connection" {
  stack_state_region = "us-east-1" # the region of Terraform backend where the state files are stored 
  account_id = local.alicloud_account_id
  role_arn = local.alicloud_role_arn
  oidc_provider_id = local.alicloud_oidc_provider_id
  name = local.alicloud_name
  description = local.alicloud_description
}

# --- Variables --- #
variable "visionone_endpoint" {
  description = "value of the Vision One endpoint"
  type        = string
  default     = "https://cloudaccounts-dev-us.visionone.trendmicro.com" # auto-fill default value by the backend  
}

variable "visionone_endpoint_type" {
  description = "value of the Vision One endpoint type"
  type        = string
  default     = "express" # auto-fill default value by the backend
}

variable "visionone_business_id" {
  description = "value of the Vision One business ID"
  type        = string
  default     = "84429025-7a92-41cf-99d1-96234765d052" # auto-fill default value by the backend
}

variable "visionone_api_key" {
  description = "value of the Vision One API key"
  type        = string
  default     = "eyJ0eXAiOiJKV1QiLCJhbGciOiJSUzI1NiJ9.eyJjaWQiOiI4NDQyOTAyNS03YTkyLTQxY2YtOTlkMS05NjIzNDc2NWQwNTIiLCJjcGlkIjoic3ZwIiwicHBpZCI6ImN1cyIsIml0IjoxNzQ1MjU0MTc5LCJldCI6MTc1MzAzMDE3OCwiaWQiOiI1YzkwMWZjOC04YzNjLTQwNWEtOTQ3Yi0wZWEyZDMxNmYyYjAiLCJ0b2tlblVzZSI6ImN1c3RvbWVyIn0.j9sqMQa9E2Ghw17__5AJQxda79sTYVNKvHX2f9aTQpf3BSGPFkr3oScj2zpejmITBIx8e6OKjGu1bfVF03R_ylGimEoD_sjGjKOxUrNx2bN4KKYIL6L5iEgDeis4bn1pdlcXnMeTGiKn1vzKIKiWJw449ZpOW_61BUOp358FTkNtn1MRlWooCaTWuimMqR7y4cD_KTszHa3X6vxzBm86B_Cyzb08hjTl14Fsug0LwMpvT9kJosTUTz9f9PUCfr1OutXdUihZX1628W1g3wGHkGzoK-OFRy3Beg0M6ucas59CkLfnDTUaHZLlbZdgmuR1BcXa_5AFJIhrNQJCnMsa07_SiBpfQa7fH1HknetSCx3qbrgqGXkGf_2ehHb02T49tQ_KTtRaMCRJzBQvlYW-fspzT6qH_ZKvme88wm_UbPk497fvWQgvDmz7CTaqZGyhHg178gn9jI-6UNFN17WyYMq5jvsM9vmcL0q6UPv8NLDF8_3adZmF3Sz1pRxh5TkK52dfQjq-GEhkwr9ZWpjKseg71PG0pRmT0223pVRWV70h-sbOM5e82dcvgO-fIMXV-Y47FeaBTWwP-ViqajfME4wYL4EN9eATUCakibf2WQYj2Xuc4Sl0woncmG7WNZOt5njLBRL5wgYjUJozwCkFeN79kjvhFQlanJ-wmrsHR24" # auto-fill default value by the backend
}

variable "visionone_region" {
  description = "value of the Vision One region"
  type        = string
  default     = "us"
}

variable "visionone_account_name" {
  description = "value of the Vision One account name"
  type        = string
  default     = "visionone-account-name" # auto-fill default value by the backend
}

variable "visionone_account_description" {
  description = "value of the Vision One account description"
  type        = string
  default     = "visionone-account-description-250428" # auto-fill default value by the backend
}
