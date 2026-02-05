# Go Kasir Umam DS

A simple yet robust cashier application written in Go, providing a REST API for managing products and categories with clean architecture.

[![Go Version](https://img.shields.io/badge/Go-1.22+-00ADD8?style=flat&logo=go)](https://go.dev/)
[![Test Coverage](https://img.shields.io/badge/coverage-87.3%25-brightgreen)](docs/coverage/CI_CD_COVERAGE.md)
[![License](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)

## 🚀 Quick Start

```bash
# Clone the repository
git clone git@github.com:fajar7xx/go-kasir-umam-ds.git
cd go-kasir-umam-ds

# Run the application
go run main.go

# Server runs on http://localhost:8080
```

## 📚 Documentation

**Complete documentation is available in the [`docs/`](docs/README.md) directory.**

### Quick Links

- 📖 [**Full Documentation Index**](docs/README.md) - Start here!
- 🏗️ [Architecture](docs/architecture/) - System design and patterns
- 🔌 [API Reference](docs/api/) - API endpoints and examples
- 🛠️ [Development Guides](docs/guides/) - How-to guides
- 🧪 [Test Coverage](docs/coverage/CI_CD_COVERAGE.md) - CI/CD & coverage config
- 📋 [Implementation Plans](docs/plans/) - Feature plans and RFCs

## 🏗️ Project Structure

```
.
├── main.go                    # Application entry point
├── handlers/                  # HTTP handlers (presentation layer)
├── internal/
│   ├── services/             # Business logic layer
│   ├── repositories/         # Data access layer
│   ├── database/             # Database configuration
│   └── mocks/                # Test mocks
├── models/                    # Data structures
├── utils/                     # Utility functions
├── config/                    # Configuration
├── Makefile                   # Build and test commands
├── .testcoverage.yml         # Coverage thresholds

```

### Architecture

The project follows **Clean Architecture** principles:

```
┌─────────────┐
│   Handler   │  ← HTTP Layer (handlers/)
└──────┬──────┘
       │
┌──────▼──────┐
│   Service   │  ← Business Logic (internal/services/)
└──────┬──────┘
       │
┌──────▼──────┐
│ Repository  │  ← Data Access (internal/repositories/)
└──────┬──────┘
       │
┌──────▼──────┐
│  Database   │  ← PostgreSQL
└─────────────┘
```

## 🔌 API Endpoints

### Health Check
- `GET /health` - Check application health

### Categories
- `GET /api/v1/categories` - Get all categories
- `GET /api/v1/categories/{id}` - Get category by ID
- `POST /api/v1/categories` - Create new category
- `PUT /api/v1/categories/{id}` - Update category
- `DELETE /api/v1/categories/{id}` - Delete category

### Products
- `GET /api/v1/products` - Get all products
- `GET /api/v1/products/{id}` - Get product by ID
- `POST /api/v1/products` - Create new product
- `PUT /api/v1/products/{id}` - Update product
- `DELETE /api/v1/products/{id}` - Delete product

**For detailed API documentation, see [docs/api/](docs/api/)**

## 🧪 Testing

```bash
# Run all tests
make test

# Check test coverage (for CI/CD)
make check-coverage

# Generate HTML coverage report
make coverage-html

# Show coverage summary
make coverage-report
```

### Test Coverage

| Package | Coverage | Status |
|---------|----------|--------|
| handlers | 83.8% | ✅ |
| services | 100% | ✅ |
| repositories | 89.6% | ✅ |
| utils | 100% | ✅ |
| **Total (Business Logic)** | **87.3%** | ✅ |

See [Test Coverage Documentation](docs/coverage/CI_CD_COVERAGE.md) for details.

## 🛠️ Development

### Prerequisites

- Go 1.22 or higher
- PostgreSQL 12 or higher
- Make (optional, for using Makefile commands)

### Environment Variables

Create a `.env` file:

```env
DB_HOST=localhost
DB_PORT=5432
DB_USER=your_user
DB_PASSWORD=your_password
DB_NAME=kasir_db
DB_SSLMODE=disable
```

### Available Commands

```bash
make help              # Show all available commands
make test              # Run tests with race detector
make check-coverage    # Run coverage check (CI/CD)
make coverage-html     # Generate HTML coverage report
make coverage-report   # Show coverage summary
```

## 🤝 Contributing

1. Check the [development guide](docs/guides/development.md) *(coming soon)*
2. Follow the coding standards in [CLAUDE.md](CLAUDE.md)
3. Ensure tests pass: `make check-coverage`
4. Update documentation in `docs/` as needed

## 📄 License

This project is licensed under the MIT License - see the LICENSE file for details.

## 📞 Contact

For questions or issues, please open an issue in the repository.

---

**📚 For complete documentation, visit [`docs/README.md`](docs/README.md)**
