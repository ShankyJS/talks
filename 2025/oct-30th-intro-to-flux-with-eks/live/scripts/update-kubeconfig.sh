#!/usr/bin/env bash
# Update kubeconfig for EKS cluster
# This script removes old cluster config and downloads the new one

set -e

# Default values - can be overridden by environment variables or arguments
CLUSTER_NAME="${1:-flux-demo}"
REGION="${2:-us-west-2}"
KUBECONFIG_FILE="${KUBECONFIG:-$HOME/.kube/config}"

echo "🧹 Cleaning up kubeconfig for cluster: $CLUSTER_NAME"

# Remove old cluster entry from kubeconfig
if kubectl config get-contexts -o name | grep -q "arn:aws:eks:$REGION:.*:cluster/$CLUSTER_NAME"; then
    CONTEXT_NAME=$(kubectl config get-contexts -o name | grep "arn:aws:eks:$REGION:.*:cluster/$CLUSTER_NAME" || true)
    if [ -n "$CONTEXT_NAME" ]; then
        echo "  → Removing old context: $CONTEXT_NAME"
        kubectl config delete-context "$CONTEXT_NAME" 2>/dev/null || true
    fi
fi

# Remove old cluster entry
if kubectl config get-clusters | grep -q "arn:aws:eks:$REGION:.*:cluster/$CLUSTER_NAME"; then
    CLUSTER_ARN=$(kubectl config get-clusters | grep "arn:aws:eks:$REGION:.*:cluster/$CLUSTER_NAME" || true)
    if [ -n "$CLUSTER_ARN" ]; then
        echo "  → Removing old cluster: $CLUSTER_ARN"
        kubectl config delete-cluster "$CLUSTER_ARN" 2>/dev/null || true
    fi
fi

# Remove old user entry
if kubectl config view -o jsonpath='{.users[*].name}' | grep -q "arn:aws:eks:$REGION:.*:cluster/$CLUSTER_NAME"; then
    USER_ARN=$(kubectl config view -o jsonpath='{.users[*].name}' | tr ' ' '\n' | grep "arn:aws:eks:$REGION:.*:cluster/$CLUSTER_NAME" || true)
    if [ -n "$USER_ARN" ]; then
        echo "  → Removing old user: $USER_ARN"
        kubectl config delete-user "$USER_ARN" 2>/dev/null || true
    fi
fi

echo ""
echo "📥 Downloading new kubeconfig for cluster: $CLUSTER_NAME"

# Update kubeconfig with new cluster credentials
aws eks update-kubeconfig --region "$REGION" --name "$CLUSTER_NAME"

echo ""
echo "✅ Kubeconfig updated successfully!"
echo ""
echo "🔍 Current context:"
kubectl config current-context

echo ""
echo "🎯 Testing connection..."
if kubectl get nodes &>/dev/null; then
    echo "✅ Successfully connected to cluster!"
    echo ""
    kubectl get nodes
else
    echo "⚠️  Could not connect to cluster. Make sure it's deployed and accessible."
    exit 1
fi
