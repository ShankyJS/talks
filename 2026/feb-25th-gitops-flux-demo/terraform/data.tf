# Get current AWS account info
data "aws_caller_identity" "current" {}
data "aws_region" "current" {}

# EKS cluster auth (for Kubernetes/Flux providers)
data "aws_eks_cluster_auth" "this" {
  name = module.eks.cluster_name
}

# 1Password - GitHub PAT for Flux
data "onepassword_vault" "this" {
  name = local.onepassword_vault
}

data "onepassword_item" "github_token" {
  vault = data.onepassword_vault.this.name
  title = local.onepassword_item
}
