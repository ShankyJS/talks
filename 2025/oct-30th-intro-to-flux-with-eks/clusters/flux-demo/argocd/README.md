# ArgoCD Deployment via Flux

This directory contains Flux manifests to deploy ArgoCD using Helm.

## What Gets Deployed

- **Namespace**: `argocd`
- **HelmRepository**: Points to the official Argo Helm chart repository
- **HelmRelease**: Deploys ArgoCD with a LoadBalancer service

## Configuration

The ArgoCD deployment is configured with:
- **Service Type**: LoadBalancer (accessible externally)
- **TLS**: Disabled for demo purposes (`--insecure` flag)
- **Dex SSO**: Disabled
- **Chart Version**: Latest 7.x

## Accessing ArgoCD

Once Flux reconciles and deploys ArgoCD, you can access it:

```bash
# Get the LoadBalancer URL
kubectl get svc argocd-server -n argocd

# Get the admin password
kubectl -n argocd get secret argocd-initial-admin-secret -o jsonpath="{.data.password}" | base64 -d; echo
```

Then access the UI at the LoadBalancer endpoint with:
- **Username**: `admin`
- **Password**: (from the command above)

## How It Works

1. **Flux watches** this directory in your Git repository
2. **HelmRepository** fetches the Argo chart metadata from GitHub
3. **HelmRelease** installs/updates ArgoCD using the specified chart version
4. **Flux automatically reconciles** changes you push to this directory

## Customization

To customize the ArgoCD deployment, edit the `values` section in `helmrelease.yaml`:

```yaml
values:
  global:
    domain: argocd.example.com  # Your domain

  server:
    service:
      type: LoadBalancer  # Or ClusterIP, NodePort
      annotations:
        service.beta.kubernetes.io/aws-load-balancer-type: "nlb"  # Optional
```

Commit and push your changes - Flux will automatically apply them!

## Monitoring

Check the HelmRelease status:
```bash
kubectl get helmrelease -n argocd
flux get helmreleases -n argocd
```

View logs:
```bash
flux logs --follow --level=info
```

## GitOps Inception 🤯

You now have Flux (GitOps tool) deploying ArgoCD (another GitOps tool) on your cluster. This demonstrates:
- Flux's ability to manage Helm charts declaratively
- How easy it is to deploy applications via Git
- The flexibility to use multiple GitOps approaches together
