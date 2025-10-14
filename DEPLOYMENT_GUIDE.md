# Deployment Guide

Complete guide for deploying Egg Price Comparison API.

## Table of Contents

1. [Local Development](#local-development)
2. [k3d (Local Kubernetes)](#k3d-local-kubernetes)
3. [Production Kubernetes](#production-kubernetes)
4. [Netlify (Serverless)](#netlify-serverless)
5. [Docker](#docker)

---

## Local Development

### Prerequisites
- Go 1.21+
- API keys from Walmart and Walgreens

### Setup

```bash
# Clone repository
git clone https://github.com/jkzilla/egg-price-compare.git
cd egg-price-compare

# Install dependencies
go mod download

# Set environment variables
export WALMART_API_KEY=your_walmart_key
export WALGREENS_API_KEY=your_walgreens_key

# Run server
go run server.go
```

Access at http://localhost:8080

---

## k3d (Local Kubernetes)

### Prerequisites
```bash
brew install k3d kubectl helm docker
```

### Quick Start

```bash
cd k3d
./setup.sh
```

This automatically:
1. Creates k3d cluster (1 server, 2 agents)
2. Installs NGINX Ingress
3. Builds and pushes Docker image
4. Deploys with Helm

### Manual Steps

```bash
# Create cluster
k3d cluster create egg-price-compare --config cluster-config.yaml

# Build image
docker build -t localhost:5000/egg-price-compare:latest .
docker push localhost:5000/egg-price-compare:latest

# Deploy
helm install egg-price-compare ./helm/egg-price-compare \
  --namespace egg-price-compare \
  --create-namespace \
  --set image.repository=localhost:5000/egg-price-compare \
  --set config.walmartApiKey=$WALMART_API_KEY \
  --set config.walgreensApiKey=$WALGREENS_API_KEY
```

### Cleanup

```bash
cd k3d
./teardown.sh
```

---

## Production Kubernetes

### Prerequisites
- Kubernetes cluster (EKS, GKE, AKS, or self-hosted)
- kubectl configured
- Helm 3 installed
- Docker registry access

### 1. Build and Push Image

```bash
# Build image
docker build -t your-registry/egg-price-compare:v1.0.0 .

# Push to registry
docker push your-registry/egg-price-compare:v1.0.0
```

### 2. Create Secrets

```bash
kubectl create secret generic egg-price-compare-secret \
  --from-literal=walmart-api-key=$WALMART_API_KEY \
  --from-literal=walgreens-api-key=$WALGREENS_API_KEY \
  -n egg-price-compare
```

### 3. Deploy with Helm

```bash
helm install egg-price-compare ./helm/egg-price-compare \
  --namespace egg-price-compare \
  --create-namespace \
  --set image.repository=your-registry/egg-price-compare \
  --set image.tag=v1.0.0 \
  --set ingress.hosts[0].host=egg-prices.yourdomain.com \
  --set config.walmartApiKey=$WALMART_API_KEY \
  --set config.walgreensApiKey=$WALGREENS_API_KEY
```

### 4. Configure DNS

Point your domain to the load balancer:

```bash
# Get load balancer IP
kubectl get svc -n ingress-nginx

# Add A record in your DNS provider
egg-prices.yourdomain.com → <load-balancer-ip>
```

### 5. Enable HTTPS (Optional)

```bash
# Install cert-manager
helm repo add jetstack https://charts.jetstack.io
helm install cert-manager jetstack/cert-manager \
  --namespace cert-manager \
  --create-namespace \
  --set installCRDs=true

# Update values.yaml
ingress:
  tls:
    - secretName: egg-prices-tls
      hosts:
        - egg-prices.yourdomain.com
```

---

## Netlify (Serverless)

### Prerequisites
```bash
npm install -g netlify-cli
```

### Setup

```bash
# Login
netlify login

# Initialize site
netlify init

# Set environment variables
netlify env:set WALMART_API_KEY your_walmart_key
netlify env:set WALGREENS_API_KEY your_walgreens_key
```

### Deploy

```bash
# Deploy to production
netlify deploy --prod

# Or deploy preview
netlify deploy
```

### Local Development

```bash
netlify dev
```

Access at http://localhost:8888

---

## Docker

### Build

```bash
docker build -t egg-price-compare .
```

### Run

```bash
docker run -p 8080:8080 \
  -e WALMART_API_KEY=$WALMART_API_KEY \
  -e WALGREENS_API_KEY=$WALGREENS_API_KEY \
  egg-price-compare
```

### Docker Compose

```bash
docker-compose up
```

---

## Monitoring

### Kubernetes

```bash
# Check pods
kubectl get pods -n egg-price-compare

# View logs
kubectl logs -f deployment/egg-price-compare -n egg-price-compare

# Check resource usage
kubectl top pods -n egg-price-compare
```

### Netlify

```bash
# View function logs
netlify functions:log graphql

# Check build logs
netlify build
```

---

## Troubleshooting

### Pods CrashLooping

```bash
kubectl describe pod <pod-name> -n egg-price-compare
kubectl logs <pod-name> -n egg-price-compare
```

Common causes:
- Missing API keys
- Invalid API keys
- Image pull errors

### API Rate Limiting

Implement caching in `api/` directory:

```go
var cache = make(map[string]*CachedPrice)
```

### CORS Issues

Update `server.go`:

```go
c := cors.New(cors.Options{
    AllowedOrigins: []string{"https://yourdomain.com"},
})
```

---

## Performance Tuning

### Kubernetes

```yaml
# values.yaml
resources:
  limits:
    cpu: 1000m
    memory: 512Mi
  requests:
    cpu: 200m
    memory: 256Mi

replicaCount: 3
```

### Caching

Add Redis for caching:

```bash
helm install redis bitnami/redis -n egg-price-compare
```

---

## Security

### 1. Use Secrets

Never commit API keys. Use Kubernetes secrets or environment variables.

### 2. Enable HTTPS

Use cert-manager for automatic SSL certificates.

### 3. Network Policies

```yaml
apiVersion: networking.k8s.io/v1
kind: NetworkPolicy
metadata:
  name: egg-price-compare-netpol
spec:
  podSelector:
    matchLabels:
      app: egg-price-compare
  ingress:
    - from:
        - namespaceSelector:
            matchLabels:
              name: ingress-nginx
```

### 4. RBAC

Create service account with minimal permissions.

---

## CI/CD

### GitHub Actions

```yaml
name: Deploy
on:
  push:
    branches: [main]
jobs:
  deploy:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      - name: Build and push
        run: |
          docker build -t ${{ secrets.REGISTRY }}/egg-price-compare:${{ github.sha }} .
          docker push ${{ secrets.REGISTRY }}/egg-price-compare:${{ github.sha }}
      - name: Deploy with Helm
        run: |
          helm upgrade --install egg-price-compare ./helm/egg-price-compare \
            --set image.tag=${{ github.sha }}
```

---

## Support

- GitHub Issues: https://github.com/jkzilla/egg-price-compare/issues
- Documentation: See README.md
- API Setup: See API_SETUP.md
