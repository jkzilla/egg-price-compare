package graph

import (
	"context"
	"fmt"
	"math"
	"time"

	"github.com/jkzilla/egg-price-compare/api"
	"github.com/jkzilla/egg-price-compare/graph/model"
)

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

func (r *queryResolver) EggPrices(ctx context.Context) (*model.EggPriceComparison, error) {
	walmartPrice, err := r.walmartAPI.GetEggPrice()
	if err != nil {
		return nil, fmt.Errorf("failed to get Walmart price: %w", err)
	}

	walgreensPrice, err := r.walgreensAPI.GetEggPrice()
	if err != nil {
		return nil, fmt.Errorf("failed to get Walgreens price: %w", err)
	}

	cheapest := "Walmart"
	if walgreensPrice.Price < walmartPrice.Price {
		cheapest = "Walgreens"
	}

	priceDiff := math.Abs(walmartPrice.Price - walgreensPrice.Price)

	// Store in history
	historyEntry := model.PriceHistoryEntry{
		Date:           time.Now().Format("2006-01-02"),
		WalmartPrice:   walmartPrice.Price,
		WalgreensPrice: walgreensPrice.Price,
	}
	
	// Add to history if it's a new day
	if len(r.priceHistory) == 0 || r.priceHistory[len(r.priceHistory)-1].Date != historyEntry.Date {
		r.priceHistory = append(r.priceHistory, historyEntry)
	}

	return &model.EggPriceComparison{
		Walmart:         walmartPrice,
		Walgreens:       walgreensPrice,
		Cheapest:        cheapest,
		PriceDifference: priceDiff,
		LastUpdated:     time.Now().Format(time.RFC3339),
	}, nil
}

func (r *queryResolver) PriceHistory(ctx context.Context, days *int) ([]*model.PriceHistoryEntry, error) {
	numDays := 7
	if days != nil {
		numDays = *days
	}

	// Return last N days
	start := len(r.priceHistory) - numDays
	if start < 0 {
		start = 0
	}

	history := make([]*model.PriceHistoryEntry, 0)
	for i := start; i < len(r.priceHistory); i++ {
		entry := r.priceHistory[i]
		history = append(history, &entry)
	}

	return history, nil
}

type queryResolver struct{ *Resolver }

func (r *Resolver) Query() QueryResolver { return &queryResolver{r} }
