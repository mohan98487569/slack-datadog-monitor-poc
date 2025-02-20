terraform {
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = ">= 5.00.0"
    }
    datadog = {
      source  = "DataDog/datadog"
      version = ">= 3.00.0"
    }
  }
  required_version = ">= 0.14"
}