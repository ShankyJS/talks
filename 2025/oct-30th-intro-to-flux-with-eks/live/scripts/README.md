# Helper Scripts

This directory contains useful scripts for managing the infrastructure.

## Kubeconfig Management

### `update-kubeconfig.sh`

Cleans up old kubeconfig entries and downloads the new configuration for your EKS cluster.

**Usage:**

```bash
# Use defaults (cluster: flux-demo, region: us-west-2)
./scripts/update-kubeconfig.sh

# Specify cluster name
./scripts/update-kubeconfig.sh my-cluster

# Specify cluster name and region
./scripts/update-kubeconfig.sh my-cluster us-east-1
```

**What it does:**
1. Removes old context, cluster, and user entries for the specified cluster
2. Downloads new kubeconfig from AWS EKS
3. Sets it as the current context
4. Tests the connection

**When to use:**
- After recreating a cluster with the same name
- When kubeconfig gets corrupted or outdated
- After cluster upgrades

### `cleanup-all-kubeconfigs.sh`

Removes specific or all AWS EKS entries from your kubeconfig file.

**Usage:**

```bash
# Clean up only flux-demo cluster (default)
./scripts/cleanup-all-kubeconfigs.sh

# Clean up a specific cluster
./scripts/cleanup-all-kubeconfigs.sh my-cluster-name

# Nuclear option: clean ALL EKS entries
./scripts/cleanup-all-kubeconfigs.sh --all
```

**What it does:**
1. Creates a backup of your current kubeconfig
2. Removes matching EKS contexts, clusters, and users
3. Provides restore instructions

**When to use:**
- After testing multiple clusters (use `--all`)
- When starting fresh
- When kubeconfig is cluttered with old entries
- Before recreating a cluster with the same name

**Safety:**
- Always creates a backup before modifying kubeconfig
- Backup location: `~/.kube/config.backup.YYYYMMDD-HHMMSS`
- Default behavior only affects flux-demo cluster

## Usage in Workflow

### Complete Deployment (First Time)

```bash
# Deploy everything with one command
./scripts/apply-all.sh
```

### Recreating Just the EKS Cluster

```bash
# 1. Destroy old cluster
cd 02-eks-cluster
terraform destroy

# 2. Apply new cluster (backend already configured)
terraform apply

# 3. Update kubeconfig
cd ..
./scripts/update-kubeconfig.sh

# 4. Verify
kubectl get nodes
```

### Recreating Everything

```bash
# Destroy all
./scripts/destroy-all.sh

# Deploy all
./scripts/apply-all.sh
```

### Manual Step-by-Step (if you prefer)

```bash
# 1. Deploy S3 backend
cd 01-s3-backend
terraform init
terraform apply

# 2. Initialize EKS with backend
cd ..
./scripts/init-with-backend.sh

# 3. Deploy EKS
cd 02-eks-cluster
terraform plan
terraform apply

# 4. Update kubeconfig
cd ..
./scripts/update-kubeconfig.sh
```

## Tips

Scripts are already executable, but if needed:

```bash
chmod +x scripts/*.sh
```

Add to your PATH or create aliases:

```bash
# In your ~/.zshrc or ~/.bashrc
alias deploy-demo='~/path/to/talks/2025/oct-30th-intro-to-flux-with-eks/live/scripts/apply-all.sh'
alias destroy-demo='~/path/to/talks/2025/oct-30th-intro-to-flux-with-eks/live/scripts/destroy-all.sh'
alias update-kube='~/path/to/talks/2025/oct-30th-intro-to-flux-with-eks/live/scripts/update-kubeconfig.sh'
```

## Why These Scripts?

**Backend Config Problem:** Terraform backend blocks don't support variables, so we can't use Terramate outputs sharing directly for backend configuration. Instead, we use scripts to pass backend config via `-backend-config` flags during `terraform init`.

**Workflow Automation:** These scripts automate the entire workflow, making it easy to spin up and tear down the infrastructure repeatedly for demos and testing.
