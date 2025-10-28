# Intro To Flux With EKS

[![en](https://img.shields.io/badge/lang-en-red.svg)](./README.md)
[![es](https://img.shields.io/badge/lang-es-yellow.svg)](./README-es.md)

## 📅 Talk Information

- **Date**: 2025-10-30
- **Event**: GitOps & Kubernetes Meetup
- **Topics**: GitOps, Flux, EKS, Terramate, Infrastructure as Code

## 📝 Description

This demo showcases a complete GitOps workflow using Flux on Amazon EKS, orchestrated with Terramate for infrastructure management. Learn how to set up automated Kubernetes deployments where Git becomes the single source of truth for your infrastructure and applications.

## 🎯 What You'll Learn

- How to provision EKS clusters with Terraform and manage them with Terramate stacks
- Setting up Flux for GitOps continuous delivery on Kubernetes
- Managing secrets securely (1Password integration example - adaptable to other providers)
- Organizing multi-stack infrastructure with proper dependencies
- Best practices for Infrastructure as Code with state management

## 🏗️ Architecture Overview

The demo is organized into 3 Terramate stacks that build upon each other:

```
┌─────────────────────────────────────────────────────────────┐
│ Stack 01: S3 Backend                                        │
│ ├── S3 Bucket (for Terraform state)                        │
│ └── DynamoDB Table (for state locking)                     │
└─────────────────────────────────────────────────────────────┘
                            ↓
┌─────────────────────────────────────────────────────────────┐
│ Stack 02: EKS Cluster                                       │
│ ├── VPC with public/private subnets                        │
│ ├── EKS Control Plane (v1.32)                              │
│ ├── Managed Node Group (t3.medium)                         │
│ ├── EBS CSI Driver                                         │
│ └── SSO Access Entry (AdministratorAccess)                 │
└─────────────────────────────────────────────────────────────┘
                            ↓
┌─────────────────────────────────────────────────────────────┐
│ Stack 03: Flux GitOps                                       │
│ ├── Flux Controllers (source, kustomize, helm, notification)│
│ ├── Git Repository Integration (GitHub)                    │
│ └── Secret Management (1Password provider - customizable)  │
└─────────────────────────────────────────────────────────────┘
```

### How Components Connect

1. **Terramate** orchestrates the deployment order and shares configuration via globals
2. **Stack 01** creates the S3 backend that Stacks 02 & 03 use for remote state
3. **Stack 02** provisions the EKS cluster and grants your SSO role admin access
4. **Stack 03** bootstraps Flux onto the cluster, connecting it to your Git repository
5. **Flux** watches your Git repo and automatically applies changes to the cluster

## 🚀 Demo

### Prerequisites

- [Terraform](https://www.terraform.io/) v1.12+ (ARM64 recommended for Apple Silicon)
- [Terramate](https://terramate.io/) CLI
- [AWS CLI](https://aws.amazon.com/cli/) configured with SSO
- [kubectl](https://kubernetes.io/docs/tasks/tools/)
- AWS Account with EKS permissions

**Secret Management** (choose one):
- [1Password CLI](https://developer.1password.com/docs/cli) (default in this demo)
- Or modify `03-flux/providers.tf` to use your preferred secret provider (AWS Secrets Manager, Vault, etc.)

### Configuration

1. **Update globals** in `live/terramate.tm.hcl`:
   ```hcl
   aws = {
     region = "us-west-2"  # Your AWS region
     sso_admin_role_name = "AWSReservedSSO_AdministratorAccess_xxxxx"  # Your SSO role
   }

   github = {
     user       = "YourGitHubUser"
     repository = "your-repo"
   }
   ```

2. **Store GitHub PAT** (if using 1Password):
   - Store your GitHub Personal Access Token in 1Password
   - Update the reference in `03-flux/main.tf` if needed

### Running the Demo

```bash
cd live/

terramate generate # Generate your config
# Initialize and apply all stacks in order
terramate run -- sh -c 'terraform init && terraform apply'

# Or apply each stack individually:
terramate run --tags s3-backend -- terraform apply
terramate run --tags eks-cluster -- terraform apply
terramate run --tags flux -- terraform apply

./scripts/update-kubeconfigs.sh
# Verify Flux installation
kubectl get pods -n flux-system

# Cleanup (destroys in reverse order)
terramate run --reverse -- terraform destroy --auto-approve
```

## 📁 Repository Structure

```
intro-to-flux-with-eks/
├── README.md              # This file
├── README-es.md           # Spanish version
├── metadata.yaml          # Talk metadata
└── live/                  # Live demo infrastructure
    ├── terramate.tm.hcl   # Global configuration
    ├── 01-s3-backend/     # S3 + DynamoDB for state
    ├── 02-eks-cluster/    # EKS cluster + VPC
    ├── 03-flux/           # Flux GitOps setup
    └── scripts/           # Helper scripts
└── /clusters/flux-demo/   # Flux Cluster spec
```

## 📚 Resources

- [Flux Documentation](https://fluxcd.io/docs/)
- [Terramate Documentation](https://terramate.io/docs/)
- [EKS Best Practices](https://aws.github.io/aws-eks-best-practices/)
- [Terraform AWS Provider](https://registry.terraform.io/providers/hashicorp/aws/latest/docs)
- [GitOps Principles](https://opengitops.dev/)

## 🎥 Recording

Recording will be available after the talk.

## 📧 Feedback

Questions or feedback? Feel free to [open an issue](https://github.com/shankyjs/talks/issues) or reach out!

---

**Presented by**: [@shankyjs](https://github.com/shankyjs)
