# Minimal Terramate root configuration
# Only contains git settings that must be at project root
# All globals and stack definitions are living in the demo folders; this file is only used in projects that use Terramate.
terramate {
  config {
    git {
      default_branch = "master"
      default_remote = "origin"
    }
  }
}
