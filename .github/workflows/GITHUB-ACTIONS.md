# 🤖 GitHub Actions Guide

This document explains the GitHub Actions workflows that automate the talks repository.

## 📋 Available Workflows

### 1. Auto-Update Talks Index 🔄

**File**: `.github/workflows/auto-update-index.yml`

**Trigger**: Automatically runs when metadata files change on master branch

**What it does**:
- Detects changes to `metadata.yaml` files
- Builds Go binaries
- Regenerates the talks index
- Commits and pushes changes if the index is out of sync
- **Perfect safety net** if you forget to run `make update-index` or bypass pre-commit hooks

**Status**: ✅ Fully automated

```yaml
on:
  push:
    branches: [master]
    paths:
      - '**/metadata.yaml'
```

**Technology**: Uses Go binaries for fast, consistent execution

### 2. Validate Talks Metadata ✅

**File**: `.github/workflows/validate-metadata.yml`

**Trigger**: Runs on every PR and push to master

**What it does**:
- Builds Go binaries
- Validates all `metadata.yaml` files
- Checks for required fields (title, date, topics)
- Verifies README files exist
- Posts validation report in PR summary

**Status**: ✅ Fully automated

**Example output**:
```
✅ All talk directories have valid metadata!
```

### 3. Check Index is Synced 🔍

**File**: `.github/workflows/check-index-sync.yml`

**Trigger**: Runs on PRs that modify metadata or READMEs

**What it does**:
- Builds Go binaries
- Regenerates index from metadata
- Compares with committed index
- **Fails the PR** if index is out of sync
- Provides clear instructions to fix

**Status**: ✅ Fully automated

**Why?**: Prevents merging PRs with outdated indexes

### 4. Generate Talk Statistics 📊

**File**: `.github/workflows/generate-stats.yml`

**Trigger**:
- Weekly (Sundays at midnight UTC)
- Manual trigger via workflow_dispatch
- Push to master

**What it does**:
- Builds Go binaries
- Counts total talks
- Breaks down by year, topic, event
- Shows upcoming vs. past talks
- Generates visualizations

**Status**: ✅ Fully automated

**View**: Check GitHub Actions summary for stats

## 🎯 Workflow Strategy

### Protection Flow

```
Developer              Pre-commit           GitHub Actions
    │                      │                      │
    ├─ Edit metadata       │                      │
    ├─ git commit ─────────┤                      │
    │                      ├─ Regenerates index   │
    │                      ├─ Validates           │
    │                      └─ Commits changes     │
    │                                             │
    ├─ git push (bypassed pre-commit) ───────────┤
    │                                             ├─ Auto-update (safety net!)
    │                                             ├─ Regenerates index
    │                                             └─ Auto-commits
    │
    └─ Open PR ──────────────────────────────────┤
                                                  ├─ Validates metadata
                                                  ├─ Checks sync
                                                  └─ Reports status
```

### Fail-Safe Mechanism

1. **Local**: Pre-commit hooks try to update index (Go binaries)
2. **Push**: Auto-update workflow catches missed updates (Go binaries)
3. **PR**: Validation ensures quality before merge (Go binaries)

**Result**: Index is ALWAYS up to date! 🎉

## 🚀 Setup Instructions

### 1. Enable Workflows

Workflows are automatically enabled when you push the `.github/workflows/` directory.

### 2. Repository Settings

For auto-update to work, ensure:

1. Go to **Settings** → **Actions** → **General**
2. Under "Workflow permissions":
   - ✅ Check "Read and write permissions"
   - ✅ Check "Allow GitHub Actions to create and approve pull requests"

### 3. Branch Protection (Optional)

Recommended branch protection rules:

1. Go to **Settings** → **Branches**
2. Add rule for `master`:
   - ✅ Require status checks: `validate`, `check-sync`
   - ✅ Require branches to be up to date

This ensures PRs can't merge with invalid metadata or out-of-sync indexes.

## 🎨 Workflow Badges

Add status badges to your README:

```markdown
[![Auto-Update Index](https://github.com/shankyjs/talks/actions/workflows/auto-update-index.yml/badge.svg)](https://github.com/shankyjs/talks/actions/workflows/auto-update-index.yml)
[![Validate Metadata](https://github.com/shankyjs/talks/actions/workflows/validate-metadata.yml/badge.svg)](https://github.com/shankyjs/talks/actions/workflows/validate-metadata.yml)
```

