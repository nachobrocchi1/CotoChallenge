package repository

import (
	"CotoChallenge/internal/domain"
	"context"
)

// SaleRepository defines the data store contract
type SaleRepository interface {
	CreateSale(ctx context.Context, sale domain.Sale) error
	GetSales(ctx context.Context) ([]domain.Sale, error)
}
