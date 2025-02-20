# lambdaFunction


This is an AWS Lambda function in Golang that listens for Slack event subscriptions(_Slack Bot_) and manages Datadog monitor muting/unmuting based on Slack replies accordingly.
When a **Datadog alert** message in a Slack channel is replied to as:

* `acknowledged`: The Datadog monitor is muted for 6 hours.
* `resolved`: The Datadog monitor is unmuted.


The function is deployed using Terraform.

## Build and Deploy Lambda Function

```
cd lambdaFunction
go mod tidy
GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o bootstrap main.go
zip function.zip main

cd ../terraform/application/
terraform init
terraform plan -input=false -out=terraform.tfplan
terraform apply -auto-approve terraform.tfplan
```
