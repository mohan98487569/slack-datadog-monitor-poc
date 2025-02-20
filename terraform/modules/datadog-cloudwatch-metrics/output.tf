output "my_aws_account_id" {
  value = data.aws_caller_identity.current.account_id
}

output "my_aws_caller_arn" {
  value = data.aws_caller_identity.current.arn
}

output "my_aws_caller_user" {
  value = data.aws_caller_identity.current.user_id
}