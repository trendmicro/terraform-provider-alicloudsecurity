terraform {
  required_providers {
    alicloud = {
      source  = "hashicorp/alicloud"
      version = "1.227.0"
    }

    alicloudsecurity = {
      source = "registry.terraform.io/trendmicro/alicloudsecurity"
      version = "0.0.1"
    }
  }
}

provider "alicloud" {
  region = var.alicloud_region
}

provider "alicloudsecurity" {
  visionone_endpoint = "autofill_by_trend"
  visionone_endpoint_type = "autofill_by_trend"
  visionone_business_id = "autofill_by_trend"
  visionone_api_key = "autofill_by_trend"
  visionone_region = var.vision_one_region
}

module "cam" {
  source = "remote_source_in_public_cloud_storage"
  cam_oidc_provider_url = var.vision_one_oidc_provider_url
  cam_subject_urn       = var.vision_one_subject_urn
  cam_oss_bucket_name   = var.vision_one_oss_bucket_name
  cam_template_version  = var.vision_one_barebone_version
}

module "feature_a" {
    source = "remote_source_in_public_cloud_storage"
    primary_region = module.cam.cam_deployed_region
    
    # ... other properties
}

resource "alicloudsecurity_connected_account" "connection" {
  depends_on = [ module.cam, module.feature_a ]

  stack_state_region = var.visionone_region
  account_id = module.cam.account_id
  role_arn = module.cam.vision_one_role_arn
  oidc_provider_id = module.cam.vision_one_oidc_provider_url
  name = var.account_alias
  description = var.account_description
}
