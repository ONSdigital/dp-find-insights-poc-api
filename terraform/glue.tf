resource "aws_iam_role" "fi-glue-execution" {
  name        = "find-insights-glue-execution"
  description = "role assumed by glue scripts at runtime"

  assume_role_policy = data.aws_iam_policy_document.glue_assume_role.json
}

data "aws_iam_policy_document" "glue_assume_role" {
  statement {
    actions = ["sts:AssumeRole"]

    principals {
      type        = "Service"
      identifiers = ["glue.amazonaws.com"]
    }
  }
}

resource "aws_iam_role_policy_attachment" "glue_can_do_glue" {
  role       = aws_iam_role.fi-glue-execution.name
  policy_arn = "arn:aws:iam::aws:policy/service-role/AWSGlueServiceRole"
}
