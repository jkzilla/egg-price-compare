# API Setup Guide - 1P Retail Pricing

This guide helps you set up API access for **1P retail/consumer pricing** (not 3P Marketplace).

## Architecture Overview

This app compares **retail consumer prices** across stores:
- **Walmart**: Uses Walmart Affiliates Product Lookup API (1P retail)
- **Walgreens**: Uses Store Inventory + Digital Offers APIs + third-party price provider

**Important**: This is a 1P (retail) price-comparison app, NOT a 3P (Marketplace seller) tool.

---

## Walmart API Setup (1P Retail)

### Step 1: Join Walmart Affiliates Program

1. Go to https://affiliates.walmart.com/
2. Click "Join Now" and complete the application
3. Provide your website/app details
4. Wait for approval (typically 1-3 business days)

### Step 2: Get Your Affiliate Credentials

1. Log in to your Walmart Affiliates account
2. Go to "API" or "Developer Tools" section
3. Copy your:
   - **Publisher ID** (Affiliate ID)
   - **API Key** (Consumer ID)
4. Save these for the `.env` file

### Step 3: Access Product Lookup API

The **Walmart Affiliates Product Lookup API** is purpose-built for price-comparison use cases.

- **API Endpoint**: `https://developer.api.walmart.com/api-proxy/service/affil/product/v2/search`
- **Documentation**: https://developer.walmart.com/api/us/affil/product/v2
- **Features**:
  - Real-time pricing
  - Product availability
  - Store-level data (limited zipcode support)
  - Special offers (clearance, special buy)

### Alternative: Third-Party Data Provider

If you can't get Walmart Affiliates access immediately, use a reputable data provider:
- **SearchAPI**: https://www.searchapi.io/docs/walmart
- **SerpApi**: https://serpapi.com/walmart-search-api

These services are designed for price-comparison apps and respect ToS.

---

## Walgreens API Setup (1P Retail)

### Reality Check

**Walgreens does NOT expose a general product pricing API** like Walmart Affiliates.

However, Walgreens provides these 1P retail APIs:
1. **Store Inventory API** - Check in-stock status
2. **Digital Offers API** - Retrieve clip-able coupons
3. **Add-to-Cart API** - Push users to checkout

### Step 1: Register for Walgreens Developer Portal

1. Go to https://developer.walgreens.com/
2. Click "Register" to create a developer account
3. Complete the registration and business verification
4. Create a new application
5. Get your API credentials:
   - **API Key**
   - **API Secret** (for OAuth 2.0)

### Step 2: Enable Required APIs

In your Walgreens Developer Portal, enable:
- **Store Inventory API** - For in-stock checks
- **Digital Offers API** - For coupon/promo data

**API Documentation**:
- Store Inventory: https://developer.walgreens.com/store-inventory
- Digital Offers: https://developer.walgreens.com/digital-offers

### Step 3: Get Pricing Data (Third-Party Provider)

Since Walgreens doesn't expose retail pricing, use a **third-party data provider**:

#### Option 1: SearchAPI (Recommended)
- **Website**: https://www.searchapi.io/
- **Docs**: https://www.searchapi.io/docs/walgreens
- **Pricing**: Pay-per-request or subscription
- **Features**: Real-time Walgreens product search with pricing

#### Option 2: SerpApi
- **Website**: https://serpapi.com/
- **Docs**: https://serpapi.com/walgreens-search-api
- **Pricing**: Free tier available

#### Option 3: Apify
- **Website**: https://apify.com/
- **Search**: "Walgreens scraper"
- **Note**: Pre-built scrapers that respect ToS

### Architecture

```
Walgreens Price Flow:
1. Third-party provider → Base price + promo price
2. Walgreens Store Inventory API → In-stock status + pickup ETA
3. Walgreens Digital Offers API → Clip-able coupons
4. Combine all → Final price with offers
```

---

## Setting Up Your Environment

### 1. Create `.env` file
```bash
cd ~/src/haileysgarden/egg-price-compare
touch .env
```

### 2. Add Your API Keys

Edit the `.env` file:

```bash
# Walmart Affiliates (1P Retail)
WALMART_AFFILIATE_ID=your_publisher_id_here
WALMART_API_KEY=your_walmart_affiliate_api_key_here

# Walgreens (1P Retail)
WALGREENS_API_KEY=your_walgreens_api_key_here
WALGREENS_API_SECRET=your_walgreens_api_secret_here

# Third-Party Price Provider (choose one)
SEARCHAPI_KEY=your_searchapi_key_here
# SERPAPI_KEY=your_serpapi_key_here
# APIFY_KEY=your_apify_key_here

# Server Configuration
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
  eggPrices(zipcode: "94102") {
    walmart {
      store
      basePrice
      promoPrice
      finalPrice
      productName
      inStock
      pickupEta
      digitalOffers {
        description
        discountAmount
      }
    }
    walgreens {
      store
      basePrice
      promoPrice
      finalPrice
      productName
      inStock
      pickupEta
      digitalOffers {
        description
        discountAmount
      }
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

### Walgreens Pricing Not Available
- **Root cause**: Walgreens doesn't expose retail pricing API
- **Solution**: Use a third-party data provider (SearchAPI, SerpApi, Apify)
- **Fallback**: Use mock data during development
- **Contact**: Walgreens Business Development for potential API access

---

## API Rate Limits

### Walmart Affiliates
- **Rate Limit**: 5 requests per second
- **Daily Limit**: Varies by affiliate tier (typically 5,000-100,000/day)
- **Recommendation**: Cache results for 5-15 minutes

### Walgreens
- **Store Inventory API**: 10 requests per second
- **Digital Offers API**: 10 requests per second
- **Recommendation**: Cache inventory for 5 minutes, offers for 1 hour

### Third-Party Providers
- **SearchAPI**: 100 requests/month (free), unlimited (paid)
- **SerpApi**: 100 searches/month (free), 5,000+ (paid)
- **Recommendation**: Cache aggressively (15-30 minutes)

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

---

## Support & Resources

### Official Support
- **Walmart Affiliates**: https://affiliates.walmart.com/support
- **Walgreens Developer**: https://developer.walgreens.com/support
- **SearchAPI**: https://www.searchapi.io/support
- **SerpApi**: https://serpapi.com/support

### Key Documentation
- Walmart Affiliates Product Lookup: https://developer.walmart.com/api/us/affil/product/v2
- Walgreens Store Inventory: https://developer.walgreens.com/store-inventory
- Walgreens Digital Offers: https://developer.walgreens.com/digital-offers

---

## Legal & Compliance

### Important Notes

1. **This is a 1P (retail) app**: You're comparing consumer prices, not managing Marketplace inventory
2. **Comply with ToS**: Follow all API Terms of Service
3. **Rate limiting**: Implement caching and respect rate limits
4. **Attribution**: Credit data sources appropriately
5. **Third-party providers**: Ensure they have proper data licensing

### When Asked "1P or 3P?"

**Answer**: "This app supports **1P retail/consumer pricing** for price comparison. It does NOT support 3P Marketplace sellers."

---

## Quick Reference

| Retailer | API Type | Purpose | Provider |
|----------|----------|---------|----------|
| Walmart | Affiliates Product Lookup | Pricing + Availability | Walmart (1P) |
| Walgreens | Store Inventory | In-stock status | Walgreens (1P) |
| Walgreens | Digital Offers | Coupons/Promos | Walgreens (1P) |
| Walgreens | Pricing Data | Retail prices | SearchAPI/SerpApi (3rd party) |
