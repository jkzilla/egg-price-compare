# k3d Local Development

Local Kubernetes development environment using k3d.

## Prerequisites

```bash
# macOS
brew install k3d kubectl helm docker

# Verify installation
k3d version
kubectl version --client
helm version
```

## Quick Start

### 1. Set API Keys

```bash
export WALMART_API_KEY=your_walmart_api_key
export WALGREENS_API_KEY=your_walgreens_api_key
```

### 2. Create Cluster and Deploy

```bash
cd k3d
./setup.sh
```

This will:
- Create a k3d cluster with 1 server and 2 agents
- Install NGINX Ingress Controller
- Build and push Docker image to local registry
- Deploy the application with Helm

### 3. Access the Application

- **GraphQL Playground**: http://localhost:8080
- **GraphQL API**: http://localhost:8080/graphql

### 4. Teardown

```bash
cd k3d
./teardown.sh
```

## Manual Steps

### Create Cluster Only

```bash
k3d cluster create egg-price-compare --config cluster-config.yaml
```

### Build and Deploy

```bash
# Build image
docker build -t localhost:5000/egg-price-compare:latest ..
docker push localhost:5000/egg-price-compare:latest

# Deploy with Helm
helm install egg-price-compare ../helm/egg-price-compare \
  --namespace egg-price-compare \
  --create-namespace \
  --set image.repository=localhost:5000/egg-price-compare
```

### Update Deployment

```bash
# Rebuild image
docker build -t localhost:5000/egg-price-compare:latest ..
docker push localhost:5000/egg-price-compare:latest

# Restart pods
kubectl rollout restart deployment/egg-price-compare -n egg-price-compare
```

## Troubleshooting

### Pods not starting

```bash
kubectl get pods -n egg-price-compare
kubectl describe pod <pod-name> -n egg-price-compare
kubectl logs <pod-name> -n egg-price-compare
```

### Ingress not working

```bash
kubectl get ingress -n egg-price-compare
kubectl describe ingress -n egg-price-compare
kubectl logs -n ingress-nginx deployment/ingress-nginx-controller
```

### Registry issues

```bash
# Check registry
docker ps | grep registry

# Test push
docker push localhost:5000/egg-price-compare:latest
```

## Cluster Configuration

- **Name**: egg-price-compare
- **Servers**: 1
- **Agents**: 2
- **Ports**: 8080 (HTTP), 8443 (HTTPS)
- **Registry**: localhost:5000

## Resources

- k3d docs: https://k3d.io
- k3s docs: https://k3s.io
- Helm docs: https://helm.sh
