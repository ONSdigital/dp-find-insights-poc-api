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