
resource "datadog_monitor" "monitor_datadog_demo_Count" {
  name                = "[Metric] - datadog-demo - Not enough count being monitored by datadog-demo"
  type                = "metric alert"
  message             = "{{#is_alert}}\nLooks at the counter1 *datadog-demo*.\n{{/is_alert}}\n{{#is_recovery}}\n\n{{/is_recovery}}\n@slack-demo-alerts-test"
  query               = "max(last_5m):max:datadog_demo.count{type:counter1} by {type} < 1"
  include_tags        = true
  no_data_timeframe   = 15
  notify_no_data      = true
  timeout_h           = 1
  priority            = "2"
  require_full_window = false
  monitor_thresholds {
    critical          = 1
    critical_recovery = 1.001
    ok                = ""
    unknown           = ""
    warning           = ""
    warning_recovery  = ""
  }
  monitor_threshold_windows {
    recovery_window = null
    trigger_window  = null
  }
  tags = ["app:sampleApp"]
}