package repository

import (
	"CotoChallenge/internal/domain"
	"context"
	"sync"

	"github.com/google/uuid"
)

type InMemorySaleRepository struct {
	sales []domain.Sale
	mu    sync.RWMutex
}

func NewInMemorySaleRepository() (SaleRepository, error) {
	return &InMemorySaleRepository{
		sales: make([]domain.Sale, 0),
	}, nil
}

func (r *InMemorySaleRepository) CreateSale(ctx context.Context, sale domain.Sale) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	sale.Id = uuid.NewString()
	r.sales = append(r.sales, sale)
	return nil
}

func (r *InMemorySaleRepository) GetSales(ctx context.Context) ([]domain.Sale, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	if len(r.sales) == 0 {
		return nil, ErrNoRows{}
	}
	// clone slice to return a copy to avoid unwanted updates on the original slice.
	return cloneSales(r.sales), nil
}

func cloneSales(sales []domain.Sale) []domain.Sale {
	clone := make([]domain.Sale, len(sales))
	copy(clone, sales)
	return clone
}
