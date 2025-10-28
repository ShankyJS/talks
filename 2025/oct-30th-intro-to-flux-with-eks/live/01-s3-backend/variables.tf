variable "project_name" {
  description = "Project name to use for naming resources"
  type        = string
  default     = "flux-eks-demo"
}

variable "environment" {
  description = "Environment name (e.g., dev, staging, prod)"
  type        = string
  default     = "demo"
}

variable "region" {
  description = "AWS region"
  type        = string
  default     = "us-west-2"
}

variable "tags" {
  description = "Common tags to apply to all resources"
  type        = map(string)
  default = {
    Project     = "flux-eks-demo"
    Environment = "demo"
    ManagedBy   = "terraform"
  }
}
