#!/bin/bash
set -e

echo "🗑️  Tearing down k3d cluster..."

# Delete Helm release
helm uninstall egg-price-compare -n egg-price-compare 2>/dev/null || true

# Delete k3d cluster
k3d cluster delete egg-price-compare

echo "✅ Cluster deleted"
