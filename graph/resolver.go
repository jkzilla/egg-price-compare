package graph

import (
	"github.com/jkzilla/egg-price-compare/api"
	"github.com/jkzilla/egg-price-compare/graph/model"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	walmartAPI   *api.WalmartAPI
	walgreensAPI *api.WalgreensAPI
	priceHistory []model.PriceHistoryEntry
}

func NewResolver() *Resolver {
	return &Resolver{
		walmartAPI:   api.NewWalmartAPI(),
		walgreensAPI: api.NewWalgreensAPI(),
		priceHistory: []model.PriceHistoryEntry{},
	}
}
