# 📊 Learning Prometheus and Grafana - Observability Stack

> A comprehensive hands-on learning project for mastering observability with Prometheus, Grafana, Loki, Jaeger, and OpenTelemetry in Go applications.

[![Go](https://img.shields.io/badge/Go-1.21+-00ADD8?style=flat&logo=go)](https://golang.org/)
[![Docker](https://img.shields.io/badge/Docker-Compose-2496ED?style=flat&logo=docker)](https://www.docker.com/)
[![Prometheus](https://img.shields.io/badge/Prometheus-Monitoring-E6522C?style=flat&logo=prometheus)](https://prometheus.io/)
[![Grafana](https://img.shields.io/badge/Grafana-Visualization-F46800?style=flat&logo=grafana)](https://grafana.com/)

## 📚 Table of Contents

- [Overview](#overview)
- [Learning Objectives](#learning-objectives)
- [Project Structure](#project-structure)
- [Technologies Used](#technologies-used)
- [Getting Started](#getting-started)
- [Projects](#projects)
  - [1. Logging Challenge](#1-logging-challenge)
  - [2. Demo App - Full Observability Stack](#2-demo-app---full-observability-stack)
- [Observability Pillars](#observability-pillars)
- [Architecture](#architecture)
- [Key Concepts Learned](#key-concepts-learned)
- [Resources](#resources)

---

## 🎯 Overview

This repository contains my learning journey in implementing a complete observability stack for modern microservices. The project demonstrates practical implementations of the **Three Pillars of Observability**:

1. **📝 Logging** - Structured logging with Zerolog, Fluent Bit, and Loki
2. **📈 Metrics** - Application metrics with Prometheus and OpenTelemetry
3. **🔍 Tracing** - Distributed tracing with Jaeger and OpenTelemetry

The repository includes two main projects:
- **Logging Challenge**: A focused implementation of logging, metrics, and tracing
- **Demo App**: A production-ready gRPC/REST API with complete observability stack

---

## 🎓 Learning Objectives

By working through this repository, I learned:

- ✅ Implementing structured logging with **Zerolog**
- ✅ Setting up **Prometheus** for metrics collection
- ✅ Creating custom metrics (Counter, Gauge, Histogram)
- ✅ Integrating **OpenTelemetry** for traces and metrics
- ✅ Using **Jaeger** for distributed tracing visualization
- ✅ Configuring **Grafana** for unified observability dashboards
- ✅ Setting up **Loki** for log aggregation
- ✅ Using **Fluent Bit** for log collection and forwarding
- ✅ Implementing context propagation with W3C Trace Context
- ✅ Building production-ready observability pipelines

---

## 📁 Project Structure

```
Learn_Prometheus_and_Grafana/
├── README.md                          # This file
├── 4b6f5b80-b2a1-11eb-90c9-2e45f036ff26.png  # Architecture diagram
│
├── challange/                         # Learning challenges
│   └── logging-challange/             # Logging & observability challenge
│       ├── main.go                    # Go application with full instrumentation
│       ├── docker-compose.yml         # Observability stack (Prometheus, Grafana, Jaeger)
│       ├── go.mod                     # Go dependencies
│       ├── logs/                      # Application logs
│       │   └── app.log
│       └── scripts/                   # Configuration files
│           ├── fluentbit/             # Fluent Bit config
│           ├── grafana/               # Grafana provisioning
│           ├── loki/                  # Loki config
│           └── prometheus/            # Prometheus config
│
└── demo-app/                          # Production-ready demo application
    ├── README.md                      # Demo app documentation
    ├── docker-compose.yml             # Full stack with dependencies
    ├── Makefile                       # Build and run commands
    ├── cmd/                           # Application entry points
    ├── course/                        # Business logic (catalog, booking)
    ├── internal/                      # Internal packages
    │   ├── instrumentation/           # Telemetry setup
    │   ├── metrics/                   # Custom metrics
    │   └── grpc/                      # gRPC instrumentation
    ├── pkg/                           # Generated protobuf code
    ├── scripts/                       # Configuration files
    │   ├── fluentbit/                 # Log collection
    │   ├── grafana/                   # Dashboards & datasources
    │   ├── loki/                      # Log aggregation
    │   ├── locust/                    # Load testing
    │   └── prometheus/                # Metrics collection
    └── third_party/                   # Swagger UI
```

---

## 🛠️ Technologies Used

### Core Application
- **Go 1.21+** - Primary programming language
- **Gin** - HTTP web framework
- **gRPC** - High-performance RPC framework
- **gRPC-Gateway** - REST API proxy for gRPC

### Observability Stack
- **Zerolog** - High-performance structured logging
- **Prometheus** - Metrics collection and storage
- **Grafana** - Visualization and dashboards
- **Loki** - Log aggregation system
- **Jaeger** - Distributed tracing platform
- **OpenTelemetry** - Unified observability framework
- **Fluent Bit** - Log processor and forwarder

### Infrastructure
- **Docker & Docker Compose** - Containerization
- **PostgreSQL** - Database
- **Redis** - Caching
- **Locust** - Load testing

---

## 🚀 Getting Started

### Prerequisites

```bash
# Required
- Go 1.21 or higher
- Docker & Docker Compose
- Make

# Optional (for load testing)
- Python 3.x
- pip & virtualenv
```

### Quick Start

1. **Clone the repository**
   ```bash
   git clone <repository-url>
   cd Learn_Prometheus_and_Grafana
   ```

2. **Choose a project to run**
   - [Logging Challenge](#1-logging-challenge) - Start here for basics
   - [Demo App](#2-demo-app---full-observability-stack) - Production-ready example

---

## 📦 Projects

### 1. Logging Challenge

**Location**: `challange/logging-challange/`

A focused learning project demonstrating the evolution of observability implementation from basic logging to full distributed tracing.

#### Features Implemented

##### 🔹 Phase 1: Basic Logging
- Structured logging with Zerolog
- Log levels (Info, Debug, Warn, Error, Fatal)
- File and console output

##### 🔹 Phase 2: HTTP Request Logging
- Request ID middleware for correlation
- Request/response logging
- Latency tracking
- Context-aware logging

##### 🔹 Phase 3: Metrics with Prometheus
- Custom metrics (Counter, Histogram, Gauge)
- HTTP request metrics
- Request duration histograms
- Active requests tracking
- `/metrics` endpoint for Prometheus scraping

##### 🔹 Phase 4: Distributed Tracing with OpenTelemetry
- Jaeger integration
- Span creation and management
- Trace context propagation
- W3C Trace Context standard
- `otelhttp` for automatic HTTP instrumentation

##### 🔹 Phase 5: Unified Observability
- Combined logs, metrics, and traces
- Prometheus exporter for metrics
- gRPC exporter for traces
- Context propagation across services
- Rich span attributes and events

#### Quick Start

```bash
cd challange/logging-challange

# Start observability stack
docker compose up -d

# Run the application
go run main.go
```

#### Available Endpoints

```bash
# Application endpoints
GET  http://localhost:8080/hello         # Simple endpoint
GET  http://localhost:8080/slow          # Simulates slow operations
GET  http://localhost:8080/error         # Error handling demo
GET  http://localhost:8080/external      # External API calls demo
GET  http://localhost:8080/propagation   # Context propagation demo
GET  http://localhost:8080/metrics       # Prometheus metrics

# Observability dashboards
http://localhost:9090     # Prometheus UI
http://localhost:16686    # Jaeger UI (Traces)
http://localhost:3000     # Grafana (Dashboards)
```

#### Testing the Implementation

```bash
# Generate traffic
curl http://localhost:8080/hello
curl http://localhost:8080/slow
curl http://localhost:8080/error
curl http://localhost:8080/external

# View metrics
curl http://localhost:8080/metrics

# Check logs
tail -f logs/app.log
```

### 2. Demo App - Full Observability Stack

**Location**: `demo-app/`

A production-ready gRPC/REST API application demonstrating enterprise-grade observability implementation for a course booking system.

#### Application Features

- **Course Catalog Management** - CRUD operations for courses
- **Booking System** - Course reservation and management
- **gRPC Server** - High-performance RPC endpoints
- **REST API** - HTTP endpoints via gRPC-Gateway
- **Swagger Documentation** - Interactive API documentation

#### Observability Features

##### 📊 Metrics
- HTTP request metrics (count, duration, status)
- gRPC metrics (calls, latency, errors)
- Database query metrics
- Redis cache metrics
- Custom business metrics

##### 📝 Logging
- Structured JSON logs with Zerolog
- Request correlation with unique IDs
- Log aggregation with Loki
- Log collection with Fluent Bit
- Grafana log exploration

##### 🔍 Tracing
- Distributed tracing with Jaeger
- gRPC interceptors for automatic tracing
- Database query tracing
- External API call tracing
- Service dependency visualization

#### Architecture

```
┌─────────────────────────────────────────┐
│         HTTP/REST Clients               │
│         (Browser, Postman, etc)         │
└──────────────┬──────────────────────────┘
               │ HTTP/REST
               ▼
┌─────────────────────────────────────────┐
│      gRPC-Gateway (HTTP Server)         │
│      Port: HTTP (e.g., 8080)            │
└──────────────┬──────────────────────────┘
               │ gRPC (internal)
               ▼
┌─────────────────────────────────────────┐
│         gRPC Server                     │
│      Port: gRPC (e.g., 9090)            │
│  ┌─────────────────────────────────┐    │
│  │  Booking Server / Catalog Server│    │
│  └──────────┬──────────────────────┘    │
└─────────────┼────────────────────────── ┘
              │
              ▼
┌─────────────────────────────────────────┐
│      Business Logic (Service Layer)     │
│  ┌────────────────┬─────────────────┐   │
│  │ BookingService │ CatalogService  │   │
│  └────────┬───────┴────────┬────────┘   │
└───────────┼────────────────┼────────── ─┘
            │                │
            ▼                ▼
┌─────────────────────────────────────────┐
│    Data Access Layer (Store/Repository) │
│  ┌────────────┬──────────────────────┐  │
│  │BookingStore│   CatalogStore       │  │
│  └─────┬──────┴──────┬───────────────┘  │
└────────┼─────────────┼──────────────────┘
         │             │
    ┌────▼──── ┐   ┌───▼────┐
    │PostgreSQL│   │ Redis  │
    └───────── ┘   └────────┘
```

#### Quick Start

```bash
cd demo-app

# Set environment variables
export REDIS_PASSWORD=your_password_here

# Start all dependencies
make bootstrap
docker compose up -d

# Run the application
make course/server

# Seed the database (in another terminal)
make course/seed
```

#### Available Services

```bash
# Application
http://localhost:8800          # REST API
http://localhost:8800/swagger  # Swagger UI
grpc://localhost:9900          # gRPC Server

# Observability Stack
http://localhost:9090          # Prometheus
http://localhost:16686         # Jaeger
http://localhost:3000          # Grafana
http://localhost:3100          # Loki

# Infrastructure
localhost:5432                 # PostgreSQL
localhost:6379                 # Redis
```

#### Load Testing

The demo app includes a Locust-based load generator for testing:

```bash
cd scripts/locust

# Setup virtual environment
virtualenv .venv
source .venv/bin/activate

# Install dependencies
pip install -r requirements.txt

# Run load generator
locust -f locustfiles --class-picker --modern-ui -H http://localhost:8800

# Open browser
http://localhost:8089
```

**Load Testing Scenarios:**
- **GeneralUser** - General API testing (use 1-5 users)
- **CompetingUser** - Concurrent booking simulation (use 10+ users)

#### API Examples

```bash
# List courses
curl http://localhost:8800/v1/courses

# Get course details
curl http://localhost:8800/v1/courses/{course_id}

# Create booking
curl -X POST http://localhost:8800/v1/bookings \
  -H "Content-Type: application/json" \
  -d '{"course_id": "...", "user_id": "..."}'

# List bookings
curl http://localhost:8800/v1/bookings
```

---

## 🔍 Observability Pillars

### 1. Logs - What happened?

**Tools**: Zerolog → Fluent Bit → Loki → Grafana

- Structured JSON logging
- Request correlation with unique IDs
- Log levels and filtering
- Centralized log aggregation
- Log-based alerting

**Example Log Entry:**
```json
{
  "level": "info",
  "request_id": "550e8400-e29b-41d4-a716-446655440000",
  "method": "GET",
  "path": "/v1/courses",
  "status": 200,
  "latency": 45,
  "time": "2024-01-25T10:30:00Z",
  "message": "request completed"
}
```

### 2. Metrics - How much/many?

**Tools**: OpenTelemetry → Prometheus → Grafana

- Request rate, error rate, duration (RED metrics)
- Resource utilization (CPU, memory)
- Custom business metrics
- Histograms and percentiles
- Alerting based on thresholds

**Example Metrics:**
```
http_server_request_count{method="GET",route="/v1/courses",status="200"} 1234
http_server_request_duration_seconds_bucket{le="0.1"} 980
http_server_requests_active 5
```

### 3. Traces - Where and why?

**Tools**: OpenTelemetry → Jaeger

- End-to-end request flow
- Service dependency mapping
- Performance bottleneck identification
- Error propagation tracking
- Context propagation across services

**Trace Example:**
```
Trace ID: 550e8400-e29b-41d4-a716-446655440000
├─ GET /v1/courses (120ms)
   ├─ database_query (80ms)
   │  └─ SELECT * FROM courses
   ├─ cache_check (10ms)
   └─ response_serialization (30ms)
```

---

## 🏗️ Architecture

### Observability Data Flow

```
┌─────────────┐
│ Application │
└──────┬──────┘
       │
       ├─────► Logs ────────► Fluent Bit ────► Loki ────► Grafana
       │
       ├─────► Metrics ─────► Prometheus ───────────────► Grafana
       │
       └─────► Traces ──────► Jaeger ──────────────────► Jaeger UI
                                                          └─► Grafana
```

### Technology Stack Diagram

```
┌──────────────────────────────────────────────────────┐
│                   Visualization Layer                 │
│              Grafana (Unified Dashboard)              │
└──────────────────────────────────────────────────────┘
                          ▲
                          │
        ┌─────────────────┼─────────────────┐
        │                 │                 │
┌───────▼──────┐  ┌──────▼──────┐  ┌──────▼──────┐
│     Loki     │  │ Prometheus  │  │   Jaeger    │
│ (Log Store)  │  │(Metric Store│  │(Trace Store)│
└───────▲──────┘  └──────▲──────┘  └──────▲──────┘
        │                │                 │
┌───────┴──────┐  ┌──────┴──────┐  ┌──────┴──────┐
│  Fluent Bit  │  │   OTel SDK  │  │   OTel SDK  │
│(Log Collect) │  │  (Metrics)  │  │  (Traces)   │
└───────▲──────┘  └──────▲──────┘  └──────▲──────┘
        │                │                 │
        └────────────────┴─────────────────┘
                         │
              ┌──────────▼──────────┐
              │   Go Application    │
              │   (Instrumented)    │
              └─────────────────────┘
```

---

## 💡 Key Concepts Learned

### 1. Structured Logging
- Using Zerolog for high-performance logging
- Adding context to logs (request ID, user ID, etc.)
- Log levels and when to use them
- Writing logs to multiple outputs

### 2. Metrics Collection
- Understanding metric types:
  - **Counter**: Monotonically increasing values (request count)
  - **Gauge**: Values that can go up/down (active connections)
  - **Histogram**: Distribution of values (request duration)
- Creating custom metrics
- Metric labels and cardinality
- Prometheus scraping and storage

### 3. Distributed Tracing
- Span creation and lifecycle
- Parent-child span relationships
- Trace context propagation
- W3C Trace Context standard
- Span attributes, events, and errors
- Using `otelhttp` for automatic instrumentation

### 4. Context Propagation
- Passing context across function calls
- HTTP header injection (traceparent, tracestate)
- Correlating logs, metrics, and traces
- Request ID propagation

### 5. Observability Best Practices
- Sampling strategies for high-traffic systems
- Metric naming conventions
- Log aggregation patterns
- Dashboard design principles
- Alert configuration

### 6. Performance Considerations
- Minimizing observability overhead
- Async metric/trace export
- Batching and buffering
- Resource limits and backpressure


## 🎯 Next Steps

After completing this learning journey, you can:

1. **Extend the Demo App**
   - Add more business metrics
   - Create custom Grafana dashboards
   - Implement alerting rules
   - Add more trace instrumentation

2. **Production Readiness**
   - Set up metric retention policies
   - Configure log rotation
   - Implement sampling strategies
   - Set up monitoring alerts

3. **Advanced Topics**
   - Service mesh integration (Istio, Linkerd)
   - Distributed tracing across multiple services
   - Custom OpenTelemetry exporters
   - Advanced Grafana dashboard design

4. **Apply to Your Projects**
   - Integrate observability into your microservices
   - Set up production monitoring
   - Create SLO/SLI dashboards
   - Implement on-call alerting

---

## 📝 License

This project is for educational purposes. Feel free to use and modify for your learning.

---

## 🙏 Acknowledgments

- Thanks to the open-source community for the amazing tools
- Inspired by production observability practices
- Built while learning modern observability patterns

---

## 📬 Contact

If you have questions or suggestions about this learning project, feel free to reach out!

---

**Happy Learning! 🚀**
