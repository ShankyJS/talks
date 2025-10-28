stack {
  name        = "flux"
  description = "Flux GitOps demo"
  tags        = ["flux"]
  after       = ["tag:eks-cluster"]
}

# Generate terraform.tfvars from globals
generate_hcl "_generated.terraform.auto.tfvars" {
  content {
    project_name = global.project_name
    environment  = global.environment
    region       = global.aws.region

    # GitHub configuration
    github_user       = global.github.user
    github_repository = global.github.repository

    # EKS cluster name (from Stack 02)
    cluster_name = global.eks.cluster_name

    tags = tm_merge(
      global.tags,
      {
        Stack = "flux"
      }
    )
  }
}
