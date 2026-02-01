#!/bin/sh
set -e

echo "Installing go-test-coverage..."
go install github.com/vladopajic/go-test-coverage/v2@latest

echo "Running tests with coverage..."
go test ./... -coverprofile=./cover.out -covermode=atomic -coverpkg=./...

echo "Checking coverage threshold..."
$(go env GOPATH)/bin/go-test-coverage --config=./.testcoverage.yml

echo "✅ All tests passed with sufficient coverage"
