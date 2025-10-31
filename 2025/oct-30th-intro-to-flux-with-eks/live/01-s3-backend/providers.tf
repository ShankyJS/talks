terraform {
  required_version = ">= 1.0"

  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 6.18"
    }
    random = {
      source  = "hashicorp/random"
      version = "~> 3.5"
    }
  }
}

terraform {
  backend "s3" {
    bucket         = "flux-eks-demo-tfstate-9qfmn2vg"
    dynamodb_table = "flux-eks-demo-tfstate-locks"
    encrypt        = true
    key            = "s3-backend/terraform.tfstate"
    region         = "us-west-2"
  }
}

provider "aws" {
  region = var.region

  default_tags {
    tags = var.tags
  }
}
