# CI/CD Test Coverage Configuration

## Overview

Project ini menggunakan automated test coverage check dengan threshold minimum untuk memastikan kualitas kode.

## Coverage Thresholds

- **File Coverage:** 80%
- **Package Coverage:** 80%
- **Total Coverage:** 67%

## Packages yang Diukur

Coverage check hanya mengukur **business logic packages**:
- ✅ `handlers/` - HTTP handlers
- ✅ `internal/services/` - Business logic services
- ✅ `internal/repositories/` - Database repositories
- ✅ `utils/` - Utility functions

## Packages yang Di-exclude

Package berikut **tidak** diukur coverage-nya:
- ❌ `main.go` - Entry point (tidak perlu unit test)
- ❌ `config/` - Configuration constants
- ❌ `internal/database/` - Database initialization
- ❌ `internal/mocks/` - Mock files untuk testing
- ❌ `models/` - Data models (struct definitions)

## Commands

### Menjalankan Tests

```bash
# Run all tests
make test

# Run tests dengan coverage check (untuk CI/CD)
make check-coverage

# Generate HTML coverage report
make coverage-html

# Show coverage summary
make coverage-report
```

## CI/CD Integration

### GitHub Actions

```yaml
name: Test Coverage

on: [push, pull_request]

jobs:
  test:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      
      - name: Set up Go
        uses: actions/setup-go@v4
        with:
          go-version: '1.22'
      
      - name: Run tests with coverage
        run: make check-coverage
```

### GitLab CI

```yaml
test:
  stage: test
  image: golang:1.22
  script:
    - make check-coverage
  coverage: '/Total test coverage: \d+\.\d+%/'
```

## Coverage Results

Current coverage status:

| Package | Coverage |
|---------|----------|
| handlers | 83.8% ✅ |
| services | 100% ✅ |
| repositories | 89.6% ✅ |
| utils | 100% ✅ |
| **Total (Business Logic)** | **87.3%** ✅ |

## Troubleshooting

### Coverage Check Fails

Jika `make check-coverage` gagal:

1. **File coverage < 80%**
   - Periksa file mana yang kurang coverage
   - Tambahkan test untuk edge cases yang belum tercover

2. **Package coverage < 80%**
   - Pastikan semua fungsi public memiliki test
   - Tambahkan test untuk error paths

3. **Total coverage < 67%**
   - Kemungkinan ada banyak code baru yang belum di-test
   - Review dan tambahkan test sesuai kebutuhan

### Update Thresholds

Edit file `.testcoverage.yml`:

```yaml
threshold:
  file: 80      # Minimal coverage per file
  package: 80   # Minimal coverage per package
  total: 67     # Minimal total coverage
```

## Best Practices

1. **Selalu run `make check-coverage` sebelum commit**
2. **Jangan turunkan threshold tanpa diskusi tim**
3. **Prioritaskan test untuk business logic**
4. **Gunakan table-driven tests untuk validation**
5. **Mock external dependencies (database, API calls)**

## Contact

Jika ada pertanyaan tentang test coverage, hubungi tim development.
