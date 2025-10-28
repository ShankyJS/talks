# Common Variables
variable "project_name" {
  description = "Name of the project"
  type        = string
}

variable "environment" {
  description = "Environment name"
  type        = string
}

variable "region" {
  description = "AWS region"
  type        = string
}

variable "tags" {
  description = "Common tags for all resources"
  type        = map(string)
  default     = {}
}

# EKS Configuration
variable "cluster_name" {
  description = "Name of the EKS cluster"
  type        = string
}

# GitHub Configuration
variable "github_user" {
  description = "GitHub username"
  type        = string
}

variable "github_repository" {
  description = "GitHub repository name"
  type        = string
}
