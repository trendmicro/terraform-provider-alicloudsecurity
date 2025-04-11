terraform {
  required_providers {
    hashicups = {
      source = "hashicorp.com/edu/hashicups"
    }
  }
}

provider "hashicups" {
  host     = "http://localhost:19090"
  username = "education"
  password = "test123"
}

locals {
  alicloud_accesskey        = "LTAI5t6DNNzPrjEDzjWiBTXi"
  alicloud_accesskey_secret = "kPX1SiQnU7Fbh0MNyXHr4qijgy4x6A" 
  alicloud_region           = "us-west-1"
}

data "hashicups_visionone_alicloud_account" "managed_account" {

  name                = "cam-test-automation"
  alicloud_accesskey        = local.alicloud_accesskey
  alicloud_accesskey_secret = local.alicloud_accesskey_secret
  alicloud_region           = local.alicloud_region
}

resource "hashicups_visionone_alicloud_account_connection" "connection" {
  role_name = "demo-visionone-cloud-security-role"
  role_description = "a demo role for creating one by sdk xxxxxxx"
  alicloud_accesskey        = local.alicloud_accesskey
  alicloud_accesskey_secret = local.alicloud_accesskey_secret
  alicloud_region           = local.alicloud_region
}

# module "feature_a" {
#   source = "./feature_a"
#   alicloud_accesskey        = "LTAI5tF1eaHNv9X2Zh6d5LLi"
#   alicloud_accesskey_secret = "4Jupc7n9A82HAS6Ww5EvEACVVhRnw3"
#   alicloud_region           = "us-west-1"
# }

# module "feature_b" {
#   source = "./feature_b"
#   alicloud_accesskey        = "LTAI5tF1eaHNv9X2Zh6d5LLi"  
#   alicloud_accesskey_secret = "4Jupc7n9A82HAS6Ww5EvEACVVhRnw3"
#   alicloud_region           = "us-west-1"
# }

output "managed_alicloud_account" {
  value = {
    id = data.hashicups_visionone_alicloud_account.managed_account.id
    name = data.hashicups_visionone_alicloud_account.managed_account.name
    email = data.hashicups_visionone_alicloud_account.managed_account.email
  }
}

output "visionone_role" {
  value = {
    role_id = resource.hashicups_visionone_alicloud_account_connection.connection.role_id
    role_name = resource.hashicups_visionone_alicloud_account_connection.connection.role_name
    role_arn = resource.hashicups_visionone_alicloud_account_connection.connection.role_arn
    role_desc = resource.hashicups_visionone_alicloud_account_connection.connection.role_description
  }
}
