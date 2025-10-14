#!/bin/bash
set -e

echo "🔨 Building for Netlify..."

# Install dependencies
go mod download

# Build the function
cd netlify/functions
GOOS=linux GOARCH=amd64 go build -o graphql graphql.go

echo "✅ Build complete"
