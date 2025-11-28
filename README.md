# Egg Price Comparison - 1P Retail Pricing

A GraphQL API that compares **retail consumer prices** for one dozen eggs from Walmart and Walgreens.

**Important**: This is a **1P (retail/consumer pricing)** app, NOT a 3P (Marketplace seller) tool.

## Features

- **1P Retail Pricing**: Walmart Affiliates API + Walgreens Store Inventory/Digital Offers APIs
- **Zipcode-based pricing**: Get prices specific to your location
- **Digital offers**: Track clip-able coupons and promotions
- **Real-time availability**: In-stock status and pickup ETA
- **GraphQL API**: Flexible querying with full price breakdown
- **Price history tracking**: Monitor price trends over time
- **Multiple deployment options**: Docker, Kubernetes, Netlify

## Tech Stack

- **Backend**: Go + GraphQL (gqlgen)
- **APIs**: 
  - Walmart Affiliates Product Lookup API (1P retail)
  - Walgreens Store Inventory API (1P retail)
  - Walgreens Digital Offers API (1P retail)
  - Third-party price provider (SearchAPI/SerpApi for Walgreens pricing)
- **Deployment**: Docker, Kubernetes (Helm), k3d, Netlify

## Architecture: 1P vs 3P

**This app is 1P (retail) focused:**
- ✅ Walmart Affiliates Product Lookup API → Consumer pricing
- ✅ Walgreens Store Inventory + Digital Offers → In-stock + coupons
- ✅ Third-party data providers → Walgreens retail pricing

**NOT 3P (Marketplace):**
- ❌ Walmart Marketplace Seller APIs (for managing seller inventory)
- ❌ 3P seller catalog/order management

When asked "Do you support 1P or 3P?": **Answer: 1P retail/consumer pricing**

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

# Set environment variables
export WALMART_AFFILIATE_ID=your_affiliate_id
export WALMART_API_KEY=your_walmart_key
export WALGREENS_API_KEY=your_walgreens_key
export SEARCHAPI_KEY=your_searchapi_key

# Run server
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
netlify env:set WALMART_AFFILIATE_ID your_affiliate_id
netlify env:set WALMART_API_KEY your_walmart_key
netlify env:set WALGREENS_API_KEY your_walgreens_key
netlify env:set SEARCHAPI_KEY your_searchapi_key
netlify deploy --prod
```

## Deployment Options

### 🐳 Docker

```bash
docker build -t egg-price-compare .
docker run -p 8080:8080 \
  -e WALMART_AFFILIATE_ID=your_affiliate_id \
  -e WALMART_API_KEY=your_walmart_key \
  -e WALGREENS_API_KEY=your_walgreens_key \
  -e SEARCHAPI_KEY=your_searchapi_key \
  egg-price-compare
```

### ☸️ Kubernetes with Helm

```bash
helm install egg-price-compare ./helm/egg-price-compare \
  --namespace egg-price-compare \
  --create-namespace \
  --set config.walmartAffiliateId=your_affiliate_id \
  --set config.walmartApiKey=your_walmart_key \
  --set config.walgreensApiKey=your_walgreens_key \
  --set config.searchapiKey=your_searchapi_key
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
netlify env:set WALMART_AFFILIATE_ID your_affiliate_id
netlify env:set WALMART_API_KEY your_walmart_key
netlify env:set WALGREENS_API_KEY your_walgreens_key
netlify env:set SEARCHAPI_KEY your_searchapi_key
netlify deploy --prod
```

See [netlify/README.md](netlify/README.md) for details.

## API Usage

### GraphQL Query

```graphql
query {
  eggPrices(zipcode: "94102") {
    walmart {
      store
      sku
      upc
      zipcode
      basePrice
      promoPrice
      finalPrice
      productName
      productUrl
      inStock
      pickupEta
      digitalOffers {
        offerId
        description
        discountAmount
        discountPercent
        expiresAt
      }
      lastUpdated
    }
    walgreens {
      store
      sku
      upc
      storeId
      zipcode
      basePrice
      promoPrice
      finalPrice
      productName
      productUrl
      inStock
      pickupEta
      digitalOffers {
        offerId
        description
        discountAmount
        expiresAt
      }
      lastUpdated
    }
    cheapest
    priceDifference
    lastUpdated
  }
}
```

### cURL Example

```bash
curl -X POST http://localhost:8080/graphql \
  -H "Content-Type: application/json" \
  -d '{"query":"{ eggPrices(zipcode: \"94102\") { walmart { finalPrice inStock } walgreens { finalPrice inStock } cheapest priceDifference } }"}'
