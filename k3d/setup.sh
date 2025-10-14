#!/bin/bash
set -e

echo "🚀 Setting up k3d cluster for Egg Price Comparison"
echo "=================================================="
echo ""

# Check if k3d is installed
if ! command -v k3d &> /dev/null; then
    echo "❌ k3d is not installed"
    echo "Install with: brew install k3d"
    exit 1
fi

# Check if kubectl is installed
if ! command -v kubectl &> /dev/null; then
    echo "❌ kubectl is not installed"
    echo "Install with: brew install kubectl"
    exit 1
fi

# Check if helm is installed
if ! command -v helm &> /dev/null; then
    echo "❌ helm is not installed"
    echo "Install with: brew install helm"
    exit 1
fi

# Create k3d cluster
echo "📦 Creating k3d cluster..."
k3d cluster create egg-price-compare \
  --config cluster-config.yaml \
  --registry-create egg-registry:0.0.0.0:5000

# Wait for cluster to be ready
echo "⏳ Waiting for cluster to be ready..."
kubectl wait --for=condition=Ready nodes --all --timeout=60s

# Install NGINX Ingress
echo "🔧 Installing NGINX Ingress Controller..."
helm repo add ingress-nginx https://kubernetes.github.io/ingress-nginx
helm repo update
helm install ingress-nginx ingress-nginx/ingress-nginx \
  --namespace ingress-nginx \
  --create-namespace \
  --set controller.service.type=NodePort \
  --set controller.hostPort.enabled=true \
  --wait

# Build and push image to local registry
echo "🐳 Building Docker image..."
cd ..
docker build -t localhost:5000/egg-price-compare:latest .
docker push localhost:5000/egg-price-compare:latest

# Install Helm chart
echo "📊 Installing Helm chart..."
helm install egg-price-compare ./helm/egg-price-compare \
  --namespace egg-price-compare \
  --create-namespace \
  --set image.repository=localhost:5000/egg-price-compare \
  --set image.tag=latest \
  --set ingress.hosts[0].host=localhost \
  --set config.walmartApiKey="${WALMART_API_KEY:-test-key}" \
  --set config.walgreensApiKey="${WALGREENS_API_KEY:-test-key}" \
  --wait

echo ""
echo "✅ Setup complete!"
echo ""
echo "📝 Access the application:"
echo "   GraphQL Playground: http://localhost:8080"
echo "   GraphQL API: http://localhost:8080/graphql"
echo ""
echo "🔍 Useful commands:"
echo "   kubectl get pods -n egg-price-compare"
echo "   kubectl logs -f deployment/egg-price-compare -n egg-price-compare"
echo "   helm list -n egg-price-compare"
echo ""
echo "🗑️  To delete the cluster:"
echo "   k3d cluster delete egg-price-compare"
