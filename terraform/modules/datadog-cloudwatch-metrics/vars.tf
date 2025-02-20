variable "datadog_api_key" {
  description = "Datadog API Key"
  type        = string
  sensitive   = true
  default     = "fd3771af7da86ddd34b33a87d0565e0b"
}

variable "datadog_buffer_seconds" {
  description = "Datadog buffering interval in seconds"
  type        = number
  default     = 60

  validation {
    condition     = var.datadog_buffer_seconds >= 0 && var.datadog_buffer_seconds <= 900
    error_message = "Allowed buffering interval between 0-900 seconds."
  }
}

variable "datadog_metric_stream_namespaces" {
  description = "List of cloudwatch metric namespaces to push to datadog"
  type        = list(string)
  default     = []

  # List cannot contain AWS/EC2 or AWS/RDS since they cannot be filtered by tag, these will need to be done via polling instead.
  validation {
    condition     = !contains(var.datadog_metric_stream_namespaces, "AWS/EC2") && !contains(var.datadog_metric_stream_namespaces, "AWS/RDS")
    error_message = "AWS/EC2 and AWS/RDS namespaces are not allowed since they cannot be filtered in a reasonable way."
  }
}

variable "aws_region" {
  description = "AWS Region"
  type        = string
  default     = "us-east-1"
}

variable "monitor_delivery_status" {
  description = "Enable firehose delivery monitoring"
  type        = bool
  default     = false
}