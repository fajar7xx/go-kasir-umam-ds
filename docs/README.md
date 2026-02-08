# 📚 Documentation Index

Welcome to the go-kasir-umam-ds project documentation!

## 📖 Table of Contents

- [Getting Started](#getting-started)
- [Architecture](#architecture)
- [API Reference](#api-reference)
- [Development Guides](#development-guides)
- [Testing & Coverage](#testing--coverage)
- [Implementation Plans](#implementation-plans)

---

## 🚀 Getting Started

New to the project? Start here:

1. **Setup & Installation** - `guides/getting-started.md` *(Coming soon)*
2. **Development Workflow** - `guides/development.md` *(Coming soon)*
3. **Project Structure** - See [Architecture Overview](#architecture)

---

## 🏗️ Architecture

System design and architectural decisions:

| Document | Description |
|----------|-------------|
| `architecture/overview.md` | High-level system architecture *(Coming soon)* |
| `architecture/clean-architecture.md` | Clean architecture implementation *(Coming soon)* |
| `architecture/database-schema.md` | Database design and relationships *(Coming soon)* |
| [`architecture/transaction-performance-optimization.md`](architecture/transaction-performance-optimization.md) | Batch INSERT optimization for transaction details ✅ |

### Current Architecture

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

---

## 🔌 API Reference

RESTful API endpoints documentation:

| Document | Description |
|----------|-------------|
| [`api/report-endpoints.md`](api/report-endpoints.md) | Report API documentation ✅ |
| `api/endpoints.md` | Complete API endpoints reference *(Coming soon)* |
| `api/request-response.md` | Request/Response examples *(Coming soon)* |
| `api/authentication.md` | Authentication & authorization *(Coming soon)* |

### Available Resources

- **Categories** - `/api/categories`
- **Products** - `/api/products`
- **Reports** - `/api/v1/reports` ✅

---

## 🛠️ Development Guides

How-to guides for developers:

| Document | Description | Status |
|----------|-------------|--------|
| `guides/getting-started.md` | Setup & installation guide | 📝 Planned |
| `guides/development.md` | Development workflow & best practices | 📝 Planned |
| `guides/testing.md` | Testing guidelines & examples | 📝 Planned |
| `guides/deployment.md` | Deployment & CI/CD setup | 📝 Planned |

---

## 🧪 Testing & Coverage

Test coverage configuration and reports:

| Document | Description | Status |
|----------|-------------|--------|
| [`coverage/CI_CD_COVERAGE.md`](coverage/CI_CD_COVERAGE.md) | CI/CD test coverage configuration | ✅ Available |
| [`coverage/COVERAGE_BUILD_SUMMARY.txt`](coverage/COVERAGE_BUILD_SUMMARY.txt) | Latest coverage build summary | ✅ Available |

### Quick Commands

```bash
# Run all tests
make test

# Check coverage (for CI/CD)
make check-coverage

# Generate HTML coverage report
make coverage-html
```

### Current Coverage

| Package | Coverage | Status |
|---------|----------|--------|
| handlers | 83.8% | ✅ |
| services | 100% | ✅ |
| repositories | 89.6% | ✅ |
| utils | 100% | ✅ |
| **Total (Business Logic)** | **87.3%** | ✅ |

---

## 📋 Implementation Plans

Implementation plans and RFCs:

| Document | Description | Status |
|----------|-------------|--------|
| `plans/2024-02-01-test-coverage-enhancement.md` | Test coverage 59.1% → 80%+ | ✅ Completed |

### Plan Naming Convention

Plans use the format: `YYYY-MM-DD-feature-name.md`

---

## 📁 Documentation Structure

```
docs/
├── README.md                      # This file - documentation index
├── architecture/                  # System architecture & design
├── api/                          # API documentation
├── guides/                       # How-to guides & tutorials
├── coverage/                     # Test coverage docs & reports
│   ├── CI_CD_COVERAGE.md        ✅
│   └── COVERAGE_BUILD_SUMMARY.txt ✅
└── plans/                        # Implementation plans & RFCs
```

---

## 🔍 Finding Documentation

### By Topic

- **Getting Started** → `guides/getting-started.md`
- **API Endpoints** → `api/endpoints.md`
- **Architecture** → `architecture/overview.md`
- **Testing** → `guides/testing.md`
- **CI/CD & Coverage** → `coverage/CI_CD_COVERAGE.md`

### By Role

**New Developer:**
1. `guides/getting-started.md` - Setup project
2. `architecture/overview.md` - Understand system
3. `guides/development.md` - Development workflow

**Backend Developer:**
1. `architecture/clean-architecture.md` - Architecture patterns
2. `api/endpoints.md` - API reference
3. `guides/testing.md` - Testing guidelines

**DevOps Engineer:**
1. `guides/deployment.md` - Deployment process
2. `coverage/CI_CD_COVERAGE.md` - CI/CD setup

---

## 📝 Contributing to Documentation

When adding new documentation:

1. Choose the correct subdirectory (`architecture/`, `api/`, `guides/`, etc.)
2. Follow naming conventions (lowercase, hyphens, descriptive)
3. Use the appropriate template (see `../CLAUDE.md`)
4. Update this index file
5. Link related documents

### Documentation Standards

See [`../CLAUDE.md`](../CLAUDE.md) for:
- File placement rules
- Naming conventions
- Document templates
- Code documentation guidelines

---

## 🔗 External Resources

- [Project Repository](#)
- [Issue Tracker](#)
- [CI/CD Pipeline](#)

---

## 📞 Contact

For questions about documentation:
- Open an issue in the repository
- Contact the development team

---

**Last Updated:** 2026-02-08
**Documentation Version:** 1.0.0
