# Create Lambda Function
resource "aws_lambda_function" "datadog_monitor_request_sns_lambda" {
  function_name    = "datadog-monitor-request"
  role             = aws_iam_role.iam_for_datadog_monitor_request_lambda.arn
  handler          = "bootstrap"
  runtime          = "provided.al2"
  filename         = "../../lambdaFunction/function.zip"
  source_code_hash = filebase64sha256("../../lambdaFunction/function.zip")
  publish          = true
  environment {
    variables = {
      DATADOG_API_KEY = var.my_datadog_api_key,
      DATADOG_APP_KEY = var.my_datadog_app_key
    }
  }
}


# IAM Role for Lambda
resource "aws_iam_role" "iam_for_datadog_monitor_request_lambda" {
  name = "golang_lambda_role"

  assume_role_policy = jsonencode({
    Version = "2012-10-17"
    Statement = [{
      Action = "sts:AssumeRole"
      Effect = "Allow"
      Principal = {
        Service = "lambda.amazonaws.com"
      }
    }]
  })
}

resource "aws_iam_policy" "datadog_monitor_request_iam_policy" {
  name = "datadog-monitor-request-iam-policy"

  policy = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Action": [
        "logs:CreateLogGroup",
        "logs:CreateLogStream",
        "logs:PutLogEvents"
      ],
      "Resource": "arn:aws:logs:*:*:*",
      "Effect": "Allow"
    }
  ]
}
EOF
}

resource "aws_iam_role_policy_attachment" "lambda_logs" {
  role       = aws_iam_role.iam_for_datadog_monitor_request_lambda.name
  policy_arn = aws_iam_policy.datadog_monitor_request_iam_policy.arn
}


#
#
#


# Create API Gateway
resource "aws_apigatewayv2_api" "slack_api" {
  name          = "SlackToAPI"
  protocol_type = "HTTP"
}

# Create API Gateway Integration for Lambda
resource "aws_apigatewayv2_integration" "slack_lambda_integration" {
  api_id           = aws_apigatewayv2_api.slack_api.id
  integration_type = "AWS_PROXY"
  integration_uri  = aws_lambda_function.datadog_monitor_request_sns_lambda.invoke_arn
}

# Create API Gateway Route (POST /slack-events)
resource "aws_apigatewayv2_route" "slack_route" {
  api_id    = aws_apigatewayv2_api.slack_api.id
  route_key = "POST /slack-events"
  target    = "integrations/${aws_apigatewayv2_integration.slack_lambda_integration.id}"
}

# Create API Gateway Deployment
resource "aws_apigatewayv2_stage" "slack_stage" {
  api_id      = aws_apigatewayv2_api.slack_api.id
  name        = "api_gateway_deployment"
  auto_deploy = true
}

# Lambda Permission to Allow API Gateway Invocation
resource "aws_lambda_permission" "apigw_lambda" {
  statement_id  = "AllowAPIGatewayInvoke"
  action        = "lambda:InvokeFunction"
  function_name = aws_lambda_function.datadog_monitor_request_sns_lambda.function_name
  principal     = "apigateway.amazonaws.com"

  source_arn = "${aws_apigatewayv2_api.slack_api.execution_arn}/*/*"
}

# Output the API Gateway Invoke URL
output "api_gateway_url" {
  value       = "${aws_apigatewayv2_stage.slack_stage.invoke_url}/slack-events"
  description = "The invoke URL for Slack to send events"
}