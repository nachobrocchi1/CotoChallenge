package service

import (
	"CotoChallenge/internal/domain"
	"CotoChallenge/internal/dto"
	"CotoChallenge/internal/repository/mocks"
	"context"
	"errors"
	"io"
	"log"
	"testing"

	"go.uber.org/mock/gomock"
)

var discardLogger = log.New(io.Discard, "", 0)

func newTestService(t *testing.T, mockRepo *mocks.MockSaleRepository) SaleService {
	t.Helper()

	s, err := NewSaleService(discardLogger, mockRepo)
	if err != nil {
		t.Fatalf("NewSaleService() error = %v", err)
	}

	return s
}

func TestSaleServiceImpl_CreateSale(t *testing.T) {
	type args struct {
		ctx     context.Context
		Vehicle domain.VehicleType
		Center  string
	}

	tests := []struct {
		name    string
		setup   func(mockRepo *mocks.MockSaleRepository)
		args    args
		wantErr bool
	}{
		{
			name: "Success - Persists sale",
			setup: func(mockRepo *mocks.MockSaleRepository) {
				mockRepo.EXPECT().
					CreateSale(gomock.Any(), gomock.Any()).
					Times(1)
			},
			args: args{
				ctx:     context.Background(),
				Center:  "Sale Center 1",
				Vehicle: domain.Sedan,
			},
			wantErr: false,
		},
		{
			name: "Success - Persists sale Sport vehicle",
			setup: func(mockRepo *mocks.MockSaleRepository) {
				mockRepo.EXPECT().
					CreateSale(gomock.Any(), gomock.Any()).
					Times(1)
			},
			args: args{
				ctx:     context.Background(),
				Center:  "Sale Center 1",
				Vehicle: domain.Sport,
			},
			wantErr: false,
		},
		{
			name: "Failure - error when database layer fails",
			setup: func(mockRepo *mocks.MockSaleRepository) {
				mockRepo.EXPECT().
					CreateSale(gomock.Any(), gomock.Any()).
					Times(1).
					Return(errors.New("write failure"))
			},
			args: args{
				ctx:     context.Background(),
				Center:  "Sale Center 1",
				Vehicle: domain.Sedan,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockRepo := mocks.NewMockSaleRepository(ctrl)

			tt.setup(mockRepo)

			s := newTestService(t, mockRepo)
			err := s.CreateSale(tt.args.ctx, tt.args.Vehicle, tt.args.Center)

			if (err != nil) != tt.wantErr {
				t.Errorf("CreateSale() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestSaleServiceImpl_GetTotalVolume(t *testing.T) {
	tests := []struct {
		name      string
		setup     func(mockRepo *mocks.MockSaleRepository)
		wantTotal float64
		wantErr   bool
	}{
		{
			name: "Success - sums all sale prices",
			setup: func(mockRepo *mocks.MockSaleRepository) {
				mockRepo.EXPECT().
					GetSales(gomock.Any()).
					Return([]domain.Sale{
						{Price: 8000, Center: "Center A"},
						{Price: 9500, Center: "Center A"},
						{Price: 12500, Center: "Center B"},
					}, nil)
			},
			wantTotal: 30000,
		},
		{
			name: "Success - returns zero when there are no sales",
			setup: func(mockRepo *mocks.MockSaleRepository) {
				mockRepo.EXPECT().
					GetSales(gomock.Any()).
					Return([]domain.Sale{}, nil)
			},
			wantTotal: 0,
		},
		{
			name: "Failure - repository error",
			setup: func(mockRepo *mocks.MockSaleRepository) {
				mockRepo.EXPECT().
					GetSales(gomock.Any()).
					Return(nil, errors.New("read failure"))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockRepo := mocks.NewMockSaleRepository(ctrl)
			tt.setup(mockRepo)

			s := newTestService(t, mockRepo)
			got, err := s.GetTotalVolume(context.Background())

			if (err != nil) != tt.wantErr {
				t.Fatalf("GetTotalVolume() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr && got != tt.wantTotal {
				t.Errorf("GetTotalVolume() = %v, want %v", got, tt.wantTotal)
			}
		})
	}
}

func TestSaleServiceImpl_GetVolumeByCenter(t *testing.T) {
	tests := []struct {
		name    string
		setup   func(mockRepo *mocks.MockSaleRepository)
		want    map[string]float64
		wantErr bool
	}{
		{
			name: "Success - aggregates volume by center",
			setup: func(mockRepo *mocks.MockSaleRepository) {
				mockRepo.EXPECT().
					GetSales(gomock.Any()).
					Return([]domain.Sale{
						{Price: 8000, Center: "Center A"},
						{Price: 9500, Center: "Center A"},
						{Price: 12500, Center: "Center B"},
					}, nil)
			},
			want: map[string]float64{
				"Center A": 17500,
				"Center B": 12500,
			},
		},
		{
			name: "Success - returns empty map when there are no sales",
			setup: func(mockRepo *mocks.MockSaleRepository) {
				mockRepo.EXPECT().
					GetSales(gomock.Any()).
					Return([]domain.Sale{}, nil)
			},
			want: map[string]float64{},
		},
		{
			name: "Failure - repository error",
			setup: func(mockRepo *mocks.MockSaleRepository) {
				mockRepo.EXPECT().
					GetSales(gomock.Any()).
					Return(nil, errors.New("read failure"))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockRepo := mocks.NewMockSaleRepository(ctrl)
			tt.setup(mockRepo)

			s := newTestService(t, mockRepo)
			got, err := s.GetVolumeByCenter(context.Background())

			if (err != nil) != tt.wantErr {
				t.Fatalf("GetVolumeByCenter() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr {
				return
			}

			if len(got) != len(tt.want) {
				t.Fatalf("GetVolumeByCenter() map length = %d, want %d", len(got), len(tt.want))
			}
			for center, wantVolume := range tt.want {
				if got[center] != wantVolume {
					t.Errorf("GetVolumeByCenter()[%q] = %v, want %v", center, got[center], wantVolume)
				}
			}
		})
	}
}

func TestSaleServiceImpl_GetSalesPercentegeByCenterOverTotalSales(t *testing.T) {
	tests := []struct {
		name    string
		setup   func(mockRepo *mocks.MockSaleRepository)
		want    map[string]map[domain.VehicleType]float64
		wantErr bool
	}{
		{
			name: "Success - calculates percentages grouped by center and model",
			setup: func(mockRepo *mocks.MockSaleRepository) {
				mockRepo.EXPECT().
					GetSales(gomock.Any()).
					Return([]domain.Sale{
						{Center: "Center A", Vehicle: domain.Sedan},
						{Center: "Center A", Vehicle: domain.Sedan},
						{Center: "Center A", Vehicle: domain.SUV},
						{Center: "Center B", Vehicle: domain.Sedan},
					}, nil)
			},
			want: map[string]map[domain.VehicleType]float64{
				"Center A": {
					domain.Sedan: 50,
					domain.SUV:   25,
				},
				"Center B": {
					domain.Sedan: 25,
				},
			},
		},
		{
			name: "Success - rounds percentages to two decimals",
			setup: func(mockRepo *mocks.MockSaleRepository) {
				mockRepo.EXPECT().
					GetSales(gomock.Any()).
					Return([]domain.Sale{
						{Center: "Center A", Vehicle: domain.Sedan},
						{Center: "Center A", Vehicle: domain.Sedan},
						{Center: "Center B", Vehicle: domain.SUV},
					}, nil)
			},
			want: map[string]map[domain.VehicleType]float64{
				"Center A": {
					domain.Sedan: 66.67,
				},
				"Center B": {
					domain.SUV: 33.33,
				},
			},
		},
		{
			name: "Success - returns empty slice when there are no sales",
			setup: func(mockRepo *mocks.MockSaleRepository) {
				mockRepo.EXPECT().
					GetSales(gomock.Any()).
					Return([]domain.Sale{}, nil)
			},
			want: map[string]map[domain.VehicleType]float64{},
		},
		{
			name: "Failure - repository error",
			setup: func(mockRepo *mocks.MockSaleRepository) {
				mockRepo.EXPECT().
					GetSales(gomock.Any()).
					Return(nil, errors.New("read failure"))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockRepo := mocks.NewMockSaleRepository(ctrl)
			tt.setup(mockRepo)

			s := newTestService(t, mockRepo)
			got, err := s.GetSalesPercentageByCenterOverTotalSales(context.Background())

			if (err != nil) != tt.wantErr {
				t.Fatalf("GetSalesPercentegeByCenterOverTotalSales() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr {
				return
			}

			assertCenterModelPercentages(t, got, tt.want)
		})
	}
}

func assertCenterModelPercentages(
	t *testing.T,
	got []dto.CenterModelPercentage,
	want map[string]map[domain.VehicleType]float64,
) {
	t.Helper()

	if len(got) != len(want) {
		t.Fatalf("result length = %d, want %d", len(got), len(want))
	}

	for _, centerResult := range got {
		wantModels, ok := want[centerResult.Center]
		if !ok {
			t.Fatalf("unexpected center %q in result", centerResult.Center)
		}

		if len(centerResult.Models) != len(wantModels) {
			t.Fatalf("center %q models length = %d, want %d", centerResult.Center, len(centerResult.Models), len(wantModels))
		}

		for _, modelResult := range centerResult.Models {
			wantPercentage, ok := wantModels[modelResult.Vehicle]
			if !ok {
				t.Fatalf("unexpected vehicle %q for center %q", modelResult.Vehicle, centerResult.Center)
			}
			if modelResult.Percentage != wantPercentage {
				t.Errorf("center %q vehicle %q percentage = %v, want %v", centerResult.Center, modelResult.Vehicle, modelResult.Percentage, wantPercentage)
			}
		}
	}
}
