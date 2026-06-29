package service

import (
	"CotoChallenge/internal/domain"
	"CotoChallenge/internal/dto"
	"CotoChallenge/internal/repository"
	"context"
	"errors"
	"log"
	"math"
	"time"
)

// SaleService defines the business logic operations
type SaleService interface {
	CreateSale(ctx context.Context, vehicle domain.VehicleType, center string) error
	GetTotalVolume(ctx context.Context) (float64, error)
	GetVolumeByCenter(ctx context.Context) (map[string]float64, error)
	GetSalesPercentegeByCenterOverTotalSales(ctx context.Context) ([]dto.CenterModelPercentage, error)
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

	err := s.repository.CreateSale(ctx, sale)
	if err != nil {
		s.logger.Printf("CreateSale - error calling repository.CreateSale %s", err.Error())
		return errors.New("error creating sale")
	}
	return nil
}

func (s *SaleServiceImpl) GetTotalVolume(ctx context.Context) (float64, error) {
	sales, err := s.repository.GetSales(ctx)
	if err != nil {
		s.logger.Printf("GetTotalVolume - error calling repository.GetSales %s", err.Error())
		return 0, errors.New("error getting total volume")
	}
	var total float64
	for _, s := range sales {
		total += s.Price
	}
	return total, err
}

func (s *SaleServiceImpl) GetVolumeByCenter(ctx context.Context) (map[string]float64, error) {
	sales, err := s.repository.GetSales(ctx)
	if err != nil {
		s.logger.Printf("GetVolumeByCenter - error calling repository.GetSales %s", err.Error())
		return nil, errors.New("error getting volume by center")
	}
	volume := make(map[string]float64)
	for _, sale := range sales {
		volume[sale.Center] += sale.Price
	}

	return volume, nil
}

func (s *SaleServiceImpl) GetSalesPercentegeByCenterOverTotalSales(ctx context.Context) ([]dto.CenterModelPercentage, error) {
	sales, err := s.repository.GetSales(ctx)
	if err != nil {
		s.logger.Printf("GetSalesPercentegeByCenterOverTotalSales - error calling repository.GetSales %s", err.Error())
		return nil, errors.New("error getting pertentage by center over total sales")
	}

	totalUnits := len(sales)
	if totalUnits == 0 {
		return []dto.CenterModelPercentage{}, nil
	}
	// group sales by center and vehicle
	grouped := make(map[string]map[domain.VehicleType]int)
	for _, sale := range sales {
		if _, ok := grouped[sale.Center]; !ok {
			grouped[sale.Center] = make(map[domain.VehicleType]int)
		}
		grouped[sale.Center][sale.Vehicle]++
	}

	result := make([]dto.CenterModelPercentage, 0, len(grouped))
	// iterate over each center
	for center, models := range grouped {
		modelPercentages := make([]dto.ModelPercentage, 0, len(models))

		// iterate over each vehicle model sold in this center
		for vehicle, count := range models {

			// percentage calculation for the model over the total units sold
			percentage := (float64(count) / float64(totalUnits)) * 100

			modelPercentages = append(modelPercentages, dto.ModelPercentage{
				Vehicle:    vehicle,
				Percentage: math.Round(percentage*100) / 100, // rounded to 2 decimales
			})
		}
		result = append(result, dto.CenterModelPercentage{
			Center: center,
			Models: modelPercentages,
		})
	}

	return result, nil
}
