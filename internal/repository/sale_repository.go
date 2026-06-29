package repository

import (
	"CotoChallenge/internal/domain"
	"context"
)

// SaleRespository defines the data store contract
type SaleRepository interface {
	CreateSale(ctx context.Context, sale domain.Sale) error
	GetSales(ctx context.Context) ([]domain.Sale, error)
	GetSalesByCenter(cctx context.Context, enter string) ([]domain.Sale, error)
	GetSalesPercentegeByCenterOverTotalSales(ctx context.Context) (float64, error)
}