## 🔧 Manual Triggers

All workflows support manual triggering:

1. Go to **Actions** tab
2. Select workflow
3. Click **Run workflow**
4. Choose branch
5. Click **Run workflow** button

## 📝 Customizing Workflows

### Change Auto-Update Commit Message

Edit `.github/workflows/auto-update-index.yml`:

```yaml
git commit -m "🤖 Auto-update talks index [skip ci]"
```

The `[skip ci]` prevents infinite loops.

### Add Slack Notifications

Add to any workflow:

```yaml
- name: Notify Slack
  uses: slackapi/slack-github-action@v1
  with:
    webhook-url: ${{ secrets.SLACK_WEBHOOK }}
    payload: |
      {
        "text": "New talk added! 🎉"
      }
```

### Add More Validations

Edit `cmd/check-metadata/main.go` to add custom checks:

```go
// Check for minimum description length
if len(metadata.Description) < 50 {
    fmt.Println("❌ Description too short")
    hasError = true
}
```

Then rebuild: `make build`

## 🐛 Troubleshooting

### Workflow Not Running?

1. Check file paths in workflow triggers
2. Verify branch name (`main` vs `master`)
3. Check Actions tab for error messages

### Permission Denied?

1. Check repository Settings → Actions → General
2. Enable "Read and write permissions"
3. Enable "Allow GitHub Actions to create PRs"

### Index Still Out of Sync?

1. Manually run: `make update-index`
2. Commit changes
3. Push to trigger workflows
4. Check workflow logs in Actions tab

### Build Failing?

1. Verify `go.mod` and `go.sum` are committed
2. Ensure `cmd/` directory is committed
3. Review workflow logs for specific errors
4. Try rebuilding locally: `make build`

## 💡 Advanced Use Cases

### Deploy to GitHub Pages

Create `.github/workflows/deploy-pages.yml`:

```yaml
name: Deploy to GitHub Pages

on:
  push:
    branches: [master]

jobs:
  deploy:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4

      - name: Set up Go
        uses: actions/setup-go@v5
        with:
          go-version: '1.22'

      - name: Build binaries
        run: make build

      - name: Generate index
        run: bin/generate-index

      - name: Deploy
        uses: peaceiris/actions-gh-pages@v3
        with:
          github_token: ${{ secrets.GITHUB_TOKEN }}
          publish_dir: ./docs
```

### Notify on New Talks

Add to auto-update workflow:

```yaml
- name: Check if new talk added
  run: |
    if git diff HEAD~1 --name-only | grep -q "metadata.yaml"; then
      echo "new_talk=true" >> $GITHUB_OUTPUT
    fi
  id: check

- name: Notify team
  if: steps.check.outputs.new_talk == 'true'
  run: |
    # Send notification
```

## 📊 Monitoring

### View Workflow Runs

1. Go to **Actions** tab
2. See all workflow runs
3. Click on any run for details
4. View logs, artifacts, summaries

### Workflow Insights

GitHub provides:
- Success/failure rates
- Run duration
- Usage statistics

## 🚀 Technology Stack

All workflows use Go for automation:

- **Fast**: Binaries compile and run in seconds
- **Consistent**: Same binaries locally and in CI/CD
- **Portable**: Works on any platform
- **Simple**: Self-contained executables

## 🎯 Best Practices

1. **Keep workflows simple** - One job = one purpose
2. **Use caching** - Cache Go modules for speed
3. **Add summaries** - Use `$GITHUB_STEP_SUMMARY` for reports
4. **Skip CI when needed** - Use `[skip ci]` in commit messages
5. **Test locally** - Run `make build` and test binaries before pushing
6. **Use `make build`** - Consistent build process locally and in CI

## 📚 Resources

- [GitHub Actions Documentation](https://docs.github.com/en/actions)
- [Workflow Syntax](https://docs.github.com/en/actions/reference/workflow-syntax-for-github-actions)
- [Go GitHub Actions](https://github.com/actions/setup-go)
- [Act - Run workflows locally](https://github.com/nektos/act)

---

**Questions?** Open an issue or check the [Actions tab](https://github.com/shankyjs/talks/actions) for logs!
