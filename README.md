# Egg Price Comparison

A GraphQL API and web application that compares the prices of one dozen eggs from Walmart and Walgreens.

## Features

- Real-time price comparison from Walmart and Walgreens APIs
- GraphQL API for querying egg prices
- React frontend with price visualization
- Automatic price updates
- Price history tracking

## Tech Stack

- **Backend**: Go + GraphQL (gqlgen)
- **Frontend**: React + TypeScript + Vite
- **APIs**: Walmart API, Walgreens API
- **Deployment**: Docker + Kubernetes + Helm

## Quick Start

```bash
# Clone the repository
git clone https://github.com/jkzilla/egg-price-compare.git
cd egg-price-compare

# Run with Docker Compose
docker-compose up

# Or run locally
cd backend && go run server.go
cd frontend && npm install && npm run dev
```

## API Usage

```graphql
query {
  eggPrices {
    walmart {
      price
      productName
      lastUpdated
    }
    walgreens {
      price
      productName
      lastUpdated
    }
    cheapest
    priceDifference
  }
}
```

## Environment Variables

```
WALMART_API_KEY=your_walmart_api_key
WALGREENS_API_KEY=your_walgreens_api_key
PORT=8080
```

## License

MIT
