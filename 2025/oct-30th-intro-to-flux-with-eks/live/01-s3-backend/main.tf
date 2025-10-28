# Random suffix to ensure unique bucket name
resource "random_string" "suffix" {
  length  = 8
  special = false
  upper   = false
}

# S3 bucket for Terraform state
module "s3_bucket" {
  source  = "terraform-aws-modules/s3-bucket/aws"
  version = "~> 4.0"

  bucket = "${var.project_name}-tfstate-${random_string.suffix.result}"

  # Block all public access
  block_public_acls       = true
  block_public_policy     = true
  ignore_public_acls      = true
  restrict_public_buckets = true

  # Enable versioning for state file history
  versioning = {
    enabled = true
  }

  # Enable server-side encryption
  server_side_encryption_configuration = {
    rule = {
      apply_server_side_encryption_by_default = {
        sse_algorithm = "AES256"
      }
    }
  }

  # Lifecycle rules to manage old versions
  lifecycle_rule = [
    {
      id      = "expire-old-versions"
      enabled = true

      noncurrent_version_expiration = {
        days = 90
      }
    }
  ]

  tags = merge(
    var.tags,
    {
      Name    = "${var.project_name}-tfstate-${random_string.suffix.result}"
      Purpose = "Terraform State Storage"
    }
  )
}

# DynamoDB table for state locking
resource "aws_dynamodb_table" "terraform_locks" {
  name         = "${var.project_name}-tfstate-locks"
  billing_mode = "PAY_PER_REQUEST"
  hash_key     = "LockID"

  attribute {
    name = "LockID"
    type = "S"
  }

  tags = merge(
    var.tags,
    {
      Name    = "${var.project_name}-tfstate-locks"
      Purpose = "Terraform State Locking"
    }
  )
}
