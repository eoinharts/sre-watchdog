# SRE Watchdog

A lightweight synthetic monitoring service written in Go.

SRE Watchdog periodically checks a configured HTTP endpoint, records availability and request latency, and exposes Prometheus metrics through a `/metrics` endpoint.

The project is being built as a practical SRE exercise covering application observability, containerisation, CI/CD and cloud deployment.

## Current Features

- Periodic HTTP endpoint checks
- Configurable target URL
- Configurable check interval
- HTTP request timeout
- Success and failure detection
- Prometheus metrics endpoint
- Request duration histogram
- Docker multi-stage build
- Non-root runtime container

## Prometheus Metrics

The service currently exposes:

| Metric | Type | Description |
|---|---|---|
| `watchdog_up` | Gauge | `1` when the target is healthy, `0` when unhealthy |
| `watchdog_checks_total` | Counter | Total number of endpoint checks |
| `watchdog_failures_total` | Counter | Total number of failed checks |
| `watchdog_request_duration_seconds` | Histogram | Distribution of endpoint check durations |

Metrics are exposed at:

```text
http://localhost:8080/metrics
```

## Project Structure

```text
sre-watchdog/
├── cmd/
│   └── watchdog/
│       └── main.go
├── .dockerignore
├── .gitignore
├── Dockerfile
├── go.mod
├── go.sum
└── README.md
```

## Run Locally

### Prerequisites

- Go
- Git

Install dependencies:

```bash
go mod tidy
```

Run with the default configuration:

```bash
go run ./cmd/watchdog
```

The default target is:

```text
https://example.com
```

The default check interval is:

```text
30 seconds
```

## Configuration

The application is configured using environment variables.

### `TARGET_URL`

The endpoint to monitor.

Example:

```bash
TARGET_URL=https://example.com go run ./cmd/watchdog
```

### `CHECK_INTERVAL`

The number of seconds between checks.

Example:

```bash
CHECK_INTERVAL=5 go run ./cmd/watchdog
```

Both can be supplied together:

```bash
TARGET_URL=https://example.com \
CHECK_INTERVAL=5 \
go run ./cmd/watchdog
```

Invalid, zero or negative check intervals fall back to the default interval of 30 seconds.

## View Prometheus Metrics

While the application is running:

```bash
curl -s http://localhost:8080/metrics | grep watchdog
```

Example output:

```text
watchdog_checks_total 8
watchdog_failures_total 0
watchdog_up 1
watchdog_request_duration_seconds_count 8
```

## Build with Docker

Build the image:

```bash
docker build -t sre-watchdog:local .
```

Run the container:

```bash
docker run --rm \
  --name sre-watchdog \
  -p 8080:8080 \
  -e CHECK_INTERVAL=5 \
  -e TARGET_URL=https://example.com \
  sre-watchdog:local
```

Verify metrics from another terminal:

```bash
curl -s http://localhost:8080/metrics | grep watchdog
```

## Failure Testing

A failed target can be tested by pointing the watchdog at an unavailable local port:

```bash
TARGET_URL=http://localhost:9999 \
CHECK_INTERVAL=5 \
go run ./cmd/watchdog
```

Expected behaviour:

- `watchdog_up` becomes `0`
- `watchdog_failures_total` increments
- the connection error is logged

## Architecture

```text
Target HTTP Endpoint
        |
        | periodic checks
        v
   SRE Watchdog
      /      \
     /        \
    v          v
Application   /metrics
Logs             |
                 v
             Prometheus
```

## Planned Work

Future iterations may include:

- Unit tests
- Jenkins CI pipeline
- Automated Docker image builds
- Amazon ECR image storage
- Kubernetes deployment
- Amazon EKS
- CloudWatch alerting
- SNS notifications
- Lambda-based alert routing
- Failure injection and rollback testing

## Purpose

This project is intended to explore practical Site Reliability Engineering concepts including:

- synthetic monitoring
- observability
- metrics
- latency measurement
- failure detection
- runtime configuration
- containerisation
- CI/CD
- deployment automation
- alerting