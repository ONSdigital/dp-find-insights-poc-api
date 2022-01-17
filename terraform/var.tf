# Using Frankfurt to avoid confusing people and to provide a bit of isolation.
#
variable "project_region" {
  default = "eu-central-1"
}

# When possible, create a "Project" tag on all resources associated with this project.
# (Unfortunately I don't have permissions to tag certain resources, so I can't use
# default tags; have to tag resources individually.)
#
variable "project_tag" {
  default = "find-insights-alpha"
}

# Assuming we will be using credentials file and profile for authentication during dev.
#
variable "aws_cred_file" {
  default = "~/.aws/credentials"
}
variable "aws_profile" {
  default = "development"
}