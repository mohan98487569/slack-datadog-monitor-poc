locals {
  name_prefix = "dd-cw-metrics-${var.aws_region}-${random_string.aws_random.result}"
}

data "aws_caller_identity" "current" {}