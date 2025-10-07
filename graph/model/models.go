package model

type EggPriceComparison struct {
	Walmart         *StorePrice `json:"walmart"`
	Walgreens       *StorePrice `json:"walgreens"`
	Cheapest        string      `json:"cheapest"`
	PriceDifference float64     `json:"priceDifference"`
	LastUpdated     string      `json:"lastUpdated"`
}

type StorePrice struct {
	Store       string  `json:"store"`
	Price       float64 `json:"price"`
	ProductName string  `json:"productName"`
	ProductURL  *string `json:"productUrl,omitempty"`
	InStock     bool    `json:"inStock"`
	LastUpdated string  `json:"lastUpdated"`
}

type PriceHistoryEntry struct {
	Date           string  `json:"date"`
	WalmartPrice   float64 `json:"walmartPrice"`
	WalgreensPrice float64 `json:"walgreensPrice"`
}
