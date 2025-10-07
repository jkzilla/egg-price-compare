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

type WalgreensAPI struct {
	apiKey string
	client *http.Client
}

type WalgreensProduct struct {
	ProductInfo struct {
		ProductName string  `json:"productName"`
		RegularPrice float64 `json:"regularPrice"`
		ProductID   string  `json:"productId"`
	} `json:"productInfo"`
	Inventory struct {
		InStock bool `json:"inStock"`
	} `json:"inventory"`
}

type WalgreensResponse struct {
	Products []WalgreensProduct `json:"products"`
}

func NewWalgreensAPI() *WalgreensAPI {
	return &WalgreensAPI{
		apiKey: os.Getenv("WALGREENS_API_KEY"),
		client: &http.Client{Timeout: 10 * time.Second},
	}
}

func (w *WalgreensAPI) GetEggPrice() (*model.StorePrice, error) {
	// For demo purposes, return mock data
	// In production, you would call the actual Walgreens API
	
	if w.apiKey == "" {
		// Return mock data for development
		productURL := "https://www.walgreens.com/store/eggs"
		return &model.StorePrice{
			Store:       "Walgreens",
			Price:       4.29,
			ProductName: "Walgreens Grade A Large White Eggs, 12 ct",
			ProductURL:  &productURL,
			InStock:     true,
			LastUpdated: time.Now().Format(time.RFC3339),
		}, nil
	}

	// Real API call would go here
	url := fmt.Sprintf("https://api.walgreens.com/v1/products/search?q=eggs&apiKey=%s", w.apiKey)
	
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

	var walgreensResp WalgreensResponse
	if err := json.Unmarshal(body, &walgreensResp); err != nil {
		return nil, err
	}

	if len(walgreensResp.Products) == 0 {
		return nil, fmt.Errorf("no products found")
	}

	product := walgreensResp.Products[0]
	productURL := fmt.Sprintf("https://www.walgreens.com/store/c/product/%s", product.ProductInfo.ProductID)
	
	return &model.StorePrice{
		Store:       "Walgreens",
		Price:       product.ProductInfo.RegularPrice,
		ProductName: product.ProductInfo.ProductName,
		ProductURL:  &productURL,
		InStock:     product.Inventory.InStock,
		LastUpdated: time.Now().Format(time.RFC3339),
	}, nil
}
