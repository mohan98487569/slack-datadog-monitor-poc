//  configure a Slack channel for receiving Datadog alerts and notifications.
resource "datadog_integration_slack_channel" "demo_alerts_test_channel" {
  account_name = "demo_DMN"
  channel_name = "#demo-alerts-test"

  display {
    message  = true
    notified = true
    snapshot = true
    tags     = true
  }
}