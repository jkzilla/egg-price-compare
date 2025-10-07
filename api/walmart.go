package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/jkzilla/egg-price-compare/graph/model"
)

type WalmartAPI struct {
	apiKey string
	client *http.Client
}

type WalmartProduct struct {
	Name      string  `json:"name"`
	SalePrice float64 `json:"salePrice"`
	Stock     string  `json:"stock"`
	ItemID    string  `json:"itemId"`
}

type WalmartResponse struct {
	Items []WalmartProduct `json:"items"`
}

func NewWalmartAPI() *WalmartAPI {
	return &WalmartAPI{
		apiKey: os.Getenv("WALMART_API_KEY"),
		client: &http.Client{Timeout: 10 * time.Second},
	}
}

func (w *WalmartAPI) GetEggPrice() (*model.StorePrice, error) {
	// For demo purposes, return mock data
	// In production, you would call the actual Walmart API
	// Example: https://developer.walmart.com/api/us/mp/items
	
	if w.apiKey == "" {
		// Return mock data for development
		productURL := "https://www.walmart.com/ip/eggs"
		return &model.StorePrice{
			Store:       "Walmart",
			Price:       3.98,
			ProductName: "Great Value Large White Eggs, 12 Count",
			ProductURL:  &productURL,
			InStock:     true,
			LastUpdated: time.Now().Format(time.RFC3339),
		}, nil
	}

	// Real API call would go here
	url := fmt.Sprintf("https://api.walmart.com/v1/search?query=eggs+dozen&apiKey=%s", w.apiKey)
	
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := w.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var walmartResp WalmartResponse
	if err := json.Unmarshal(body, &walmartResp); err != nil {
		return nil, err
	}

	if len(walmartResp.Items) == 0 {
		return nil, fmt.Errorf("no products found")
	}

	product := walmartResp.Items[0]
	productURL := fmt.Sprintf("https://www.walmart.com/ip/%s", product.ItemID)
	
	return &model.StorePrice{
		Store:       "Walmart",
		Price:       product.SalePrice,
		ProductName: product.Name,
		ProductURL:  &productURL,
		InStock:     product.Stock == "Available",
		LastUpdated: time.Now().Format(time.RFC3339),
	}, nil
}
