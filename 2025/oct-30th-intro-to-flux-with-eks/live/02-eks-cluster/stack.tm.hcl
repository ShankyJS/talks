stack {
  name        = "eks-cluster"
  description = "EKS cluster for Flux GitOps demo"
  tags        = ["eks-cluster"]
  after       = ["tag:s3-backend"]
}

# Generate terraform.tfvars from globals
generate_hcl "_generated.terraform.auto.tfvars" {
  content {
    project_name = global.project_name
    environment  = global.environment
    region       = global.aws.region

    # EKS Configuration
    cluster_name    = global.eks.cluster_name
    cluster_version = global.eks.cluster_version

    # VPC Configuration
    vpc_cidr            = global.vpc.cidr
    vpc_azs             = global.vpc.azs
    vpc_private_subnets = global.vpc.private_subnets
    vpc_public_subnets  = global.vpc.public_subnets

    # Node Group Configuration
    node_group_instance_types = global.node_group.instance_types
    node_group_desired_size   = global.node_group.desired_size
    node_group_min_size       = global.node_group.min_size
    node_group_max_size       = global.node_group.max_size

    # SSO Admin Role for EKS access
    admin_sso_role_name = global.aws.sso_admin_role_name

    tags = tm_merge(
      global.tags,
      {
        Stack = "eks-cluster"
      }
    )
  }
}
