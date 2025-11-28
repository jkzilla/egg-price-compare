#!/bin/bash
set -e

echo "🔨 Building for Netlify..."

# Install Go dependencies
go mod download

# Build frontend (React + Vite) into public/
echo "📦 Installing frontend dependencies..."
cd frontend
npm install
echo "🧱 Building frontend..."
npm run build

# Build the GraphQL function
cd ../netlify/functions
GOOS=linux GOARCH=amd64 go build -o graphql graphql.go

echo "✅ Build complete"
