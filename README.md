# GoHex

A Go microservice example based on Hexagonal Architecture (Ports and Adapters Architecture) that demonstrates how to build a maintainable, testable and high-performance system.

## 🌟 Features

- Hexagonal Architecture based on Domain-Driven Design (DDD)
- CQRS pattern for command and query separation
- Event Sourcing with optimistic concurrency control
- Distributed tracing with OpenTelemetry
- Metrics collection with Prometheus
- Structured logging with Zap
- JWT-based authentication
- MySQL for persistence
- Redis for caching
- Kafka for event bus
- Graceful shutdown
- Comprehensive error handling
- Unit of Work pattern
- Validation using validator/v10
- API documentation with Swagger

## 🚀 Quick Start

### Prerequisites

- Go 1.21+
- Docker & Docker Compose
- MySQL 8.0+
- Redis 6.0+
- Kafka 2.8+

### Installation

```bash
# Clone repository
git clone https://github.com/gohex/gohex

# Enter project directory
cd gohex

# Start infrastructure services
docker-compose up -d

# Install dependencies
go mod download

# Run migrations
make migrate

# Start service
make run
```

## 📁 Project Structure

```
📦 gohex
├── 📂 cmd                            # Application entries
│   └── 📂 api                       # API server
├── 📂 internal                      # Internal packages
│   ├── 📂 domain                    # Domain layer
│   │   ├── 📂 aggregate            # Aggregates
│   │   ├── 📂 event               # Domain events
│   │   ├── 📂 service            # Domain services
│   │   └── 📂 vo                # Value objects
│   ├── 📂 application             # Application layer
│   │   ├── 📂 command           # Command handlers
│   │   ├── 📂 query            # Query handlers
│   │   ├── 📂 dto             # Data transfer objects
│   │   └── 📂 port           # Ports (interfaces)
│   └── 📂 infrastructure        # Infrastructure layer
│       ├── 📂 adapter         # Adapters
│       ├── 📂 bus            # Command/Query/Event buses
│       └── 📂 config        # Configuration
└── 📂 pkg                    # Public packages
    ├── 📂 errors           # Error handling
    ├── 📂 tracer          # Distributed tracing
    └── 📂 validator      # Validation utilities
```

## 🔧 Configuration

Configuration is handled through environment variables or config file:

```yaml
# config/config.yaml
app:
  name: gohex
  env: development
  version: 1.0.0

http:
  port: 8080
  timeout: 30s

database:
  driver: mysql
  host: localhost
  port: 3306
  name: gohex
  user: root
  password: secret

redis:
  host: localhost
  port: 6379

kafka:
  brokers:
    - localhost:9092
  group: gohex

jwt:
  secret: your-secret-key
  duration: 24h

log:
  level: debug
  format: json
```

## 📖 Documentation

- [API Documentation](docs/api.md)
- [Architecture Overview](docs/architecture.md)
- [Development Guide](docs/development.md)
- [Deployment Guide](docs/deployment.md)

## 🤝 Contributing

Please read our [Contributing Guide](CONTRIBUTING.md) before submitting a Pull Request.

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details. 