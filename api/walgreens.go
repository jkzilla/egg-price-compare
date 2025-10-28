package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/jkzilla/egg-price-compare/graph/model"
)

// WalgreensAPI handles Walgreens Store Inventory + Digital Offers APIs (1P retail)
// Note: Walgreens does not expose a general product pricing API
// Price data must come from a third-party provider (SearchAPI, SerpApi, Apify, etc.)
type WalgreensAPI struct {
	apiKey              string // Walgreens API Key
	apiSecret           string // Walgreens API Secret (OAuth)
	thirdPartyAPIKey    string // Third-party price provider API key (SearchAPI, SerpApi, etc.)
	client              *http.Client
}

// WalgreensInventoryResponse represents Store Inventory API response
type WalgreensInventoryResponse struct {
	StoreID     string `json:"storeId"`
	ProductID   string `json:"productId"`
	InStock     bool   `json:"inStock"`
	Quantity    int    `json:"quantity"`
	PickupReady bool   `json:"pickupReady"`
	PickupETA   string `json:"pickupEta"`
}

// WalgreensDigitalOffer represents a digital coupon/offer
type WalgreensDigitalOffer struct {
	OfferID         string  `json:"offerId"`
	Description     string  `json:"description"`
	DiscountAmount  float64 `json:"discountAmount"`
	DiscountPercent float64 `json:"discountPercent"`
	ExpiresAt       string  `json:"expiresAt"`
	Clippable       bool    `json:"clippable"`
}

// WalgreensOffersResponse represents Digital Offers API response
type WalgreensOffersResponse struct {
	Offers []WalgreensDigitalOffer `json:"offers"`
}

// ThirdPartyPriceResponse represents price data from third-party provider
// This could be SearchAPI, SerpApi, Apify, or custom scraper
type ThirdPartyPriceResponse struct {
	ProductName  string  `json:"product_name"`
	Price        float64 `json:"price"`
	RegularPrice float64 `json:"regular_price"`
	OnSale       bool    `json:"on_sale"`
	UPC          string  `json:"upc"`
	SKU          string  `json:"sku"`
	ProductURL   string  `json:"product_url"`
	Available    bool    `json:"available"`
}

func NewWalgreensAPI() *WalgreensAPI {
	return &WalgreensAPI{
		apiKey:           os.Getenv("WALGREENS_API_KEY"),
		apiSecret:        os.Getenv("WALGREENS_API_SECRET"),
		thirdPartyAPIKey: os.Getenv("SEARCHAPI_KEY"), // or SERPAPI_KEY, APIFY_KEY, etc.
		client:           &http.Client{Timeout: 15 * time.Second},
	}
}

