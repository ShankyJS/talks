# Makefile for OTEL, Jaeger, and Go Services

# Variables
KUBECONFIG := $(shell pwd)/.kube/config
HELMFILE_IMAGE := ghcr.io/helmfile/helmfile:v1.2.0

.PHONY: all cluster-up deploy cluster-down clean help kubeconfig-export port-forward-jaeger port-forward-app

help: ## Show this help message
	@echo "📚 OTEL & Jaeger Setup - Available Commands"
	@echo ""
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-20s\033[0m %s\n", $$1, $$2}'
	@echo ""

all: cluster-up deploy ## Create cluster and deploy all services

# Create the Kind cluster
cluster-up: ## Create the Kind cluster
	@echo "Creating Kind cluster..."
	mkdir -p .kube
	kind create cluster --config ./conf/kind-config.yaml --kubeconfig $(KUBECONFIG)

# Deploy services using Helmfile (running in Docker)
deploy: ## Deploy services using Helmfile
	@echo "Deploying services..."
	docker run --rm --net=host \
		-v $(shell pwd):/wd \
		-w /wd \
		-e KUBECONFIG=/wd/.kube/config \
		$(HELMFILE_IMAGE) \
		helmfile apply --file conf/helmfile.yaml

# Build Docker images for backend and frontend
build-images: ## Build Docker images for backend and frontend
	@echo "Building backend image..."
	docker build -t backend:latest ./apps/backend
	@echo "Building frontend image..."
	docker build -t frontend:latest ./apps/frontend

# Load images into Kind cluster
load-images: build-images ## Build and load images into Kind cluster
	@echo "Loading backend image into Kind..."
	kind load docker-image backend:latest
	@echo "Loading frontend image into Kind..."
	kind load docker-image frontend:latest

# Build, load images, and deploy
deploy-apps: load-images deploy ## Build images, load into Kind, and deploy

# Destroy the Kind cluster
cluster-down: ## Destroy the Kind cluster
	@echo "Destroying Kind cluster..."
	kind delete cluster --kubeconfig $(KUBECONFIG)
	rm -rf .kube

# Export kubeconfig to default location (~/.kube/config)
kubeconfig-export: ## Export kubeconfig to default location (~/.kube/config)
	kind export kubeconfig --name kind

# Port forward Jaeger UI
port-forward-jaeger: ## Port forward Jaeger UI to http://localhost:16686
	@echo "Jaeger UI available at http://localhost:16686"
	kubectl --kubeconfig $(KUBECONFIG) port-forward -n monitoring $$(kubectl --kubeconfig $(KUBECONFIG) get pod -n monitoring -l app.kubernetes.io/component=all-in-one -o jsonpath='{.items[0].metadata.name}') 16686:16686

# Port forward Frontend App
port-forward-app: ## Port forward Frontend App to http://localhost:8080
	@echo "Frontend App available at http://localhost:8080"
	kubectl --kubeconfig $(KUBECONFIG) port-forward svc/frontend-generic-service 8080:80

# Connect to Cassandra shell
cassandra-shell: ## Connect to Cassandra CQL shell
	@echo "Connecting to Cassandra..."
	kubectl --kubeconfig $(KUBECONFIG) exec -it -n monitoring jaeger-cassandra-0 -- cqlsh

# Check Jaeger data in Cassandra
cassandra-check-data: ## Check if Jaeger data exists in Cassandra
	@echo "Checking Jaeger keyspace and tables..."
	kubectl --kubeconfig $(KUBECONFIG) exec -n monitoring jaeger-cassandra-0 -- cqlsh -e "DESCRIBE KEYSPACES;"
	@echo ""
	@echo "Checking for traces in Jaeger keyspace..."
	kubectl --kubeconfig $(KUBECONFIG) exec -n monitoring jaeger-cassandra-0 -- cqlsh -e "USE jaeger_v1_test; SELECT COUNT(*) FROM traces;"

# Wait for Cassandra to be ready
cassandra-wait: ## Wait for Cassandra to be ready
	@echo "Waiting for Cassandra to be ready..."
	kubectl --kubeconfig $(KUBECONFIG) wait --for=condition=ready pod -n monitoring -l app=cassandra --timeout=300s

# Destroy all helmfile deployments
destroy: ## Delete all helmfile deployments
	@echo "Destroying all helmfile deployments..."
	docker run --rm --net=host \
		-v $(shell pwd):/wd \
		-w /wd \
		-e KUBECONFIG=/wd/.kube/config \
		$(HELMFILE_IMAGE) \
		helmfile destroy --file conf/helmfile.yaml

# Clean up everything
clean: cluster-down ## Clean up everything
