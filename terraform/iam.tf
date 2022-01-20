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
