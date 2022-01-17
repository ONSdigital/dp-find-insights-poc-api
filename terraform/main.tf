terraform {
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 3.0"
    }
  }

  # The backend bucket for holding state is created out of band.
  # Good idea to enable versioning and lifecycle management.
  backend "s3" {
    bucket  = "find-insights-alpha-tfstate"
    key     = "tfstate"
    region  = "eu-central-1"
    profile = "development"
  }
}

# There are other ways to specify region and authentication, but for now do it this way.
#
provider "aws" {
  region                  = var.project_region
  shared_credentials_file = pathexpand(var.aws_cred_file)
  profile                 = var.aws_profile
}
