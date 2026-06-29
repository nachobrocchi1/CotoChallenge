package service

import (
	"CotoChallenge/internal/domain"
	"CotoChallenge/internal/repository/mocks"
	"context"
	"errors"
	"io"
	"log"
	"testing"

	"go.uber.org/mock/gomock"
)

func TestSaleServiceImpl_CreateSale(t *testing.T) {
	discardLogger := log.New(io.Discard, "", 0)

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

			s, _ := NewSaleService(discardLogger, mockRepo)
			err := s.CreateSale(tt.args.ctx, tt.args.Vehicle, tt.args.Center)

			if (err != nil) != tt.wantErr {
				t.Errorf("CreateSale() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
