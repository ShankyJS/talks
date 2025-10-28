# Get current AWS account ID
data "aws_caller_identity" "current" {}

# Get current AWS region
data "aws_region" "current" {}

# Locals for common references
locals {
  account_id = data.aws_caller_identity.current.account_id
  region     = data.aws_region.current.name
}

# Get information about the EKS cluster from AWS
data "aws_eks_cluster" "this" {
  name = var.cluster_name
}

# Get the auth info for the EKS cluster (required for Kubernetes provider)
data "aws_eks_cluster_auth" "this" {
  name = var.cluster_name
}

data "onepassword_vault" "this" {
  name = "Private"
}

# Get Github PAT from 1Password
data "onepassword_item" "this" {
  vault = data.onepassword_vault.this.name
  title = "flux-cd-token-github"
}
