stack {
  name        = "s3-backend"
  description = "S3 bucket for Terraform remote state backend"
  tags        = ["s3-backend"]
  after       = []
}

# Generate terraform.tfvars from globals
generate_hcl "_generated.terraform.auto.tfvars" {
  content {
    project_name = global.project_name
    environment  = global.environment
    region       = global.aws.region

    tags = tm_merge(
      global.tags,
      {
        Stack = "s3-backend"
      }
    )
  }
}
