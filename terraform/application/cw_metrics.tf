// sends sampleApp custom metrics to datadog
module "test_cw_metrics" {
  source                           = "../modules/datadog-cloudwatch-metrics"
  datadog_api_key                  = var.my_datadog_api_key
  datadog_buffer_seconds           = 60
  datadog_metric_stream_namespaces = ["datadog-demo"]
  monitor_delivery_status          = true
  aws_region                       = var.region
}
