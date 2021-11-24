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
# But that doesn't play well with terraform.
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
  policy_arn = "arn:aws:iam::aws:policy/service-role/AWSLambdaBasicExecutionRole"
}

# Policy to allow reading pg password
#
resource "aws_iam_policy" "fi-read-pg-password" {
  name        = "fi-read-pg-password"
  description = "allows reading postgress password"
  policy      = <<EOF
{
    "Version": "2012-10-17",
    "Statement": [
        {
            "Effect": "Allow",
            "Action": [
                "secretsmanager:GetResourcePolicy",
                "secretsmanager:GetSecretValue",
                "secretsmanager:DescribeSecret",
                "secretsmanager:ListSecretVersionIds"
            ],
            "Resource": [
              "${aws_secretsmanager_secret.fi-pg.arn}"
            ]
        },
        {
            "Effect": "Allow",
            "Action": "secretsmanager:ListSecrets",
            "Resource": "*"
        }
    ]
}
EOF

  # AccessDenied: User: arn:aws:iam::352437599875:user/DanielLawrence is not authorized to perform: iam:TagPolicy on resource: policy fi-read-pg-password
  #tags = {
  #  Project = var.project_tag
  #}
}

# Attach pg password policy to runtime role
#
resource "aws_iam_role_policy_attachment" "fi-lambda-pg-password" {
  role       = aws_iam_role.fi-lambda-execution.name
  policy_arn = aws_iam_policy.fi-read-pg-password.arn
}

# The lambda itself
#
resource "aws_lambda_function" "fi-hello" {
  function_name    = "find-insights-hello"
  description      = "hello world lambda"
  role             = aws_iam_role.fi-lambda-execution.arn
  handler          = "hello"
  runtime          = "go1.x"
  filename         = "../build/hello.zip"
  source_code_hash = filebase64sha256("../build/hello.zip")
  memory_size      = 256

  environment {
    variables = {
      PGHOST          = "fi-database-1.cbhpmcuqy9vo.eu-central-1.rds.amazonaws.com"
      PGPORT          = "54322"
      PGUSER          = "insights"
      PGDATABASE      = "census"
      FI_PG_SECRET_ID = aws_secretsmanager_secret.fi-pg.id
    }
  }

  tags = {
    project = var.project_tag
  }
}

# External API for our lambda
#
resource "aws_api_gateway_rest_api" "fi-hello" {
  name                     = "find-insights-api"
  description              = "api for find insights alpha lambda"
  minimum_compression_size = 2097152
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

# /hello/{dataset}
#
resource "aws_api_gateway_resource" "fi-hello-dataset" {
  rest_api_id = aws_api_gateway_rest_api.fi-hello.id
  parent_id   = aws_api_gateway_resource.fi-hello.id
  path_part   = "{dataset+}"
}

# GET /hello/{dataset}
# (I can't get POST to work yet.)
#
resource "aws_api_gateway_method" "fi-get-hello" {
  rest_api_id   = aws_api_gateway_rest_api.fi-hello.id
  resource_id   = aws_api_gateway_resource.fi-hello-dataset.id
  http_method   = "GET"
  authorization = "NONE"
}

# configure CORS on /hello/{dataset}
module "cors" {
  source = "squidfunk/api-gateway-enable-cors/aws"
  version = "0.3.3"

  api_id          = aws_api_gateway_rest_api.fi-hello.id
  api_resource_id = aws_api_gateway_resource.fi-hello-dataset.id
  allow_headers = [
  "Authorization",
  "Content-Type",
  "Content-Encoding",
  "X-Amz-Date",
  "X-Amz-Security-Token",
  "X-Api-Key"
]
}


# Integrate GET /hello/{dataset} method with lambda
#
resource "aws_api_gateway_integration" "fi-get-hello" {
  rest_api_id = aws_api_gateway_rest_api.fi-hello.id
  resource_id = aws_api_gateway_resource.fi-hello-dataset.id
  http_method = aws_api_gateway_method.fi-get-hello.http_method

  # lambda methods can only be invoked with POST integration_http_method
  integration_http_method = "POST"

  # AWS_PROXY type required for lambda integration
  type = "AWS_PROXY"

  uri = aws_lambda_function.fi-hello.invoke_arn
}

resource "aws_api_gateway_deployment" "fi-hello" {
  rest_api_id = aws_api_gateway_rest_api.fi-hello.id

  # must be redeployed after any change, so set triggers to detect changes in
  # all dependencies
  # https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/api_gateway_deployment
  #
  triggers = {
    redeployment = sha1(jsonencode([
      aws_api_gateway_resource.fi-hello.id,
      aws_api_gateway_resource.fi-hello-dataset.id,
      aws_api_gateway_method.fi-get-hello.id,
      aws_api_gateway_integration.fi-get-hello.id,
    ]))
  }

  lifecycle {
    create_before_destroy = true
  }
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

# Create a KMS customer master key to encrypt/decrypt sensitive data needed by the lambda,
# such as the postgres password.
#
# The first statement preserves permissions for normal users to access the key.
# The second statement is specific to the lambda.
#
resource "aws_kms_key" "fi-cmk" {
  description = "master key for find insights lambda passwords"
  policy      = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Sid": "Enable IAM User Permissions",
      "Effect": "Allow",
      "Principal": {
        "AWS": "arn:aws:iam::352437599875:root"
      },
      "Action": "kms:*",
      "Resource": "*"
    },
    {
      "Sid": "Allow lambda to use the key",
      "Effect": "Allow",
      "Principal": {
        "AWS": [
          "${aws_iam_role.fi-lambda-execution.arn}"
        ]
      },
      "Action": [
        "kms:Encrypt",
        "kms:Decrypt",
        "kms:ReEncrypt*",
        "kms:GenerateDataKey*",
        "kms:DescribeKey"
      ],
      "Resource": "*"
    }
  ]
}
EOF

  tags = {
    Project = var.project_tag
  }
}

# Secret holding pg password
#
resource "aws_secretsmanager_secret" "fi-pg" {
  name       = "fi-pg"
  kms_key_id = aws_kms_key.fi-cmk.id
  tags = {
    Project = var.project_tag
  }
}
