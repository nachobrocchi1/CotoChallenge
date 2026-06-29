package service

import (
	"CotoChallenge/internal/domain"
	"CotoChallenge/internal/repository"
	"context"
	"log"
	"time"
)

// SaleService defines the business logic operations
type SaleService interface {
	CreateSale(ctx context.Context, vehicle domain.VehicleType, center string) error
	GetSales(ctx context.Context) ([]domain.Sale, error)
	GetSalesByCenter(ctx context.Context, center string) ([]domain.Sale, error)
	GetSalesPercentegeByCenterOverTotalSales(ctx context.Context) (float64, error)
}

type SaleServiceImpl struct {
	logger     *log.Logger
	repository repository.SaleRepository
}

func NewSaleService(logger *log.Logger, repository repository.SaleRepository) (SaleService, error) {
	return &SaleServiceImpl{
		logger:     logger,
		repository: repository,
	}, nil
}

func (s *SaleServiceImpl) CreateSale(ctx context.Context, vehicle domain.VehicleType, center string) error {
	sale := domain.Sale{
		Vehicle: vehicle,
		Price:   vehicle.FinalPrice(),
		Date:    time.Now(),
		Center:  center,
	}

	return s.repository.CreateSale(ctx, sale)
}

func (s *SaleServiceImpl) GetSales(ctx context.Context) ([]domain.Sale, error) {
	return s.repository.GetSales(ctx)
}

func (s *SaleServiceImpl) GetSalesByCenter(ctx context.Context, center string) ([]domain.Sale, error) {
	return s.repository.GetSalesByCenter(ctx, center)
}

func (s *SaleServiceImpl) GetSalesPercentegeByCenterOverTotalSales(ctx context.Context) (float64, error) {
	return s.repository.GetSalesPercentegeByCenterOverTotalSales(ctx)
}
