####
# AWS Metrics Stream & Firehose
####

resource "random_string" "aws_random" {
  length  = 10
  numeric = true
  special = false
  upper   = false
}


##########
# Bucket
##########

resource "aws_s3_bucket" "s3_bucket" {
  bucket        = local.name_prefix
  force_destroy = true
}

resource "aws_s3_bucket_public_access_block" "s3_bucket_access_block" {
  bucket                  = aws_s3_bucket.s3_bucket.id
  block_public_acls       = true
  block_public_policy     = true
  ignore_public_acls      = true
  restrict_public_buckets = true
}

resource "aws_s3_bucket_logging" "logging" {
  bucket        = aws_s3_bucket.s3_bucket.id
  target_bucket = aws_s3_bucket.s3_bucket.id
  target_prefix = "accesslog/"
}

resource "aws_s3_bucket_lifecycle_configuration" "lifecycle" {
  bucket = aws_s3_bucket.s3_bucket.id

  rule {
    id     = "expire"
    status = "Enabled"
    expiration {
      days = 1
    }
  }
}

##########
# Metrics Stream
##########

resource "aws_cloudwatch_metric_stream" "metric_stream" {
  name          = "${local.name_prefix}-stream"
  role_arn      = aws_iam_role.metrics_role.arn
  firehose_arn  = aws_kinesis_firehose_delivery_stream.metrics_delivery_stream.arn
  output_format = "opentelemetry0.7"

  dynamic "include_filter" {
    for_each = var.datadog_metric_stream_namespaces
    iterator = item

    content {
      namespace = item.value
    }
  }
}

##########
# Kinesis Stream
##########

resource "aws_kinesis_firehose_delivery_stream" "metrics_delivery_stream" {
  name        = "${local.name_prefix}-firehose-delivery"
  destination = "http_endpoint"

  server_side_encryption {
    enabled = false
  }

  http_endpoint_configuration {
    role_arn           = aws_iam_role.firehose_role.arn
    url                = "https://awsmetrics-intake.datadoghq.com/v1/input"
    access_key         = var.datadog_api_key
    name               = "${local.name_prefix}-endpoint"
    buffering_size     = 4
    buffering_interval = var.datadog_buffer_seconds
    retry_duration     = 60
    s3_backup_mode     = "FailedDataOnly"

    cloudwatch_logging_options {
      enabled = false
    }

    processing_configuration {
      enabled = false
    }

    request_configuration {
      content_encoding = "GZIP"
    }

    s3_configuration {
      role_arn            = aws_iam_role.firehose_role.arn
      bucket_arn          = aws_s3_bucket.s3_bucket.arn
      buffering_interval  = 300 # seconds
      buffering_size      = 5   # MB
      compression_format  = "UNCOMPRESSED"
      error_output_prefix = "metrics/"

      cloudwatch_logging_options {
        enabled = false
      }
    }
  }
}

##########
# IAM (Firehose)
##########

resource "aws_iam_role" "firehose_role" {
  name = "${local.name_prefix}-firehose"
  assume_role_policy = templatefile("${path.module}/../templates/firehose_assume_role.tmpl", {
    AWS_ACCOUNT_ID = data.aws_caller_identity.current.account_id,
  })
}

resource "aws_iam_policy" "firehose_s3_upload_policy" {
  policy = templatefile("${path.module}/../templates/firehose_s3_upload_policy.tmpl", {
    BUCKET_NAME = aws_s3_bucket.s3_bucket.id
  })
}

resource "aws_iam_policy" "firehose_delivery_policy" {
  policy = templatefile("${path.module}/../templates/firehose_delivery_policy.tmpl", {
    KINESIS_FIREHOSE_ARN = aws_kinesis_firehose_delivery_stream.metrics_delivery_stream.arn
  })
}

resource "aws_iam_role_policy_attachment" "firehose_policy_attach" {
  policy_arn = aws_iam_policy.firehose_delivery_policy.arn
  role       = aws_iam_role.firehose_role.name
}

resource "aws_iam_role_policy_attachment" "firehose_s3_policy_attach" {
  policy_arn = aws_iam_policy.firehose_s3_upload_policy.arn
  role       = aws_iam_role.firehose_role.name
}

##########
# IAM (Metric Stream)
##########

resource "aws_iam_role" "metrics_role" {
  name               = "${local.name_prefix}-metrics-stream"
  assume_role_policy = templatefile("${path.module}/../templates/metrics_assume_role.tmpl", {})
}

resource "aws_iam_policy" "metrics_policy" {
  policy = templatefile("${path.module}/../templates/metrics_policy.tmpl", {
    AWS_REGION  = var.aws_region
    AWS_ACCOUNT = data.aws_caller_identity.current.account_id
    ROLE_NAME   = aws_iam_role.metrics_role.name
  })
}

resource "aws_iam_role_policy_attachment" "metrics_policy_attach" {
  policy_arn = aws_iam_policy.metrics_policy.arn
  role       = aws_iam_role.metrics_role.name
}

##########
# Firehose Alarms
##########

# NB: Intentionally not in Datadog as we want to be alerted if the stream is not sending data to Datadog
resource "aws_cloudwatch_metric_alarm" "firehose_delivery_success_less_than" {
  count               = var.monitor_delivery_status ? 1 : 0
  alarm_name          = "${local.name_prefix}-delivery-success"
  comparison_operator = "LessThanThreshold"
  evaluation_periods  = 1
  metric_name         = "DeliveryToHttpEndpoint.Success"
  namespace           = "AWS/Firehose"
  period              = 300
  statistic           = "Minimum"
  dimensions = {
    DeliveryStreamName = aws_kinesis_firehose_delivery_stream.metrics_delivery_stream.name
  }
  threshold         = 1
  alarm_description = "Alarm if the Firehose delivery to the HTTP Datadog endpoint success rate is less than 100% for 5 minute"
  alarm_actions     = [aws_sns_topic.alarm_sns_topic[0].arn]
}

# SNS Topic for alerts
resource "aws_sns_topic" "alarm_sns_topic" {
  count = var.monitor_delivery_status ? 1 : 0
  name  = "${local.name_prefix}-topic"
}