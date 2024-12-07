# GoHex

A Go microservice example based on Hexagonal Architecture (Ports and Adapters Architecture) that demonstrates how to build a maintainable, testable and high-performance system.

## 🌟 Features

- Hexagonal Architecture based on Domain-Driven Design (DDD)
- CQRS pattern for command and query separation
- Complete event sourcing support
- Distributed tracing and monitoring
- Cache optimization and read/write separation
- Robust error handling
- Transaction management with Unit of Work pattern

## 🚀 Quick Start

### Prerequisites

- Go 1.21+
- Docker & Docker Compose
- Make

### Installation

```bash
# Clone repository
git clone https://github.com/yourusername/gohex

# Enter project directory
cd gohex

# Install dependencies
make install

# Start dependent services
make docker-up

# Run service
make run
```

### Basic Usage

```bash
# Run tests
make test

# Run lint
make lint

# Build
make build
```

## 📁 Project Structure

```
📦 gohex
├── 📂 cmd                  # Application entries
├── 📂 internal            # Internal code
│   ├── 📂 domain         # Domain layer
│   ├── 📂 application    # Application layer
│   └── 📂 infrastructure # Infrastructure layer
└── 📂 pkg                # Public packages
```

## 🔧 Configuration

Configure through environment variables or configuration files:

```yaml
app:
  name: gohex
  version: 1.0.0

http:
  port: 8080
  timeout: 30s

database:
  driver: postgres
  dsn: postgres://user:pass@localhost:5432/dbname

cache:
  driver: redis
  address: localhost:6379
```

## 📖 Documentation

For detailed documentation, please refer to:

- [Architecture Design](docs/architecture.md)
- [API Documentation](docs/api.md)
- [Contributing Guide](CONTRIBUTING.md)

## 🤝 Contributing

Contributions are welcome! Please read our [Contributing Guide](CONTRIBUTING.md) for details.

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details. 