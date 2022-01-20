# Create a KMS customer master key to encrypt/decrypt sensitive data needed by the census api,
# such as the postgres password.
#
# The Statement preserves permissions for normal users to access the key.
# Additional statements should be added for other users of the key.
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
    }
  ]
}
EOF

  tags = {
    Project = var.project_tag
  }
}
