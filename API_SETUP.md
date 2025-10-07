# API Setup Guide

This guide will help you set up API access for Walmart and Walgreens to get real-time egg prices.

## Walmart API Setup

### Step 1: Register for Walmart Developer Account
1. Go to https://developer.walmart.com/
2. Click "Sign Up" and create an account
3. Verify your email address

### Step 2: Create an Application
1. Log in to your Walmart Developer account
2. Go to "My Account" → "Add New Application"
3. Fill in the application details:
   - **Application Name**: Egg Price Comparison
   - **Description**: Compare egg prices across stores
   - **Website**: Your website URL (or use http://localhost:8080)
4. Accept the Terms of Service
5. Click "Create Application"

### Step 3: Get Your API Key
1. Once created, you'll see your application dashboard
2. Copy the **API Key** (also called Consumer ID)
3. Save this key - you'll need it for the `.env` file

### API Endpoints Used
- **Product Search API**: `https://developer.api.walmart.com/api-proxy/service/affil/product/v2/search`
- **Documentation**: https://developer.walmart.com/doc/us/mp/us-mp-items/

## Walgreens API Setup

### Option 1: Walgreens Developer Portal (Recommended)
1. Go to https://developer.walgreens.com/
2. Click "Register" to create a developer account
3. Complete the registration process
4. Create a new application in the developer portal
5. Get your API credentials (API Key and Secret)

### Option 2: Alternative - Use Web Scraping (Legal Considerations)
**Note**: Web scraping should only be used if:
- You have permission from Walgreens
- You comply with their Terms of Service
- You implement rate limiting and respectful scraping practices

If Walgreens doesn't offer a public API, you may need to:
1. Contact Walgreens Business Development for API access
2. Use their affiliate program if available
3. Implement web scraping with proper rate limiting

### API Endpoints
- **Product Search**: `https://services.walgreens.com/api/products/search`
- **Authentication**: OAuth 2.0 (requires client credentials)

## Setting Up Your Environment

### 1. Create `.env` file
```bash
cd ~/src/egg-price-compare
cp .env.example .env
```

### 2. Add Your API Keys
Edit the `.env` file:
```bash
# Walmart API Key
WALMART_API_KEY=your_walmart_api_key_here

# Walgreens API Key
WALGREENS_API_KEY=your_walgreens_api_key_here

# Server port
PORT=8080
```

### 3. Test the Integration
```bash
# Run with Docker
docker-compose up

# Or run directly
go run server.go
```

### 4. Verify API Calls
Open http://localhost:8080 and run this query:
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

## Troubleshooting

### "API Key is invalid"
- Double-check your API key is correct
- Ensure there are no extra spaces in the `.env` file
- Verify your Walmart developer account is active

### "Rate limit exceeded"
- Walmart API has rate limits (typically 5 requests per second)
- Implement caching to reduce API calls
- Add delays between requests

### "No products found"
- The search query might need adjustment
- Try different search terms: "eggs dozen", "large eggs", etc.
- Check if the product is available in your region

### Walgreens API Not Working
- Walgreens API access may require business partnership
- Consider contacting Walgreens Developer Support
- Alternative: Use mock data or implement web scraping (with permission)

## API Rate Limits

### Walmart
- **Rate Limit**: 5 requests per second
- **Daily Limit**: Varies by plan (typically 5,000-100,000 requests/day)
- **Recommendation**: Cache results for 1-5 minutes

### Walgreens
- **Rate Limit**: Varies (check your API plan)
- **Recommendation**: Cache results and implement exponential backoff

## Caching Strategy

To avoid hitting rate limits, implement caching:

```go
// Cache prices for 5 minutes
type PriceCache struct {
    price       *model.StorePrice
    lastUpdated time.Time
}

var cache = make(map[string]*PriceCache)

func getCachedPrice(store string, fetchFunc func() (*model.StorePrice, error)) (*model.StorePrice, error) {
    if cached, ok := cache[store]; ok {
        if time.Since(cached.lastUpdated) < 5*time.Minute {
            return cached.price, nil
        }
    }
    
    price, err := fetchFunc()
    if err != nil {
        return nil, err
    }
    
    cache[store] = &PriceCache{
        price:       price,
        lastUpdated: time.Now(),
    }
    
    return price, nil
}
```

## Production Considerations

1. **Error Handling**: Add retry logic for failed API calls
2. **Monitoring**: Track API usage and errors
3. **Fallback**: Use cached data if APIs are unavailable
4. **Compliance**: Ensure you comply with API Terms of Service
5. **Security**: Never commit API keys to version control

## Support

- **Walmart API Support**: https://developer.walmart.com/support
- **Walgreens**: Contact their business development team

## Legal Notice

This application is for educational/personal use. Ensure you:
- Comply with API Terms of Service
- Don't violate rate limits
- Respect data usage policies
- Attribute data sources appropriately
