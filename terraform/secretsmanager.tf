# Secret holding pg password
#
resource "aws_secretsmanager_secret" "fi-pg" {
  name       = "fi-pg"
  kms_key_id = aws_kms_key.fi-cmk.id
  tags = {
    Project = var.project_tag
  }
}
