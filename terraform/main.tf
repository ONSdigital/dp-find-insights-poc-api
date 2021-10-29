# Very basic infra for hello world lambda.
#
# This is just a demo, so there are no links to the rest of provisioning.

terraform {
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 3.0"
    }
  }

  # The backend bucket for holding state is created out of band.
  # Good idea to enable versioning and lifecycle management.
  backend "s3" {
    bucket  = "find-insights-alpha-tfstate"
    key     = "tfstate"
    region  = "eu-central-1"
    profile = "development"
  }
}

# Using Frankfurt to avoid confusing people and to provide a bit of isolation.
#
variable "project_region" {
  default = "eu-central-1"
}

# When possible, create a "Project" tag on all resources associated with this project.
# (Unfortunately I don't have permissions to tag certain resources, so I can't use
# default tags; have to tag resources individually.)
#
variable "project_tag" {
  default = "find-insights-alpha"
}

# Assuming we will be using credentials file and profile for authentication during dev.
#
variable "aws_cred_file" {
  default = "~/.aws/credentials"
}
variable "aws_profile" {
  default = "development"
}

# There are other ways to specify region and authentication, but for now do it this way.
#
provider "aws" {
  region                  = var.project_region
  shared_credentials_file = pathexpand(var.aws_cred_file)
  profile                 = var.aws_profile
}

# I wanted to use something like the aws cli to upload new versions of the lambda into S3.
# But tha doesn't play well with terraform.
# So for now using terraform to push new artifacts.
# Maybe there is a way to do this with lambda aliases?

# S3 bucket to hold lambda deployment packages (zips)
#
#resource "aws_s3_bucket" "fi-lambdas" {
#  bucket = "find-insights-lambdas"
#  acl    = "private"
#  versioning {
#    enabled = true
#  }
#  lifecycle_rule {
#    enabled = true
#    noncurrent_version_expiration {
#      days = 7
#    }
#  }
#  tags = {
#    Project = var.project_tag
#  }
#}

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

# Policy allowing CloudWatch logging
# 
resource "aws_iam_policy" "fi-logging" {
  name        = "find-insights-logging"
  description = "policy to allow find insights lambda to log to cloudwatch"
  path        = "/"
  policy      = <<EOF
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
  # AccessDenied: User: arn:aws:iam::352437599875:user/DanielLawrence is not authorized to perform: iam:TagPolicy on resource: policy find-insights-logging
  #tags = {
  #  Project = var.project_tag
  #}
}

# Role which AWS Lambda can assume when running our lambda
#
resource "aws_iam_role" "fi-lambda-execution" {
  name        = "find-insights-execution"
  description = "role assumed by lambda at runtime"

  assume_role_policy = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
      {
        "Effect": "Allow",
        "Principal": {
            "Service": "lambda.amazonaws.com"
        },
        "Action": "sts:AssumeRole"
      }
  ]
}
EOF

  # AccessDenied: User: arn:aws:iam::352437599875:user/DanielLawrence is not authorized to perform: iam:TagRole on resource: arn:aws:iam::352437599875:role/find-insights-execution
  #  tags = {
  #    project = var.project_tag
  #  }
}

# Attach basic policy to runtime role
#
resource "aws_iam_role_policy_attachment" "fi-lambda-basic" {
  role       = aws_iam_role.fi-lambda-execution.name
  policy_arn = aws_iam_policy.fi-logging.arn
}

# Attach logging policy to runtime role
#
resource "aws_iam_role_policy_attachment" "fi-lambda-logs" {
  role       = aws_iam_role.fi-lambda-execution.name
  policy_arn = aws_iam_policy.fi-logging.arn
}

# The lambda itself
#
resource "aws_lambda_function" "fi-hello" {
  function_name    = "find-insights-hello"
  description      = "hello world lambda"
  role             = aws_iam_role.fi-lambda-execution.arn
  handler          = "hello"
  runtime          = "go1.x"
  filename         = "../hello.zip"
  source_code_hash = filebase64sha256("../hello.zip")

  # Permissions to log to cloudwatch have to be set up before the lambda runs,
  # but this resource doesn't directly reference the policy attachment, so
  # terraform's dag doesn't know that.
  depends_on = [
    aws_iam_role_policy_attachment.fi-lambda-logs
  ]

  tags = {
    project = var.project_tag
  }
}

# External API for our lambda
#
resource "aws_api_gateway_rest_api" "fi-hello" {
  name        = "find-insights-api"
  description = "api for find insights alpha lambda"
  tags = {
    Project = var.project_tag
  }
}

# /hello
#
resource "aws_api_gateway_resource" "fi-hello" {
  rest_api_id = aws_api_gateway_rest_api.fi-hello.id
  parent_id   = aws_api_gateway_rest_api.fi-hello.root_resource_id
  path_part   = "hello"
}

# GET /hello
# (I can't get POST to work yet.)
#
resource "aws_api_gateway_method" "fi-get-hello" {
  rest_api_id   = aws_api_gateway_rest_api.fi-hello.id
  resource_id   = aws_api_gateway_resource.fi-hello.id
  http_method   = "GET"
  authorization = "NONE"
}

# Integrate GET /hello method with lambda
#
resource "aws_api_gateway_integration" "fi-get-hello" {
  rest_api_id = aws_api_gateway_rest_api.fi-hello.id
  resource_id = aws_api_gateway_resource.fi-hello.id
  http_method = aws_api_gateway_method.fi-get-hello.http_method

  # lambda methods can only be invoked with POST integration_http_method
  integration_http_method = "POST"

  # AWS_PROXY type required for lambda integration
  type = "AWS_PROXY"

  uri = aws_lambda_function.fi-hello.invoke_arn
}

resource "aws_api_gateway_deployment" "fi-hello" {
  rest_api_id = aws_api_gateway_rest_api.fi-hello.id
}

resource "aws_api_gateway_stage" "dev" {
  deployment_id = aws_api_gateway_deployment.fi-hello.id
  rest_api_id   = aws_api_gateway_rest_api.fi-hello.id
  stage_name    = "dev"
}

# Permissions to allow apigateway to invoke lambda
#
# More: http://docs.aws.amazon.com/apigateway/latest/developerguide/api-gateway-control-access-using-iam-policies-to-invoke-api.html
#
resource "aws_lambda_permission" "fi-allow-apigw" {
  statement_id  = "AllowExecutionByAPIGateway"
  action        = "lambda:InvokeFunction"
  function_name = aws_lambda_function.fi-hello.function_name
  principal     = "apigateway.amazonaws.com"
  source_arn    = "${aws_api_gateway_rest_api.fi-hello.execution_arn}/*"
}
