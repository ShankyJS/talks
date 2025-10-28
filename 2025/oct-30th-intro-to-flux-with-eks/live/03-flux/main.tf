resource "flux_bootstrap_git" "this" {
  embedded_manifests = true
  path               = "2025/oct-30th-intro-to-flux-with-eks/clusters/flux-demo"
}
