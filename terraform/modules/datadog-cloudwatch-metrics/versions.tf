terraform {
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = ">= 5.00.0"
    }
    random = {
      source  = "hashicorp/random"
      version = ">= 3.0.0"
    }
  }
  required_version = ">= 0.14"
}