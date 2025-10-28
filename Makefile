.PHONY: help build install create-talk new-talk update-index check generate-stats clean stats regen new

# Binary locations
BIN_DIR = bin
CREATE_TALK = $(BIN_DIR)/create-talk
GENERATE_INDEX = $(BIN_DIR)/generate-index
CHECK_METADATA = $(BIN_DIR)/check-metadata
GENERATE_STATS = $(BIN_DIR)/generate-stats

help: ## Show this help message
	@echo "📚 Talks Repository - Available Commands"
	@echo ""
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2}'
	@echo ""
	@echo "📝 Usage Examples:"
	@echo "  make create-talk DATE=2025-11-15 SLUG=kubernetes-scaling"
	@echo "  make new DATE=2026-01-20 SLUG=docker-security"

build: ## Build all Go binaries
	@echo "🔨 Building Go binaries..."
	@mkdir -p $(BIN_DIR)
	@go build -o $(CREATE_TALK) ./cmd/create-talk
	@go build -o $(GENERATE_INDEX) ./cmd/generate-index
	@go build -o $(CHECK_METADATA) ./cmd/check-metadata
	@go build -o $(GENERATE_STATS) ./cmd/generate-stats
	@echo "✅ Binaries built in $(BIN_DIR)/"

install: build ## Build binaries (alias for build)
	@echo "✅ Installation complete!"
	@echo ""
	@echo "💡 Binaries are located in $(BIN_DIR)/"

$(CREATE_TALK):
	@$(MAKE) build

$(GENERATE_INDEX):
	@$(MAKE) build

$(CHECK_METADATA):
	@$(MAKE) build

$(GENERATE_STATS):
	@$(MAKE) build

create-talk: $(CREATE_TALK) ## Create a new talk (requires DATE=YYYY-MM-DD SLUG=talk-name)
ifndef DATE
	@echo "❌ Error: DATE is required"
	@echo "Usage: make create-talk DATE=2025-11-15 SLUG=kubernetes-scaling"
	@exit 1
endif
ifndef SLUG
	@echo "❌ Error: SLUG is required"
	@echo "Usage: make create-talk DATE=2025-11-15 SLUG=kubernetes-scaling"
	@exit 1
endif
	@echo "🎤 Creating new talk..."
	@$(CREATE_TALK) -date $(DATE) -slug $(SLUG)

new-talk: create-talk ## Alias for create-talk

update-index: $(GENERATE_INDEX) ## Regenerate the talks index from metadata files
	@echo "🔄 Regenerating talks index..."
	@$(GENERATE_INDEX)
	@echo "✅ Index updated"

generate-stats: $(GENERATE_STATS) ## Generate the talk statistics
	@echo "🔄 Generating talk statistics..."
	@$(GENERATE_STATS)
	@echo "✅ Statistics generated"

check: $(CHECK_METADATA) ## Verify all talks have metadata files
	@echo "🔍 Checking for missing metadata files..."
	@$(CHECK_METADATA)

clean: ## Remove generated files and binaries
	@echo "🧹 Cleaning up..."
	@rm -rf $(BIN_DIR)
	@echo "✅ Cleanup complete"

# Quick aliases
regen: update-index ## Alias for update-index
new: create-talk ## Alias for create-talk
stats: generate-stats ## Alias for generate-stats
