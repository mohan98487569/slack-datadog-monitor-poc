variable "region" {
  default = "us-east-1"
}

variable "my_datadog_api_key" {
  type      = string
  sensitive = true
  default   = "testapikey"
}

variable "my_datadog_app_key" {
  type      = string
  sensitive = true
  default   = "testappkey"
}

