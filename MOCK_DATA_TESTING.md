# Mock Data Testing Guide

The application now includes enhanced mock data for both Walmart and Walgreens APIs when API keys are not configured.

## How Mock Data Works

Both APIs check for the presence of API keys. If no keys are found:
- **Walmart**: Returns mock data when `WALMART_API_KEY` is not set
- **Walgreens**: Returns mock data when both `WALGREENS_API_KEY` and `SEARCHAPI_KEY` are not set

## Mock Data Features

### Walmart Mock Data
- **Base Price**: Varies between $3.48 and $4.98 based on zipcode
- **Promotions**: 60% chance of having a rollback/promo
- **Digital Coupons**: 30% chance of having a digital coupon ($0.25 off)
- **Stock Status**: 90% in stock
- **Product**: Great Value Large White Eggs, 12 Count

### Walgreens Mock Data
- **Base Price**: Varies between $3.99 and $5.29 based on zipcode
- **Promotions**: 70% chance of having a sale price
- **Digital Coupons**: 50% chance of having a digital coupon ($0.50 off)
- **Rewards**: 20% chance of having a myWalgreens 10% off offer
- **Stock Status**: 85% in stock
- **Product**: Walgreens Grade A Large White Eggs, 12 ct

## Testing Different Scenarios

The mock data uses a simple hash of the zipcode to generate consistent but varied results. Try different zipcodes to see different prices, promotions, and stock statuses:

```graphql
query {
  eggPrices(zipcode: "10001") {
    walmart {
      basePrice
      promoPrice
      finalPrice
      inStock
      digitalOffers {
        description
        discountAmount
      }
    }
    walgreens {
      basePrice
      promoPrice
      finalPrice
      inStock
      digitalOffers {
        description
        discountAmount
        discountPercent
      }
    }
    cheapest
    priceDifference
  }
}
```

### Example Zipcodes to Try
- `10001` - New York, NY
- `90210` - Beverly Hills, CA
- `60601` - Chicago, IL
- `33101` - Miami, FL
- `98101` - Seattle, WA

Each zipcode will produce different:
- Base prices
- Promotional discounts
- Digital offer availability
- Stock status
- Pickup ETAs

## Running the Server

```bash
# Make sure no API keys are set (or unset them)
unset WALMART_API_KEY
unset WALGREENS_API_KEY
unset SEARCHAPI_KEY

# Run the server
go run server.go
```

Then visit http://localhost:8080/ to access the GraphQL Playground.

## Switching to Real APIs

When you're ready to use real APIs, simply set the environment variables:

```bash
export WALMART_API_KEY="your-walmart-key"
export WALGREENS_API_KEY="your-walgreens-key"
export SEARCHAPI_KEY="your-searchapi-key"
```

The application will automatically switch from mock data to real API calls.
