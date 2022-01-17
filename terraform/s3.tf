# bucket to hold db dumps
#
resource "aws_s3_bucket" "db-dumps" {
  bucket = "find-insights-db-dumps"
  acl    = "private"

  tags = {
    Project = var.project_tag
  }
}