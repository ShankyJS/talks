# Flux with EKS - Infrastructure Setup

This directory contains Terraform configurations organized into Terramate stacks for deploying an EKS cluster with Flux GitOps.

## Prerequisites

- [Terraform](https://www.terraform.io/downloads.html) >= 1.0
- [Terramate](https://terramate.io/docs/cli/installation/) (required for generating configuration from globals)
- [AWS CLI](https://aws.amazon.com/cli/) configured with appropriate credentials
- [kubectl](https://kubernetes.io/docs/tasks/tools/)
- [flux CLI](https://fluxcd.io/flux/installation/)

## Terramate Globals

This project uses **Terramate globals** to share common configuration across all stacks. All configuration values are defined in the `terramate.tm.hcl` file:

- **AWS Configuration**: Region, account information
- **Project Settings**: Project name, environment, tags
- **EKS Settings**: Cluster name, version
- **VPC Settings**: CIDR blocks, subnets, availability zones
- **Node Group Settings**: Instance types, sizing

Each stack automatically generates its `terraform.auto.tfvars` file from these globals using Terramate's code generation feature.

### Customizing Values

To change any configuration value:

1. Edit the `terramate.tm.hcl` file in the root directory
2. Run `terramate generate` to regenerate the tfvars files
3. Apply changes with Terraform

**Example:** To change the AWS region for all stacks:

```hcl
# In terramate.tm.hcl
globals {
  aws = {
    region = "us-east-1"  # Changed from us-west-2
  }
}
```

Then run:
```bash
terramate generate
```

## Stack Architecture

The infrastructure is divided into two stacks that must be applied in order:

### Stack 01: S3 Backend (`01-s3-backend/`)
- Creates an S3 bucket for storing Terraform state
- Creates a DynamoDB table for state locking
- Uses local state (no remote backend)
- **Must be applied first**

### Stack 02: EKS Cluster (`02-eks-cluster/`)
- Creates a VPC with public and private subnets
- Deploys an EKS cluster with managed node groups
- Configures EBS CSI driver with IRSA
- Uses S3 backend from Stack 01
- **Must be applied after Stack 01**

## Deployment Steps

### Step 0: Generate Configuration from Globals

Before deploying, generate the Terraform variable files from Terramate globals:

```bash
# From the live/ directory
terramate generate

# This will create _generated.terraform.auto.tfvars files in each stack
```

The generated files will contain all configuration values from the global definitions.

### Step 1: Deploy S3 Backend

```bash
cd 01-s3-backend

# Initialize Terraform
terraform init

# Review the plan
terraform plan

# Apply the configuration
terraform apply

# Save the outputs - you'll need these for Stack 02
terraform output -json > ../backend-outputs.json
```

**Important:** Note the following outputs:
- `s3_bucket_name` - Use this for the S3 backend configuration
- `dynamodb_table_name` - Use this for state locking
- `aws_account_id` - Your AWS account ID
- `aws_region` - The AWS region being used

### Step 2: Configure Backend for EKS Stack

Edit `02-eks-cluster/providers.tf` and uncomment the backend configuration:

```hcl
backend "s3" {
  bucket         = "<s3_bucket_name_from_step_1>"
  key            = "eks-cluster/terraform.tfstate"
  region         = "us-west-2"
  dynamodb_table = "<dynamodb_table_name_from_step_1>"
  encrypt        = true
}
```

### Step 3: Deploy EKS Cluster

```bash
cd ../02-eks-cluster

# Initialize Terraform (will configure S3 backend)
terraform init

# Review the plan
terraform plan

# Apply the configuration
terraform apply
```

### Step 4: Configure kubectl

After the EKS cluster is created, configure kubectl:

```bash
# Use the output from terraform
terraform output configure_kubectl

# Or run directly:
aws eks update-kubeconfig --region us-west-2 --name flux-demo

# Verify connection
kubectl get nodes
```

## Using Terramate for Full Stack Management

Terramate makes it easy to manage both stacks together with proper ordering:

```bash
# From the live/ directory

# List all stacks and their order
terramate list

# Generate configuration from globals
terramate generate

# Run terraform init in all stacks
terramate run terraform init

# Run terraform plan in all stacks (respects order)
terramate run terraform plan

# Run terraform apply in all stacks (respects order)
# Note: You'll still need to manually configure the backend after stack 01
terramate run terraform apply
```

Terramate will automatically execute stacks in the correct order based on the `after` dependencies defined in each `stack.tm.hcl` file.

## Customization

All customizable values are defined in the root `terramate.tm.hcl` file under the `globals` block. After making changes, regenerate configuration files:

### Available Global Configuration

**AWS Settings:**
```hcl
globals {
  aws = {
    region = "us-west-2"  # Change AWS region for all stacks
  }
}
```

**EKS Settings:**
```hcl
globals {
  eks = {
    cluster_name    = "flux-demo"
    cluster_version = "1.31"  # Kubernetes version
  }
}
```

**VPC Settings:**
```hcl
globals {
  vpc = {
    cidr            = "10.0.0.0/16"
    azs             = ["us-west-2a", "us-west-2b", "us-west-2c"]
    private_subnets = ["10.0.1.0/24", "10.0.2.0/24", "10.0.3.0/24"]
    public_subnets  = ["10.0.101.0/24", "10.0.102.0/24", "10.0.103.0/24"]
  }
}
```

**Node Group Settings:**
```hcl
globals {
  node_group = {
    instance_types = ["t3.medium"]
    desired_size   = 2
    min_size       = 1
    max_size       = 4
  }
}
```

### Scaling Example

To change the node group size:

1. Edit `terramate.tm.hcl`:
```hcl
globals {
  node_group = {
    desired_size = 3  # Changed from 2
    min_size     = 2  # Changed from 1
    max_size     = 6  # Changed from 4
  }
}
```

2. Regenerate configuration:
```bash
terramate generate
```

3. Apply changes:
```bash
cd 02-eks-cluster
terraform apply
```

## Cleanup

**Important:** Destroy resources in reverse order!

### Step 1: Destroy EKS Cluster

```bash
cd 02-eks-cluster
terraform destroy
```

### Step 2: Destroy S3 Backend

```bash
cd ../01-s3-backend
terraform destroy
```

**Note:** If you get an error destroying the S3 bucket (because it contains state files), you may need to:

1. Empty the bucket first:
   ```bash
   aws s3 rm s3://<bucket-name> --recursive
   ```
2. Then run `terraform destroy` again

## Outputs

### Stack 01 Outputs
- `s3_bucket_name` - S3 bucket for Terraform state
- `dynamodb_table_name` - DynamoDB table for state locking
- `backend_config` - Complete backend configuration

### Stack 02 Outputs
- `cluster_name` - EKS cluster name
- `cluster_endpoint` - Kubernetes API endpoint
- `configure_kubectl` - Command to configure kubectl
- `vpc_id` - VPC ID
- `oidc_provider_arn` - OIDC provider ARN (for IRSA)
- Many more (see `outputs.tf`)

## Next Steps

After the infrastructure is deployed:

1. **Install Flux:**
   ```bash
   flux install
   ```

2. **Bootstrap Flux with your Git repository:**
   ```bash
   flux bootstrap github \
     --owner=<your-github-username> \
     --repository=<your-repo> \
     --branch=main \
     --path=./clusters/flux-demo \
     --personal
   ```

3. **Deploy applications using GitOps!**

## Troubleshooting

### Issue: Backend initialization fails

**Solution:** Make sure you've applied Stack 01 first and updated the backend configuration in `02-eks-cluster/providers.tf` with the correct bucket name.

### Issue: kubectl can't connect to cluster

**Solution:** Run the configure_kubectl command from the terraform outputs:
```bash
aws eks update-kubeconfig --region us-west-2 --name flux-demo
```

### Issue: Node group fails to create

**Solution:** Check that your AWS account has sufficient EC2 service quotas for the instance types you're using.

## Cost Considerations

Running this infrastructure will incur AWS costs:
- EKS cluster: ~$0.10/hour
- EC2 instances (t3.medium × 2): ~$0.08/hour
- NAT Gateway: ~$0.045/hour
- Other resources: minimal costs

**Estimated total: ~$0.25/hour or ~$180/month**

Remember to destroy resources when not in use!

## Resources

- [EKS Best Practices](https://aws.github.io/aws-eks-best-practices/)
- [Flux Documentation](https://fluxcd.io/docs/)
- [Terraform AWS EKS Module](https://registry.terraform.io/modules/terraform-aws-modules/eks/aws/latest)
- [Terramate Documentation](https://terramate.io/docs/)
