package model

type EggPriceComparison struct {
	Walmart         *RetailerPrice `json:"walmart"`
	Walgreens       *RetailerPrice `json:"walgreens"`
	Cheapest        string         `json:"cheapest"`
	PriceDifference float64        `json:"priceDifference"`
	LastUpdated     string         `json:"lastUpdated"`
}

type RetailerPrice struct {
	Store         string          `json:"store"`
	Sku           *string         `json:"sku,omitempty"`
	Upc           *string         `json:"upc,omitempty"`
	StoreID       *string         `json:"storeId,omitempty"`
	Zipcode       string          `json:"zipcode"`
	BasePrice     float64         `json:"basePrice"`
	PromoPrice    *float64        `json:"promoPrice,omitempty"`
	FinalPrice    float64         `json:"finalPrice"`
	ProductName   string          `json:"productName"`
	ProductURL    *string         `json:"productUrl,omitempty"`
	InStock       bool            `json:"inStock"`
	PickupEta     *string         `json:"pickupEta,omitempty"`
	DigitalOffers []*DigitalOffer `json:"digitalOffers,omitempty"`
	LastUpdated   string          `json:"lastUpdated"`
}

type DigitalOffer struct {
	OfferID         string   `json:"offerId"`
	Description     string   `json:"description"`
	DiscountAmount  *float64 `json:"discountAmount,omitempty"`
	DiscountPercent *float64 `json:"discountPercent,omitempty"`
	ExpiresAt       *string  `json:"expiresAt,omitempty"`
}

// Legacy type - kept for backward compatibility
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
