# Egg Price Comparison

A GraphQL API and web application that compares the prices of one dozen eggs from Walmart and Walgreens.

## Features

- Real-time price comparison from Walmart and Walgreens APIs
- GraphQL API for querying egg prices
- React frontend with price visualization
- Automatic price updates
- Price history tracking
- Multiple deployment options (Docker, Kubernetes, Netlify)

## Tech Stack

- **Backend**: Go + GraphQL (gqlgen)
- **Frontend**: React + TypeScript + Vite
- **APIs**: Walmart API, Walgreens API
- **Deployment**: Docker, Kubernetes (Helm), k3d, Netlify

## Quick Start

### Option 1: Docker Compose (Fastest)

```bash
git clone https://github.com/jkzilla/egg-price-compare.git
cd egg-price-compare
docker-compose up
```

Access at http://localhost:8080

### Option 2: Local Development

```bash
# Install dependencies
go mod download

# Run server
export WALMART_API_KEY=your_key
export WALGREENS_API_KEY=your_key
go run server.go
```

### Option 3: k3d (Local Kubernetes)

```bash
cd k3d
./setup.sh
```

Access at http://localhost:8080

### Option 4: Netlify

```bash
netlify init
netlify env:set WALMART_API_KEY your_key
netlify env:set WALGREENS_API_KEY your_key
netlify deploy --prod
```

## Deployment Options

### 🐳 Docker

```bash
docker build -t egg-price-compare .
docker run -p 8080:8080 \
  -e WALMART_API_KEY=your_key \
  -e WALGREENS_API_KEY=your_key \
  egg-price-compare
```

### ☸️ Kubernetes with Helm

```bash
helm install egg-price-compare ./helm/egg-price-compare \
  --namespace egg-price-compare \
  --create-namespace \
  --set config.walmartApiKey=your_key \
  --set config.walgreensApiKey=your_key
```

### 🏠 k3d Local Cluster

```bash
cd k3d
./setup.sh  # Creates cluster and deploys app
./teardown.sh  # Removes cluster
```

See [k3d/README.md](k3d/README.md) for details.

### 🌐 Netlify Functions

```bash
netlify init
netlify env:set WALMART_API_KEY your_key
netlify deploy --prod
```

See [netlify/README.md](netlify/README.md) for details.

## API Usage

### GraphQL Query

```graphql
query {
  eggPrices {
    walmart {
      price
      productName
      inStock
      lastUpdated
    }
    walgreens {
      price
      productName
      inStock
      lastUpdated
    }
    cheapest
    priceDifference
  }
}
```

### cURL Example

```bash
curl -X POST http://localhost:8080/graphql \
  -H "Content-Type: application/json" \
  -d '{"query":"{ eggPrices { walmart { price } walgreens { price } cheapest } }"}'
```

## Environment Variables

```bash
WALMART_API_KEY=your_walmart_api_key
WALGREENS_API_KEY=your_walgreens_api_key
PORT=8080
```

## Project Structure

```
egg-price-compare/
├── server.go                 # Main server
├── graph/                    # GraphQL schema and resolvers
├── api/                      # External API clients
├── helm/                     # Helm charts for Kubernetes
│   └── egg-price-compare/
├── k3d/                      # Local k3d setup
│   ├── setup.sh
│   ├── teardown.sh
│   └── cluster-config.yaml
├── netlify/                  # Netlify deployment
│   ├── functions/
│   └── README.md
├── Dockerfile                # Docker configuration
├── docker-compose.yml        # Docker Compose setup
└── README.md
```

## Development

### Run Tests

```bash
go test ./...
```

### Generate GraphQL Code

```bash
go run github.com/99designs/gqlgen generate
```

### Build Binary

```bash
go build -o egg-price-compare .
```

## API Setup

See [API_SETUP.md](API_SETUP.md) for detailed instructions on obtaining API keys from Walmart and Walgreens.

## Monitoring

### Kubernetes

```bash
# Check pods
kubectl get pods -n egg-price-compare

# View logs
kubectl logs -f deployment/egg-price-compare -n egg-price-compare

# Check ingress
kubectl get ingress -n egg-price-compare
```

### Docker

```bash
docker logs -f <container-id>
```

## Troubleshooting

### API Rate Limits

Both Walmart and Walgreens APIs have rate limits. Implement caching:

```go
// Cache prices for 5 minutes
var priceCache = make(map[string]*CachedPrice)
```

### CORS Issues

Update CORS settings in `server.go`:

```go
c := cors.New(cors.Options{
    AllowedOrigins: []string{"https://your-domain.com"},
})
```

### Kubernetes Issues

```bash
# Describe pod
kubectl describe pod <pod-name> -n egg-price-compare

# Check events
kubectl get events -n egg-price-compare --sort-by='.lastTimestamp'
```

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Submit a pull request

## License

MIT

## Resources

- [Walmart API Docs](https://developer.walmart.com/)
- [Walgreens API Docs](https://developer.walgreens.com/)
- [GraphQL Docs](https://graphql.org/)
- [Helm Docs](https://helm.sh/docs/)
- [k3d Docs](https://k3d.io/)
- [Netlify Docs](https://docs.netlify.com/)