// GetEggPrice fetches egg prices using:
// 1. Walgreens Store Inventory API for in-stock status
// 2. Walgreens Digital Offers API for clip-able coupons
// 3. Third-party data provider for actual pricing (SearchAPI, SerpApi, Apify, etc.)
func (w *WalgreensAPI) GetEggPrice(zipcode string) (*model.RetailerPrice, error) {
	if w.apiKey == "" && w.thirdPartyAPIKey == "" {
		// Return mock data for development
		// Vary prices slightly based on zipcode for more realistic testing
		zipcodeHash := 0
		for _, c := range zipcode {
			zipcodeHash += int(c)
		}
		
		// Base price varies between $3.99 and $5.29
		basePrice := 3.99 + float64(zipcodeHash%130)/100.0
		
		// Sometimes there's a sale price (70% of the time)
		var promoPrice *float64
		var offers []*model.DigitalOffer
		hasPromo := zipcodeHash%10 < 7
		
		finalPrice := basePrice
		if hasPromo {
			discount := 0.20 + float64(zipcodeHash%60)/100.0
			promo := basePrice - discount
			promoPrice = &promo
			finalPrice = promo
		}
		
		// Digital coupons (50% of the time)
		if zipcodeHash%10 < 5 {
			couponDiscount := 0.50
			offers = append(offers, &model.DigitalOffer{
				OfferID:        "WAG-DIGITAL-001",
				Description:    "Digital Coupon: Save $0.50",
				DiscountAmount: floatPtr(couponDiscount),
				ExpiresAt:      strPtr(time.Now().AddDate(0, 0, 7).Format(time.RFC3339)),
			})
			finalPrice -= couponDiscount
		}
		
		// Rewards program offer (20% of the time)
		if zipcodeHash%10 < 2 {
			offers = append(offers, &model.DigitalOffer{
				OfferID:         "WAG-REWARDS-002",
				Description:     "myWalgreens: 10% off",
				DiscountPercent: floatPtr(10.0),
			})
			finalPrice *= 0.90
		}
		
		// Stock status varies (85% in stock)
		inStock := zipcodeHash%20 < 17
		pickupEta := "Ready in 1 hour"
		if !inStock {
			pickupEta = "Out of stock at nearby stores"
		} else if zipcodeHash%5 == 0 {
			pickupEta = "Ready in 2-3 hours"
		}
		
		// Store ID varies
		storeID := fmt.Sprintf("%d", 10000+(zipcodeHash%5000))
		
		productURL := "https://www.walgreens.com/store/c/walgreens-grade-a-large-white-eggs/ID=prod6378461"
		
		return &model.RetailerPrice{
			Store:          "Walgreens",
			Sku:            strPtr("prod6378461"),
			Upc:            strPtr("041220993758"),
			StoreID:        strPtr(storeID),
			Zipcode:        zipcode,
			BasePrice:      basePrice,
			PromoPrice:     promoPrice,
			FinalPrice:     finalPrice,
			ProductName:    "Walgreens Grade A Large White Eggs, 12 ct",
			ProductURL:     &productURL,
			InStock:        inStock,
			PickupEta:      strPtr(pickupEta),
			DigitalOffers:  offers,
			LastUpdated:    time.Now().Format(time.RFC3339),
		}, nil
	}

	// Step 1: Get price data from third-party provider
	// This is necessary because Walgreens doesn't expose retail pricing API
	priceData, err := w.fetchPriceFromThirdParty(zipcode)
	if err != nil {
		return nil, fmt.Errorf("walgreens: failed to fetch price data: %w", err)
	}

	// Step 2: Get inventory status from Walgreens Store Inventory API
	inventory, err := w.fetchInventory(priceData.SKU, zipcode)
	if err != nil {
		// Log error but continue with price data
		fmt.Printf("walgreens: inventory check failed: %v\n", err)
		inventory = &WalgreensInventoryResponse{
			InStock:   priceData.Available,
			PickupETA: "Check store availability",
		}
	}

	// Step 3: Get digital offers from Walgreens Digital Offers API
	offers, err := w.fetchDigitalOffers(priceData.SKU)
	if err != nil {
		// Log error but continue without offers
		fmt.Printf("walgreens: digital offers fetch failed: %v\n", err)
		offers = []*model.DigitalOffer{}
	}

	// Calculate final price with digital offers
	basePrice := priceData.RegularPrice
	if basePrice == 0 {
		basePrice = priceData.Price
	}
	finalPrice := priceData.Price
	var promoPrice *float64
	
	if priceData.OnSale {
		promoPrice = &priceData.Price
	}
	
	// Apply digital offer discounts
	for _, offer := range offers {
		if offer.DiscountAmount != nil {
			finalPrice -= *offer.DiscountAmount
		}
	}

	return &model.RetailerPrice{
		Store:         "Walgreens",
		Sku:           &priceData.SKU,
		Upc:           &priceData.UPC,
		StoreID:       &inventory.StoreID,
		Zipcode:       zipcode,
		BasePrice:     basePrice,
		PromoPrice:    promoPrice,
		FinalPrice:    finalPrice,
		ProductName:   priceData.ProductName,
		ProductURL:    &priceData.ProductURL,
		InStock:       inventory.InStock,
		PickupEta:     &inventory.PickupETA,
		DigitalOffers: offers,
		LastUpdated:   time.Now().Format(time.RFC3339),
	}, nil
}

