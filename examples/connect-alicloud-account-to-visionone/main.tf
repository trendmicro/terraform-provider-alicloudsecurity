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
  alicloud_account_id = "__module_cam_account_id__" # obtained from the output of CAM module
  alicloud_role_arn = "__module_cam_role_arn__" # obtained from the output of CAM module
  alicloud_oidc_provider_id = "__module_cam_oidc_provider_id__" # obtained from the output of CAM module 
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

data "alicloudsecurity_connected_account" "connected" {
  account_id = local.alicloud_account_id
}

# --- Variables --- #
variable "visionone_endpoint" {
  description = "value of the Vision One endpoint"
  type        = string
  default     = "__auto_fill_by_backend_" # auto-fill default value by the backend  
}

variable "visionone_endpoint_type" {
  description = "value of the Vision One endpoint type"
  type        = string
  default     = "__auto_fill_by_backend_" # auto-fill default value by the backend
}

variable "visionone_business_id" {
  description = "value of the Vision One business ID"
  type        = string
  default     = "__auto_fill_by_backend_" # auto-fill default value by the backend
}

variable "visionone_api_key" {
  description = "value of the Vision One API key"
  type        = string
  default     = "__auto_fill_by_backend_" # auto-fill default value by the backend
}

variable "visionone_region" {
  description = "value of the Vision One region"
  type        = string
  default     = "us"
}

variable "visionone_account_name" {
  description = "value of the Vision One account name"
  type        = string
  default     = "__auto_fill_by_backend_" # auto-fill default value by the backend
}

variable "visionone_account_description" {
  description = "value of the Vision One account description"
  type        = string
  default     = "__auto_fill_by_backend_" # auto-fill default value by the backend
}
