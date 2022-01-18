# bucket to hold db dumps
#
resource "aws_s3_bucket" "db-dumps" {
  bucket = "find-insights-db-dumps"
  acl    = "private"

  tags = {
    Project = var.project_tag
  }
}

# bucket to hold precious input files
#
resource "aws_s3_bucket" "input-data-files" {
  bucket = "find-insights-input-data-files"
  acl    = "private"
  versioning {
    enabled = true
  }

  tags = {
    Project = var.project_tag
  }
}
