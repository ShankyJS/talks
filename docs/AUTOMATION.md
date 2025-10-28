# 🤖 Automation Guide

This document explains the automation system for managing the talks index.

## 📋 Overview

The talks repository uses an automated system to generate the index from metadata files. This means you never have to manually update the README index tables - they're generated automatically!

## 🗂️ How It Works

### 1. Metadata Files

Each talk directory contains a `metadata.yaml` file with talk information:

```yaml
title: "Flux with AWS"
date: "2025-10-30"
event: "Cloud Native Vancouver"
topics:
  - GitOps
  - Flux CD
  - AWS
  - Kubernetes
description: "Demo of Flux CD on AWS"
```

### 2. Index Generation

The `bin/generate-index` binary:
- Scans all year directories (2024, 2025, etc.)
- Reads metadata from each talk
- Generates markdown tables
- Updates both English and Spanish README files
- Generates statistics section

### 3. Pre-commit Hooks

When you commit changes, pre-commit hooks automatically:
- Validate metadata files
- Regenerate the index
- Check for missing files

## 🚀 Setup

### Initial Setup

```bash
# Build automation tools
make build

# This compiles all automation tools:
# - bin/create-talk (create new talk directories)
# - bin/generate-index (regenerate talks index)
# - bin/check-metadata (validate metadata files)
# - bin/generate-stats (generate statistics)

# Install pre-commit hooks
pip install pre-commit  # or brew install pre-commit
pre-commit install
```

### Creating a New Talk

Use the Makefile command:

```bash
make create-talk DATE=2025-11-15 SLUG=kubernetes-scaling
```

This creates:
- Talk directory with proper naming
- Template metadata.yaml
- Template README.md and README-es.md

## 📝 Daily Usage

### Adding a New Talk

1. Use the command: `make create-talk DATE=2025-11-15 SLUG=kubernetes-scaling`
2. Edit `metadata.yaml` with talk details
3. Add your content to README files
4. Regenerate index: `make update-index`

### Updating an Existing Talk

1. Edit the `metadata.yaml` file
2. Run `make update-index` to regenerate the index
3. Commit changes (pre-commit will also update the index)

### Checking for Issues

```bash
# Check all talks have valid metadata
make check
```

### Generating Statistics

```bash
# Generate detailed talk statistics
make generate-stats
# or
make stats
```

## 🛠️ Available Commands

```bash
make help           # Show all available commands
make build          # Build automation tools
make install        # Alias for build
make create-talk    # Create new talk (requires DATE and SLUG)
make new-talk       # Alias for create-talk
make new            # Short alias for create-talk
make update-index   # Regenerate talks index
make generate-stats # Generate statistics
make stats          # Alias for generate-stats
make check          # Verify metadata files
make clean          # Remove generated files
make regen          # Alias for update-index
```

## 📁 Required File Structure

```
talks/
├── 2025/
│   └── oct-30th-flux-with-aws/
│       ├── metadata.yaml          # Required
│       ├── README.md              # Required
│       ├── README-es.md           # Required
│       └── [demo files...]
├── cmd/
│   ├── create-talk/               # Go source code
│   │   ├── main.go
│   │   └── templates/             # Go templates
│   ├── generate-index/
│   │   └── main.go
│   ├── check-metadata/
│   │   └── main.go
│   └── generate-stats/
│       └── main.go
├── bin/                           # Compiled binaries (gitignored)
│   ├── create-talk
│   ├── generate-index
│   ├── check-metadata
│   └── generate-stats
├── Makefile                       # Commands
└── .pre-commit-config.yaml        # Git hooks
```

## 🎯 Metadata Fields

### Required Fields

- `title`: Talk title
- `date`: Date in YYYY-MM-DD format
- `topics`: List of topics/technologies
- `event`: Event name or "TBA"

### Optional Fields

- `description`: Brief description
- `slides_url`: Link to slides
- `video_url`: Link to recording

## 🔄 Workflow Example

```bash
# 1. Create new talk
make create-talk DATE=2025-12-01 SLUG=docker-best-practices

# 2. Edit metadata
vim 2025/dec-1st-docker-best-practices/metadata.yaml

# 3. Add your content
vim 2025/dec-1st-docker-best-practices/README.md

# 4. Regenerate index
make update-index

# 5. Generate stats (optional)
make generate-stats

# 6. Commit (pre-commit will validate)
git add .
git commit -m "Add Docker best practices talk"
```

## 🐛 Troubleshooting

### Index not updating?

```bash
# Manually regenerate
bin/generate-index
```

### Pre-commit not running?

```bash
# Reinstall hooks
pre-commit install
pre-commit run --all-files
```

### Metadata validation errors?

```bash
# Check what's wrong
make check
```

### Binaries not found?

```bash
# Rebuild tools
make build
```

## 💡 Tips

1. **Always use the Makefile** to create new talks - ensures consistent naming
2. **Run `make update-index`** after any metadata changes
3. **Use pre-commit hooks** to catch issues before pushing
4. **Keep metadata.yaml simple** - the automation handles the rest
5. **Index generation is fast** - nearly instant

## 📚 Adding More Languages

To add a new language (e.g., French):

1. Update `cmd/generate-index/main.go` to support the new language
2. Add language-specific strings
3. Create `docs/README-fr.md`
4. Update the program to generate that README

## 🎨 Customizing

The automation system is flexible. You can customize:

- Table format in `cmd/generate-index/main.go`
- Required metadata fields in `cmd/check-metadata/main.go`
- Pre-commit hooks in `.pre-commit-config.yaml`
- Template content in `cmd/create-talk/templates/`

## 🤝 Contributing

If you improve the automation system, please:
1. Update this guide
2. Rebuild tools: `make build`
3. Test thoroughly