// fetchPriceFromThirdParty gets pricing data from SearchAPI, SerpApi, Apify, or similar
// These services are designed for price-comparison use cases and respect ToS
func (w *WalgreensAPI) fetchPriceFromThirdParty(zipcode string) (*ThirdPartyPriceResponse, error) {
	if w.thirdPartyAPIKey == "" {
		return nil, fmt.Errorf("third-party API key not configured")
	}

	// Example: SearchAPI for Walgreens product search
	// Documentation: https://www.searchapi.io/docs/walgreens
	baseURL := "https://www.searchapi.io/api/v1/search"
	
	params := url.Values{}
	params.Add("engine", "walgreens")
	params.Add("q", "eggs dozen large white")
	params.Add("api_key", w.thirdPartyAPIKey)
	params.Add("location", zipcode)
	
	requestURL := fmt.Sprintf("%s?%s", baseURL, params.Encode())
	
	req, err := http.NewRequest("GET", requestURL, nil)
	if err != nil {
		return nil, err
	}

	resp, err := w.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("third-party API returned status %d: %s", resp.StatusCode, string(body))
	}

	// Parse response (format depends on provider)
	// This is a simplified example - actual implementation depends on provider
	var result struct {
		Products []ThirdPartyPriceResponse `json:"products"`
	}
	
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, err
	}

	if len(result.Products) == 0 {
		return nil, fmt.Errorf("no products found")
	}

	return &result.Products[0], nil
}

// fetchInventory calls Walgreens Store Inventory API
func (w *WalgreensAPI) fetchInventory(sku, zipcode string) (*WalgreensInventoryResponse, error) {
	if w.apiKey == "" {
		return nil, fmt.Errorf("walgreens API key not configured")
	}

	// Walgreens Store Inventory API
	// Documentation: https://developer.walgreens.com/store-inventory
	url := fmt.Sprintf("https://services.walgreens.com/api/stores/inventory?sku=%s&zip=%s", sku, zipcode)
	
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	
	req.Header.Set("apikey", w.apiKey)
	// OAuth token would go here if required
	
	resp, err := w.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("inventory API returned status %d", resp.StatusCode)
	}

	var inventory WalgreensInventoryResponse
	if err := json.NewDecoder(resp.Body).Decode(&inventory); err != nil {
		return nil, err
	}

	return &inventory, nil
}

// fetchDigitalOffers calls Walgreens Digital Offers API
func (w *WalgreensAPI) fetchDigitalOffers(sku string) ([]*model.DigitalOffer, error) {
	if w.apiKey == "" {
		return []*model.DigitalOffer{}, nil
	}

	// Walgreens Digital Offers API
	// Documentation: https://developer.walgreens.com/digital-offers
	url := fmt.Sprintf("https://services.walgreens.com/api/offers?sku=%s", sku)
	
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	
	req.Header.Set("apikey", w.apiKey)
	
	resp, err := w.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return []*model.DigitalOffer{}, nil // No offers available
	}

	var offersResp WalgreensOffersResponse
	if err := json.NewDecoder(resp.Body).Decode(&offersResp); err != nil {
		return nil, err
	}

	// Convert to model format
	var offers []*model.DigitalOffer
	for _, offer := range offersResp.Offers {
		modelOffer := &model.DigitalOffer{
			OfferID:     offer.OfferID,
			Description: offer.Description,
			ExpiresAt:   &offer.ExpiresAt,
		}
		
		if offer.DiscountAmount > 0 {
			modelOffer.DiscountAmount = &offer.DiscountAmount
		}
		if offer.DiscountPercent > 0 {
			modelOffer.DiscountPercent = &offer.DiscountPercent
		}
		
		offers = append(offers, modelOffer)
	}

	return offers, nil
}

func strPtr(s string) *string {
	return &s
}

func floatPtr(f float64) *float64 {
	return &f
}
