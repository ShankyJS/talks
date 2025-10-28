# Global variables available to all stacks
globals {
  # Project configuration
  project_name = "flux-eks-demo"
  environment  = "demo"

  # AWS configuration
  aws = {
    region = "us-west-2"
    # Account ID will be dynamically retrieved via data source in each stack

    # SSO Admin role for EKS access
    sso_admin_role_name = "AWSReservedSSO_AdministratorAccess_999416fc7acbcbb2"
  }

  # EKS configuration
  eks = {
    cluster_name    = "flux-demo"
    cluster_version = "1.31"
  }

  # VPC configuration
  vpc = {
    cidr            = "10.0.0.0/16"
    azs             = ["us-west-2a", "us-west-2b", "us-west-2c"]
    private_subnets = ["10.0.1.0/24", "10.0.2.0/24", "10.0.3.0/24"]
    public_subnets  = ["10.0.101.0/24", "10.0.102.0/24", "10.0.103.0/24"]
  }

  # Node group configuration
  node_group = {
    instance_types = ["t3.medium"]
    desired_size   = 2
    min_size       = 1
    max_size       = 4
  }

  # GitHub configuration for Flux
  github = {
    user       = "ShankyJS"
    repository = "talks"
  }

  # Common tags for all resources
  tags = {
    Project     = global.project_name
    Environment = global.environment
    ManagedBy   = "terraform"
    Talk        = "intro-to-flux-with-eks"
    Date        = "2025-10-30"
  }

  # Backend configuration
  backend = {
    bucket         = "flux-eks-demo-tfstate-0mw37gze"
    dynamodb_table = "flux-eks-demo-tfstate-locks"
    region         = "us-west-2"
  }
}
