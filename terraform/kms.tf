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