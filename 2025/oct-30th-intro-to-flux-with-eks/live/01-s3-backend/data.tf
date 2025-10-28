# Get current AWS account ID
data "aws_caller_identity" "current" {}

# Get current AWS region
data "aws_region" "current" {}

# Locals for common references
locals {
  account_id = data.aws_caller_identity.current.account_id
  region     = data.aws_region.current.name
}
