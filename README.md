# 🎤 Talks & Demos Repository

[![en](https://img.shields.io/badge/lang-en-red.svg)](./README.md)
[![es](https://img.shields.io/badge/lang-es-yellow.svg)](./docs/README-es.md)

## 👋 About Me

Hi! I'm **Shanky** ([@shankyjs](https://github.com/shankyjs)), a Sr. Platform Engineer passionate about cloud-native technologies, DevOps, and open-source software. I love sharing knowledge through talks, workshops, and demos.

<div align="center">
  <a href="https://github.com/shankyjs">
    <img src="https://github.com/shankyjs.png" width="150" alt="Shanky"/>
  </a>
</div>

**Community Involvement:**
- 🇨🇦 Organizer of [Cloud Native Vancouver](https://community.cncf.io/cloud-native-vancouver/)
- 🇸🇻 Organizer of [Cloud Native San Salvador](https://community.cncf.io/cloud-native-san-salvador/)

## 📚 About This Repository

Welcome to my talks and demos repository! This is where I collect and share all the presentations, demonstrations, and code examples I've created over the years. Whether it's from conferences, meetups, workshops, or community events, you'll find the resources here.

Each talk includes:
- 📝 Presentation materials and slides
- 💻 Demo code and examples
- 📖 Step-by-step instructions
- 🔗 Additional resources and references

## 📊 Statistics

- 🎤 **Total Talks**: 2
- ✅ **Past**: 1
- 🔜 **Upcoming**: 1
- 🏷️ **Top Topics**: AWS (1), GitOps (1), Go (1)

## 📑 Talks Index

Browse all talks by year, topic, and event. Click on any talk to access the full demo, code, and materials.

### 2025

| Date | Talk Title | Topics | Event/Location | Materials |
|------|------------|--------|----------------|-----------|
| 2025-11-19 | [**Otel Jaeger Go Services**](./2025/nov-19th-otel-jaeger-go-services) | Otel, Jaeger, Go | Cloud Native Vancouver: Nov 2025 | [EN](./2025/nov-19th-otel-jaeger-go-services/README.md) / [ES](./2025/nov-19th-otel-jaeger-go-services/README-es.md) |
| 2025-10-30 | [**Intro To Flux With EKS**](./2025/oct-30th-intro-to-flux-with-eks) | GitOps, AWS, Kubernetes | October 30th Cloud Native Vancouver event | [EN](./2025/oct-30th-intro-to-flux-with-eks/README.md) / [ES](./2025/oct-30th-intro-to-flux-with-eks/README-es.md) |


### Coming Soon 🚀

More talks and demos will be added here as they happen!

---

## 🏷️ Browse by Topic

- **AWS**: [Intro To Flux With EKS (2025)](./2025/oct-30th-intro-to-flux-with-eks)
- **GitOps**: [Intro To Flux With EKS (2025)](./2025/oct-30th-intro-to-flux-with-eks)
- **Go**: [Otel Jaeger Go Services (2025)](./2025/nov-19th-otel-jaeger-go-services)
- **Jaeger**: [Otel Jaeger Go Services (2025)](./2025/nov-19th-otel-jaeger-go-services)
- **Kubernetes**: [Intro To Flux With EKS (2025)](./2025/oct-30th-intro-to-flux-with-eks)
- **Otel**: [Otel Jaeger Go Services (2025)](./2025/nov-19th-otel-jaeger-go-services)


## 🤝 Contributing

Found a typo or want to improve something? Feel free to open an issue or submit a pull request!

## 📫 Get in Touch

- GitHub: [@shankyjs](https://github.com/shankyjs)
- Feel free to reach out if you have questions about any of the demos or talks!

## 📄 License

Unless otherwise specified, all content in this repository is available for educational purposes. Please reference this repository if you use any materials.

---

⭐ If you find these resources helpful, consider giving this repository a star!

## Contribute

```bash
# 1. Build automation tools
make build

# This compiles all automation tools:
# - create-talk (create new talk directories)
# - generate-index (update talks index)
# - check-metadata (validate metadata files)
# - generate-stats (generate statistics)

# 2. Install pre-commit hooks (optional but recommended)
pip install pre-commit  # or brew install pre-commit
pre-commit install
```

> **Note**: All automation is built using Go. Run `make build` to compile the binaries.

### Creating a New Talk

```bash
# Use the Makefile command
make create-talk DATE=2025-11-15 SLUG=my-awesome-talk

# Or use the shorter alias
make new DATE=2025-11-15 SLUG=my-awesome-talk

# This creates:
# - 2025/nov-15th-my-awesome-talk/
# - metadata.yaml (edit this!)
# - README.md
# - README-es.md
```

### Updating the Index

```bash
# After creating or editing talks
make update-index

# Or simply
make regen
```

### Pre-commit Hooks

Once installed, pre-commit hooks will:
- ✅ Auto-generate index on commit
- ✅ Validate metadata files
- ✅ Check for missing files
- ✅ Fix trailing whitespace

```bash
# Manual run
pre-commit run --all-files
```

### Quick Commands

```bash
make help           # Show all commands
make build          # Build automation tools
make install        # Alias for build
make create-talk    # Create new talk (requires DATE and SLUG)
make new            # Alias for create-talk
make update-index   # Regenerate index
make generate-stats # Generate statistics
make check          # Validate metadata
make clean          # Cleanup
```

### Example Workflow

```bash
# 1. Create talk
make create-talk DATE=2025-12-10 SLUG=kubernetes-secrets

# 2. Edit metadata
vim 2025/dec-10th-kubernetes-secrets/metadata.yaml

# 3. Add content
vim 2025/dec-10th-kubernetes-secrets/README.md

# 4. Update index
make update-index

# 5. Commit (pre-commit does the rest!)
git add .
git commit -m "feat: Add Kubernetes secrets talk"
```
