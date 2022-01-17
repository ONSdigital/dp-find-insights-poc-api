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
  timeout          = 30

  environment {
    variables = {
      PGHOST          = "fi-database-2.cbhpmcuqy9vo.eu-central-1.rds.amazonaws.com"
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
