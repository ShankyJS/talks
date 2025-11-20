# OTEL, Jaeger, and Go Services

This project sets up a local development environment with OpenTelemetry, Jaeger (all-in-one with Cassandra storage), and Golang services (Frontend & Backend) using Kind and Helmfile.

## Prerequisites

- Docker
- `kind` (installed on host, or ensure `make cluster-up` can run it)
- `make`

## Usage

The entire workflow is managed via the `Makefile`.

### 1. Create Cluster & Deploy

```bash
make all
```

This command will:
1.  Create a Kind cluster (if not exists).
2.  Deploy Jaeger (all-in-one with Cassandra), OTEL Collector, Backend, and Frontend services using Helmfile (running in Docker).

**Note:** Cassandra takes a couple of minutes to initialize. Jaeger components may restart until Cassandra is ready.

### 2. Access Services

You can use the following make targets to easily access the services:

- **Jaeger UI**:
  ```bash
  make port-forward-jaeger
  ```
  Open [http://localhost:16686](http://localhost:16686).

- **Frontend App**:
  ```bash
  make port-forward-app
  ```
  Open [http://localhost:8080](http://localhost:8080).

### 3. Cassandra Commands

- **Check Cassandra Data**:
  ```bash
  make cassandra-check-data
  ```
  Shows keyspaces and a trace count from the Jaeger keyspace.

- **Connect to Cassandra Shell**:
  ```bash
  make cassandra-shell
  ```
  Opens an interactive CQL shell inside the Jaeger Cassandra pod.

- **Wait for Cassandra**:
  ```bash
  make cassandra-wait
  ```
  Waits until the Cassandra pod is ready before you generate traffic.

### 4. Kubeconfig

By default, this project uses a local `.kube/config` file to keep your host environment clean and ensure the Dockerized tools work correctly.

If you want to use `kubectl` from your host without specifying the config file, run:

```bash
make kubeconfig-export
```

This will merge the cluster config into your default `~/.kube/config`.

### 5. Cleanup

- **Delete all deployments** (keeps cluster running):
  ```bash
  make destroy
  ```

- **Delete everything** (cluster + deployments):
  ```bash
  make clean
  ```

## Architecture

- **Jaeger**: Distributed tracing backend using the all-in-one image with Cassandra storage
  - **Collector**: Receives traces from OTEL Collector (inside the all-in-one pod)
  - **Query**: Provides UI and API for viewing traces
  - **Cassandra**: Stateful backend where traces and service dependencies are stored
- **OTEL Collector**: Receives traces from apps and forwards to Jaeger
- **Backend**: Golang service (simulated with Nginx for now)
- **Frontend**: Golang service (simulated with Nginx for now)

## Configuration

- `conf/helmfile.yaml`: Main deployment descriptor
- `conf/values/`: Value files for charts
- `charts/generic-service/`: Generic Helm chart using helmet library

## Next Steps

To see traces in Jaeger:
1. Replace the Nginx placeholder apps with your actual Go applications
2. Instrument them with OpenTelemetry SDK
3. Configure them to send traces to `otel-collector-opentelemetry-collector.monitoring.svc.cluster.local:4317`
4. Traces will flow: **App → OTEL Collector → Jaeger (all-in-one with Cassandra)**
5. View traces in Jaeger UI and see the Dependencies graph (System Architecture) populate from Cassandra data!
