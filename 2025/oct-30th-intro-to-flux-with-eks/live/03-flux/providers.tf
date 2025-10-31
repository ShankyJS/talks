terraform {
  required_version = ">= 1.0"

  required_providers {
    flux = {
      source  = "fluxcd/flux"
      version = "1.7.4"
    }
    onepassword = {
      source  = "1Password/onepassword"
      version = "2.2.0"
    }
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
    key            = "flux/terraform.tfstate"
    region         = "us-west-2"
  }
}


provider "onepassword" {
  account = "my.1password.ca"
}


provider "aws" {
  region = var.region

  default_tags {
    tags = var.tags
  }
}

provider "flux" {
  kubernetes = {
    host                   = data.aws_eks_cluster.this.endpoint
    token                  = data.aws_eks_cluster_auth.this.token
    cluster_ca_certificate = base64decode(data.aws_eks_cluster.this.certificate_authority[0].data)
  }
  git = {
    url    = "https://github.com/${var.github_user}/${var.github_repository}.git"
    branch = "master"
    http = {
      username = "git" # This can be any string when using a personal access token
      password = data.onepassword_item.this.credential
    }
  }
}
