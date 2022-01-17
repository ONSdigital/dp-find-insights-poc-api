# CloudWatch Log Group to hold lambda messages
#
# Lambdas already log to /aws/lambda/<function-name>
# But we set up the resource here so we can specify a retention policy and tags.
#
resource "aws_cloudwatch_log_group" "fi-group" {
  name              = "/aws/lambda/${aws_lambda_function.fi-hello.function_name}"
  retention_in_days = 3
  tags = {
    Project = var.project_tag
  }
}