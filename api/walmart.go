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

// WalmartAPI handles Walmart Affiliates Product Lookup API (1P retail pricing)
type WalmartAPI struct {
	affiliateID string // Walmart Affiliates Publisher ID
	apiKey      string // Walmart Affiliates API Key
	client      *http.Client
}

// WalmartAffiliateProduct represents a product from Walmart Affiliates Product Lookup API
type WalmartAffiliateProduct struct {
	ItemID              string  `json:"itemId"`
	Name                string  `json:"name"`
	SalePrice           float64 `json:"salePrice"`
	MSRP                float64 `json:"msrp"`
	UPC                 string  `json:"upc"`
	Stock               string  `json:"stock"`
	AvailableOnline     bool    `json:"availableOnline"`
	ProductURL          string  `json:"productUrl"`
	ShortDescription    string  `json:"shortDescription"`
	SpecialBuy          bool    `json:"specialBuy"`
	Clearance           bool    `json:"clearance"`
	PreOrder            bool    `json:"preOrder"`
	ShippingPassEligible bool   `json:"shippingPassEligible"`
}

// WalmartAffiliateResponse represents the response from Walmart Affiliates API
type WalmartAffiliateResponse struct {
	Items      []WalmartAffiliateProduct `json:"items"`
	TotalResults int                      `json:"totalResults"`
	Start        int                      `json:"start"`
	NumItems     int                      `json:"numItems"`
}

func NewWalmartAPI() *WalmartAPI {
	return &WalmartAPI{
		affiliateID: os.Getenv("WALMART_AFFILIATE_ID"),
		apiKey:      os.Getenv("WALMART_API_KEY"),
		client:      &http.Client{Timeout: 15 * time.Second},
	}
}

// GetEggPrice fetches egg prices using Walmart Affiliates Product Lookup API
// This is the official 1P retail pricing API for price-comparison use cases
func (w *WalmartAPI) GetEggPrice(zipcode string) (*model.RetailerPrice, error) {
	if w.apiKey == "" {
		// Return mock data for development
		// Vary prices slightly based on zipcode for more realistic testing
		zipcodeHash := 0
		for _, c := range zipcode {
			zipcodeHash += int(c)
		}
		
		// Base price varies between $3.48 and $4.98
		basePrice := 3.48 + float64(zipcodeHash%150)/100.0
		
		// Sometimes there's a promo (60% of the time)
		var promoPrice *float64
		var offers []*model.DigitalOffer
		hasPromo := zipcodeHash%10 < 6
		
		finalPrice := basePrice
		if hasPromo {
			discount := 0.30 + float64(zipcodeHash%70)/100.0
			promo := basePrice - discount
			promoPrice = &promo
			finalPrice = promo
			
			offers = append(offers, &model.DigitalOffer{
				OfferID:        "WMT-PROMO-001",
				Description:    "Rollback: Save on eggs",
				DiscountAmount: floatPtr(discount),
			})
		}
		
		// Occasionally add a digital coupon (30% of the time)
		if zipcodeHash%10 < 3 {
			offers = append(offers, &model.DigitalOffer{
				OfferID:        "WMT-DIGITAL-002",
				Description:    "Digital Coupon: Extra $0.25 off",
				DiscountAmount: floatPtr(0.25),
			})
			finalPrice -= 0.25
		}
		
		// Stock status varies (90% in stock)
		inStock := zipcodeHash%10 != 0
		pickupEta := "Available today"
		if !inStock {
			pickupEta = "Out of stock"
		}
		
		productURL := "https://www.walmart.com/ip/Great-Value-Large-White-Eggs-12-Count/10450114"
		
		return &model.RetailerPrice{
			Store:          "Walmart",
			Sku:            strPtr("10450114"),
			Upc:            strPtr("078742370842"),
			StoreID:        nil,
			Zipcode:        zipcode,
			BasePrice:      basePrice,
			PromoPrice:     promoPrice,
			FinalPrice:     finalPrice,
			ProductName:    "Great Value Large White Eggs, 12 Count",
			ProductURL:     &productURL,
			InStock:        inStock,
			PickupEta:      strPtr(pickupEta),
			DigitalOffers:  offers,
			LastUpdated:    time.Now().Format(time.RFC3339),
		}, nil
	}

	// Walmart Affiliates Product Lookup API
	// Documentation: https://developer.walmart.com/api/us/affil/product/v2
	baseURL := "https://developer.api.walmart.com/api-proxy/service/affil/product/v2/search"
	
	params := url.Values{}
	params.Add("query", "eggs dozen large white")
	params.Add("apiKey", w.apiKey)
	params.Add("format", "json")
	params.Add("numItems", "5") // Get top 5 results to find best match
	
	// Note: Walmart Affiliates API doesn't directly support zipcode filtering
	// For store-specific pricing, you may need to use SearchAPI or similar service
	// that provides zipcode-based Walmart pricing
	
	requestURL := fmt.Sprintf("%s?%s", baseURL, params.Encode())
	
	req, err := http.NewRequest("GET", requestURL, nil)
	if err != nil {
		return nil, fmt.Errorf("walmart: failed to create request: %w", err)
	}

	resp, err := w.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("walmart: request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("walmart: API returned status %d: %s", resp.StatusCode, string(body))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("walmart: failed to read response: %w", err)
	}

	var walmartResp WalmartAffiliateResponse
	if err := json.Unmarshal(body, &walmartResp); err != nil {
		return nil, fmt.Errorf("walmart: failed to parse response: %w", err)
	}

	if len(walmartResp.Items) == 0 {
		return nil, fmt.Errorf("walmart: no products found")
	}

	// Find the best match (first available product)
	product := walmartResp.Items[0]
	
	// Calculate final price (handle promos, clearance, special buy)
	basePrice := product.MSRP
	if basePrice == 0 {
		basePrice = product.SalePrice
	}
	finalPrice := product.SalePrice
	var promoPrice *float64
	
	if product.SalePrice < basePrice {
		promoPrice = &product.SalePrice
	}
	
	// Build digital offers list
	var offers []*model.DigitalOffer
	if product.Clearance {
		offers = append(offers, &model.DigitalOffer{
			OfferID:     fmt.Sprintf("WMT-CLR-%s", product.ItemID),
			Description: "Clearance Item",
		})
	}
	if product.SpecialBuy {
		offers = append(offers, &model.DigitalOffer{
			OfferID:     fmt.Sprintf("WMT-SPECIAL-%s", product.ItemID),
			Description: "Special Buy",
		})
	}
	
	inStock := product.Stock == "Available" && product.AvailableOnline
	
	return &model.RetailerPrice{
		Store:       "Walmart",
		Sku:         &product.ItemID,
		Upc:         &product.UPC,
		Zipcode:     zipcode,
		BasePrice:   basePrice,
		PromoPrice:  promoPrice,
		FinalPrice:  finalPrice,
		ProductName: product.Name,
		ProductURL:  &product.ProductURL,
		InStock:     inStock,
		PickupEta:   strPtr("Check store availability"),
		DigitalOffers: offers,
		LastUpdated: time.Now().Format(time.RFC3339),
	}, nil
}
