# Go Services with OpenTelemetry

This directory contains two Go microservices instrumented with OpenTelemetry:

## Backend Service

**Location:** `apps/backend/`

A simple HTTP service that returns random inspirational words.

**Endpoints:**
- `GET /words` - Returns 5 random words with timestamp
- `GET /health` - Health check endpoint

**Environment Variables:**
- `OTEL_EXPORTER_OTLP_ENDPOINT` - OTEL Collector endpoint (default: localhost:4317)
- `PORT` - Server port (default: 8080)

## Frontend Service

**Location:** `apps/frontend/`

A web UI that fetches and displays random words from the backend.

**Endpoints:**
- `GET /` - Web UI with button to generate words
- `GET /api/words` - Proxies request to backend
- `GET /health` - Health check endpoint

**Environment Variables:**
- `OTEL_EXPORTER_OTLP_ENDPOINT` - OTEL Collector endpoint (default: localhost:4317)
- `BACKEND_URL` - Backend service URL (default: http://localhost:8080)
- `PORT` - Server port (default: 3000)

## OpenTelemetry Instrumentation

Both services are instrumented with:
- **Trace propagation** - Context is propagated from frontend → backend
- **Custom spans** - Each operation creates detailed spans
- **Attributes** - Rich metadata attached to spans
- **Error recording** - Errors are captured in traces

## Building and Running

### Local Development

```bash
# Backend
cd apps/backend
go mod download
go run main.go

# Frontend
cd apps/frontend
go mod download
go run main.go
```

### Docker Build

```bash
# Build images
make build-images

# Build and load into Kind
make load-images

# Build, load, and deploy
make deploy-apps
```

## Trace Flow

```
User clicks button
  ↓
Frontend: homeHandler span
  ↓
Frontend: apiWordsHandler span
  ↓
Frontend: fetchWordsFromBackend span (HTTP call with trace context)
  ↓
Backend: wordsHandler span (receives trace context)
  ↓
Backend: getRandomWords span
  ↓
All spans sent to OTEL Collector → Jaeger → Cassandra
```

## Viewing Traces

1. Deploy the services: `make deploy-apps`
2. Port-forward frontend: `make port-forward-app`
3. Open http://localhost:8080 and click "Generate Words"
4. Port-forward Jaeger: `make port-forward-jaeger`
5. Open http://localhost:16686 to see traces in Jaeger UI
6. Check the Dependencies graph to see service relationships!
