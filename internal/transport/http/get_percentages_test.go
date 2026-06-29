package handler

import (
	"CotoChallenge/internal/domain"
	"CotoChallenge/internal/dto"
	"CotoChallenge/internal/service/mocks"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"go.uber.org/mock/gomock"
)

func TestGetSalesPercentegeByCenterOverTotalSales(t *testing.T) {
	tests := []struct {
		name           string
		mockSetup      func(m *mocks.MockSaleService)
		expectedStatus int
		expectedBody   []dto.CenterModelPercentage
	}{
		{
			name: "returns percentages successfully",
			mockSetup: func(m *mocks.MockSaleService) {
				m.EXPECT().
					GetSalesPercentageByCenterOverTotalSales(gomock.Any()).
					Return([]dto.CenterModelPercentage{
						{
							Center: "C1",
							Models: []dto.ModelPercentage{
								{Vehicle: domain.Sedan, Percentage: 20.00},
								{Vehicle: domain.SUV, Percentage: 10.00},
							},
						},
						{
							Center: "C2",
							Models: []dto.ModelPercentage{
								{Vehicle: domain.Sport, Percentage: 15.00},
							},
						},
					}, nil)
			},
			expectedStatus: http.StatusOK,
			expectedBody: []dto.CenterModelPercentage{
				{
					Center: "C1",
					Models: []dto.ModelPercentage{
						{Vehicle: domain.Sedan, Percentage: 20.00},
						{Vehicle: domain.SUV, Percentage: 10.00},
					},
				},
				{
					Center: "C2",
					Models: []dto.ModelPercentage{
						{Vehicle: domain.Sport, Percentage: 15.00},
					},
				},
			},
		},
		{
			name: "returns empty slice when no sales exist",
			mockSetup: func(m *mocks.MockSaleService) {
				m.EXPECT().
					GetSalesPercentageByCenterOverTotalSales(gomock.Any()).
					Return([]dto.CenterModelPercentage{}, nil)
			},
			expectedStatus: http.StatusOK,
			expectedBody:   []dto.CenterModelPercentage{},
		},
		{
			name: "returns 500 when service fails",
			mockSetup: func(m *mocks.MockSaleService) {
				m.EXPECT().
					GetSalesPercentageByCenterOverTotalSales(gomock.Any()).
					Return(nil, errors.New("internal error"))
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockService := mocks.NewMockSaleService(ctrl)
			tt.mockSetup(mockService)

			handler, err := NewSaleHandler(mockService)
			if err != nil {
				t.Errorf("error creating handler: %v", err)
			}
			req := httptest.NewRequest(http.MethodGet, "/api/v1/sales/percentage/centers", nil)
			w := httptest.NewRecorder()

			handler.GetSalesPercentegeByCenterOverTotalSales(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, w.Code)
			}

			// skip body assertion on error responses
			if tt.expectedBody == nil {
				return
			}

			var got []dto.CenterModelPercentage
			if err := json.Unmarshal(w.Body.Bytes(), &got); err != nil {
				t.Fatalf("could not unmarshal response body: %v", err)
			}

			if len(got) != len(tt.expectedBody) {
				t.Errorf("expected %d centers, got %d", len(tt.expectedBody), len(got))
				return
			}

			for i, center := range got {
				if center.Center != tt.expectedBody[i].Center {
					t.Errorf("expected center %s, got %s", tt.expectedBody[i].Center, center.Center)
				}
				if len(center.Models) != len(tt.expectedBody[i].Models) {
					t.Errorf("center %s: expected %d models, got %d", center.Center, len(tt.expectedBody[i].Models), len(center.Models))
				}
			}
		})
	}
}
