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

# /ckmeans
#
resource "aws_api_gateway_resource" "fi-ckmeans" {
  rest_api_id = aws_api_gateway_rest_api.fi-hello.id
  parent_id   = aws_api_gateway_rest_api.fi-hello.root_resource_id
  path_part   = "ckmeans"
}

# GET /ckmeans
#
resource "aws_api_gateway_method" "fi-get-ckmeans" {
  rest_api_id   = aws_api_gateway_rest_api.fi-hello.id
  resource_id   = aws_api_gateway_resource.fi-ckmeans.id
  http_method   = "GET"
  authorization = "NONE"
}

# Integrate GET /hello/{dataset} method with lambda
#
resource "aws_api_gateway_integration" "fi-get-ckmeans" {
  rest_api_id = aws_api_gateway_rest_api.fi-hello.id
  resource_id = aws_api_gateway_resource.fi-ckmeans.id
  http_method = aws_api_gateway_method.fi-get-ckmeans.http_method

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
      aws_api_gateway_resource.fi-ckmeans.id,
      aws_api_gateway_method.fi-get-hello.id,
      aws_api_gateway_method.fi-get-ckmeans.id,
      aws_api_gateway_integration.fi-get-hello.id,
      aws_api_gateway_integration.fi-get-ckmeans.id,
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