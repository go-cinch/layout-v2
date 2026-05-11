# Go-Cinch Layout V2

A production-ready Go microservice template generator based on [Kratos](https://github.com/go-kratos/kratos) framework with Wire dependency injection.

## Features

- 🚀 **Quick Start**: Generate production-ready microservices in seconds with presets
- 🔧 **Flexible Configuration**: Multiple presets and customizable options
- 🎯 **Clean Architecture**: Clear separation of layers (cmd/internal/api)
- 💉 **Dependency Injection**: Wire-based automatic code generation
- 🗄️ **Database Support**: PostgreSQL/MySQL with GORM ORM
- 🔑 **Distributed ID**: Sonyflake integration for unique ID generation
- 📊 **Observability**: OpenTelemetry tracing support
- 🔄 **Auto-generation**: GORM models from database schema
- 🎭 **Multiple Templates**: Full CRUD or Simple GET operations
- 🚦 **Production Features**: Health checks, middleware, caching, task scheduling

## Quick Start

### 1. Install scaffold

```bash
go install github.com/hay-kot/scaffold@v0.12.0
```

### 2. Create New Service

#### Using Presets (Recommended)

Presets provide pre-configured settings for common use cases:

##### Default Preset - Full Features

Full-featured microservice with Redis, cache, idempotent, task scheduler, and tracing:

```bash
scaffold new https://github.com/go-cinch/layout-v2 \
  --output-dir=. \
  --run-hooks=always \
  --no-prompt \
  --preset default \
  Project=myservice
```

**Features:**
- ✅ Full CRUD operations (Create/Get/Find/Update/Delete)
- ✅ GORM ORM with auto-generated models from database
- ✅ Sonyflake distributed ID generator
- ✅ Redis connection support
- ✅ Multi-layer cache system
- ✅ Idempotent middleware (prevent duplicate requests)
- ✅ Task/Cron worker scheduler
- ✅ OpenTelemetry tracing
- ✅ Transaction support with rollback
- ✅ Health check endpoints (HTTP/gRPC)
- ✅ Header middleware
- 📦 Binary size: ~32MB

**Generated Structure:**
```
myservice/
├── api/              # Protobuf definitions
├── cmd/              # Application entry points
├── configs/          # Configuration files
├── internal/
│   ├── biz/         # Business logic layer
│   ├── data/        # Data access layer (with Sonyflake)
│   ├── server/      # HTTP/gRPC servers
│   └── service/     # Service layer
└── Makefile         # Build automation
```

##### Simple Preset - Minimal

Lightweight microservice without Redis features:

```bash
scaffold new https://github.com/go-cinch/layout-v2 \
  --output-dir=. \
  --run-hooks=always \
  --no-prompt \
  --preset simple \
  Project=myservice
```

**Features:**
- ✅ Simple Get operation only
- ✅ GORM ORM with auto-generated models
- ✅ Sonyflake distributed ID generator
- ✅ OpenTelemetry tracing
- ✅ Transaction support
- ✅ Health check endpoints
- ❌ No Redis/Cache/Task features
- 📦 Binary size: ~25MB

##### NoRedis Preset - Full CRUD without Redis

Full CRUD operations without any Redis dependencies:

```bash
scaffold new https://github.com/go-cinch/layout-v2 \
  --output-dir=. \
  --run-hooks=always \
  --no-prompt \
  --preset noredis \
  Project=myservice
```

**Features:**
- ✅ Full CRUD operations
- ✅ GORM ORM with auto-generated models
- ✅ Sonyflake distributed ID generator
- ✅ OpenTelemetry tracing
- ✅ Transaction support
- ✅ Health check endpoints
- ❌ No Redis connection
- ❌ No Cache layer
- ❌ No Idempotent middleware
- ❌ No Task/Cron worker
- 📦 Binary size: ~27MB

### 3. Presets Comparison

| Feature | Default | Simple | NoRedis |
|---------|---------|--------|---------|
| **Proto Template** | Full CRUD | Simple GET | Full CRUD |
| **CRUD Operations** | ✅ C/R/U/D | ✅ R only | ✅ C/R/U/D |
| **Sonyflake ID** | ✅ | ✅ | ✅ |
| **Database (GORM)** | ✅ PostgreSQL | ✅ PostgreSQL | ✅ PostgreSQL |
| **Auto-gen Models** | ✅ | ✅ | ✅ |
| **Transaction** | ✅ | ✅ | ✅ |
| **Redis** | ✅ | ❌ | ❌ |
| **Cache Layer** | ✅ | ❌ | ❌ |
| **Idempotent** | ✅ | ❌ | ❌ |
| **Task/Cron** | ✅ | ❌ | ❌ |
| **OpenTelemetry** | ✅ | ✅ | ✅ |
| **Health Check** | ✅ | ✅ | ✅ |
| **Binary Size** | ~32MB | ~25MB | ~27MB |
| **Use Case** | Production | Learning | Stateless API |

### 4. Customization

#### Override Preset Options

You can override specific preset options:

```bash
# Enable WebSocket on default preset
scaffold new https://github.com/go-cinch/layout-v2 \
  --preset default \
  Project=myservice \
  enable_ws=true

# Use MySQL instead of PostgreSQL
scaffold new https://github.com/go-cinch/layout-v2 \
  --preset default \
  Project=myservice \
  db_type=mysql

# Change HTTP/gRPC ports
scaffold new https://github.com/go-cinch/layout-v2 \
  --preset simple \
  Project=myservice \
  http_port=9090 \
  grpc_port=9190

# Customize service name
scaffold new https://github.com/go-cinch/layout-v2 \
  --preset default \
  Project=game \
  service_name=user
```

#### Available Configuration Options

| Option | Values | Default | Description |
|--------|--------|---------|-------------|
| `service_name` | string | Project name | Service name (used in logs, metrics) |
| `module_name` | string | service_name | Go module name |
| `http_port` | string | `8080` | HTTP server port |
| `grpc_port` | string | `8180` | gRPC server port |
| `openapi_dev_url` | string | `https://example.com/api/{service}` | Development server URL in generated OpenAPI docs |
| `proto_template` | `full`/`simple` | varies | API template (full CRUD or simple GET) |
| `db_type` | `postgres`/`mysql` | `postgres` | Database type |
| `orm_type` | `gorm`/`none` | `gorm` | ORM framework |
| `enable_ws` | `true`/`false` | `false` | Enable WebSocket support |
| `enable_redis` | `true`/`false` | varies | Enable Redis connection |
| `enable_cache` | `true`/`false` | varies | Enable cache layer |
| `enable_idempotent` | `true`/`false` | varies | Enable idempotent middleware |
| `enable_task` | `true`/`false` | varies | Enable task/cron scheduler |
| `enable_trace` | `true`/`false` | `true` | Enable OpenTelemetry tracing |
| `enable_i18n` | `true`/`false` | `false` | Enable i18n support |

### 5. Interactive Mode

Answer prompts to configure all options:

```bash
scaffold new https://github.com/go-cinch/layout-v2 \
  --output-dir=. \
  --run-hooks=always \
  Project=myservice
```

## Building and Running

### 1. Generate Code

```bash
cd myservice
make all  # Install tools, generate proto/wire/config, lint
```

### 2. Database Setup

**Start PostgreSQL (Docker):**
```bash
docker run -d --name postgres \
  -e POSTGRES_USER=root \
  -e POSTGRES_PASSWORD=password \
  -p 5432:5432 \
  postgres:17
```

**Configure in `configs/db.yaml`:**
```yaml
db:
  driver: postgres
  dsn: "host=localhost user=root password=password dbname=myservice port=5432 sslmode=disable TimeZone=UTC"
  migrate: true  # Auto-run migrations
```

### 3. Redis Setup (for default preset)

```bash
docker run -d --name redis \
  -p 6379:6379 \
  redis:7
```

**Configure in `configs/redis.yaml`:**
```yaml
redis:
  dsn: "redis://:password@localhost:6379/0"
```

### 4. Build

```bash
make build  # Output: ./bin/myservice
```

### 5. Run

```bash
./bin/myservice -c ./configs
```

**Endpoints:**
- HTTP: http://localhost:8080
- gRPC: localhost:8180
- Health: http://localhost:8080/health
- API Docs: http://localhost:8080/docs/

## Development Workflow

```bash
# Generate API from proto
make api

# Generate Wire dependency injection
make wire

# Run linter
make lint

# Run tests
make test

# Complete build pipeline
make all

# Clean generated files
make clean
```

## Project Naming Guidelines

### ⚠️ Important Naming Rules

**Avoid project names ending with 's'** (e.g., `users`, `items`, `orders`)

**Why?**
GORM's gentool singularizes table names incorrectly for names ending in 's':
- `users` → generates `User` ❌ (expected `Users`)
- `orders` → generates `Order` ❌ (expected `Orders`)

**Recommended Naming:**
- ✅ `user`, `order`, `pay`
- ❌ `users`, `items`, `orders` (will cause compilation errors)

**If you must use plural names:**
Manually correct the generated model struct names in `internal/data/model/*.gen.go` after generation.

## Contributing

Issues and pull requests are welcome!

## License

Apache 2.0
