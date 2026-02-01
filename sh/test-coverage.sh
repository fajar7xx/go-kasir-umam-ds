#!/bin/sh
set -e

echo "🔍 Listing files in current directory..."
ls -la | head -n 20

echo "📦 Installing go-test-coverage..."
go install github.com/vladopajic/go-test-coverage/v2@latest

echo "🧪 Running tests with coverage..."
go test ./... -coverprofile=cover.out -covermode=atomic -coverpkg=./...

echo "📊 Coverage summary:"
go tool cover -func=cover.out | tail -n 1

echo "🔍 Checking coverage threshold..."
$(go env GOPATH)/bin/go-test-coverage --config=.testcoverage.yml || {
  echo "❌ Coverage check failed. Low coverage files:"
  go tool cover -func=cover.out | grep -v 100.0% | sort -k3 -n | head -n 20
  exit 1
}

echo "✅ All coverage thresholds met!"
