#!/usr/bin/env bash
# Clean up specific EKS cluster entries from kubeconfig
# Default: Only cleans flux-demo cluster
# Use --all flag to clean ALL EKS clusters

set -e

KUBECONFIG_FILE="${KUBECONFIG:-$HOME/.kube/config}"
CLUSTER_NAME="${1:-flux-demo}"
CLEAN_ALL=false

# Check for --all flag
if [ "$1" = "--all" ]; then
    CLEAN_ALL=true
    CLUSTER_NAME=""
fi

if [ "$CLEAN_ALL" = true ]; then
    echo "🧹 Cleaning up ALL AWS EKS entries from kubeconfig..."
    SEARCH_PATTERN='arn:aws:eks:'
else
    echo "🧹 Cleaning up kubeconfig entries for cluster: $CLUSTER_NAME"
    SEARCH_PATTERN="cluster/$CLUSTER_NAME"
fi

echo "   Kubeconfig file: $KUBECONFIG_FILE"
echo ""

# Backup current kubeconfig
BACKUP_FILE="${KUBECONFIG_FILE}.backup.$(date +%Y%m%d-%H%M%S)"
cp "$KUBECONFIG_FILE" "$BACKUP_FILE"
echo "📦 Backup created: $BACKUP_FILE"
echo ""

# Get matching EKS contexts
EKS_CONTEXTS=$(kubectl config get-contexts -o name | grep "$SEARCH_PATTERN" || true)

if [ -z "$EKS_CONTEXTS" ]; then
    echo "ℹ️  No matching EKS contexts found in kubeconfig"
else
    echo "🗑️  Removing EKS contexts:"
    for context in $EKS_CONTEXTS; do
        echo "  → $context"
        kubectl config delete-context "$context" 2>/dev/null || true
    done
    echo ""
fi

# Get matching EKS clusters
EKS_CLUSTERS=$(kubectl config get-clusters | grep "$SEARCH_PATTERN" || true)

if [ -z "$EKS_CLUSTERS" ]; then
    echo "ℹ️  No matching EKS clusters found in kubeconfig"
else
    echo "🗑️  Removing EKS clusters:"
    for cluster in $EKS_CLUSTERS; do
        echo "  → $cluster"
        kubectl config delete-cluster "$cluster" 2>/dev/null || true
    done
    echo ""
fi

# Get matching EKS users
EKS_USERS=$(kubectl config view -o jsonpath='{.users[*].name}' | tr ' ' '\n' | grep "$SEARCH_PATTERN" || true)

if [ -z "$EKS_USERS" ]; then
    echo "ℹ️  No matching EKS users found in kubeconfig"
else
    echo "🗑️  Removing EKS users:"
    for user in $EKS_USERS; do
        echo "  → $user"
        kubectl config delete-user "$user" 2>/dev/null || true
    done
    echo ""
fi

echo "✅ Cleanup complete!"
echo ""
echo "💡 To restore your previous config:"
echo "   cp $BACKUP_FILE $KUBECONFIG_FILE"
