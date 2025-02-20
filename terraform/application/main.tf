provider "aws" {
  region = var.region
  # access_key = var.my_aws_access_key
  # secret_key = var.my_aws_secret_key
}

provider "datadog" {
  api_key = var.my_datadog_api_key
  app_key = var.my_datadog_app_key
}

terraform {
  backend "s3" {
    bucket               = "test01-1"
    key                  = "terraform/datadog_demo_application/terraform.tfstate"
    region               = "us-east-1"
    workspace_key_prefix = "env"
  }
}

data "aws_caller_identity" "current" {}