```

## Environment Variables

```bash
# Walmart Affiliates (1P Retail)
WALMART_AFFILIATE_ID=your_publisher_id
WALMART_API_KEY=your_walmart_affiliate_api_key

# Walgreens (1P Retail)
WALGREENS_API_KEY=your_walgreens_api_key
WALGREENS_API_SECRET=your_walgreens_api_secret

# Third-Party Price Provider (choose one)
SEARCHAPI_KEY=your_searchapi_key
# SERPAPI_KEY=your_serpapi_key
# APIFY_KEY=your_apify_key

# Server Configuration
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

### API Documentation (1P Retail)
- [Walmart Affiliates Program](https://affiliates.walmart.com/)
- [Walmart Affiliates Product Lookup API](https://developer.walmart.com/api/us/affil/product/v2)
- [Walgreens Developer Portal](https://developer.walgreens.com/)
- [Walgreens Store Inventory API](https://developer.walgreens.com/store-inventory)
- [Walgreens Digital Offers API](https://developer.walgreens.com/digital-offers)

### Third-Party Data Providers
- [SearchAPI - Walmart](https://www.searchapi.io/docs/walmart)
- [SearchAPI - Walgreens](https://www.searchapi.io/docs/walgreens)
- [SerpApi - Walmart](https://serpapi.com/walmart-search-api)
- [SerpApi - Walgreens](https://serpapi.com/walgreens-search-api)
- [Apify Scrapers](https://apify.com/)

### Development Tools
- [GraphQL Docs](https://graphql.org/)
- [gqlgen Docs](https://gqlgen.com/)
- [Helm Docs](https://helm.sh/docs/)
- [k3d Docs](https://k3d.io/)
- [Netlify Docs](https://docs.netlify.com/)

## Minimal Architecture (Go Handler Example)

Here's the core flow for merging Walmart + Walgreens data:

```go
type RetailerPrice struct {
    SKU         string
    UPC         string
    StoreID     string
    Zipcode     string
    BasePrice   float64
    PromoPrice  float64
    FinalPrice  float64
    InStock     bool
    PickupETA   string
    Link        string
}

func GetEggPrices(zipcode string) ([]RetailerPrice, error) {
    // 1. Walmart: Affiliates Product Lookup
    walmartPrice := fetchWalmartAffiliates(zipcode)
    
    // 2. Walgreens: Third-party price + Store Inventory + Digital Offers
    walgreensPrice := fetchWalgreensPrice(zipcode)
    walgreensInventory := fetchWalgreensInventory(zipcode)
    walgreensOffers := fetchWalgreensOffers()
    
    // 3. Merge and normalize
    prices := []RetailerPrice{walmartPrice, walgreensPrice}
    
    // 4. Rank by finalPrice
    sort.Slice(prices, func(i, j int) bool {
        return prices[i].FinalPrice < prices[j].FinalPrice
    })
    
    return prices, nil
}
```

## Key Clarifications

### When Asked: "1P or 3P?"

**Answer**: "This app supports **1P retail/consumer pricing** for cross-retailer price comparison. It does NOT support 3P Marketplace sellers."

### API Breakdown

| Retailer | API | Type | Purpose |
|----------|-----|------|----------|
| Walmart | Affiliates Product Lookup | 1P Retail | Price + availability by zipcode |
| Walgreens | Store Inventory | 1P Retail | In-stock signal |
| Walgreens | Digital Offers | 1P Retail | Clip-able coupons |
| Walgreens | SearchAPI/SerpApi | 3rd Party | Retail price (Walgreens doesn't expose pricing API) |
