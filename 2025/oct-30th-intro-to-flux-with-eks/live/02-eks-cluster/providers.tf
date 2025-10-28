terraform {
  required_version = ">= 1.0"

  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 5.0"
    }
    kubernetes = {
      source  = "hashicorp/kubernetes"
      version = "~> 2.23"
    }
  }
}

terraform {
  backend "s3" {
    bucket         = "flux-eks-demo-tfstate-0mw37gze"
    dynamodb_table = "flux-eks-demo-tfstate-locks"
    encrypt        = true
    key            = "eks-cluster/terraform.tfstate"
    region         = "us-west-2"
  }
}

provider "aws" {
  region = var.region

  default_tags {
    tags = var.tags
  }
}

# Kubernetes provider configuration
provider "kubernetes" {
  host                   = module.eks.cluster_endpoint
  cluster_ca_certificate = base64decode(module.eks.cluster_certificate_authority_data)

  exec {
    api_version = "client.authentication.k8s.io/v1beta1"
    command     = "aws"
    args        = ["eks", "get-token", "--cluster-name", module.eks.cluster_name]
  }
}
