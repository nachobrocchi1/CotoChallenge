package handler

import (
	"CotoChallenge/internal/service/mocks"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"go.uber.org/mock/gomock"
)

func TestGetTotalVolume(t *testing.T) {
	tests := []struct {
		name           string
		mockSetup      func(m *mocks.MockSaleService)
		expectedStatus int
		expectedVolume float64
		expectError    bool
	}{
		{
			name: "returns total volume successfully",
			mockSetup: func(m *mocks.MockSaleService) {
				m.EXPECT().
					GetTotalVolume(gomock.Any()).
					Return(150000.00, nil)
			},
			expectedStatus: http.StatusOK,
			expectedVolume: 150000.00,
		},
		{
			name: "returns zero when no sales exist",
			mockSetup: func(m *mocks.MockSaleService) {
				m.EXPECT().
					GetTotalVolume(gomock.Any()).
					Return(0.0, nil)
			},
			expectedStatus: http.StatusOK,
			expectedVolume: 0.0,
		},
		{
			name: "returns 500 when service fails",
			mockSetup: func(m *mocks.MockSaleService) {
				m.EXPECT().
					GetTotalVolume(gomock.Any()).
					Return(0.0, errors.New("internal error"))
			},
			expectedStatus: http.StatusInternalServerError,
			expectError:    true,
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
			req := httptest.NewRequest(http.MethodGet, "/api/v1/sales/volume", nil)
			w := httptest.NewRecorder()

			handler.GetTotalVolume(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, w.Code)
			}

			if tt.expectError {
				return
			}

			var got GetTotalVolumeResponse
			if err := json.Unmarshal(w.Body.Bytes(), &got); err != nil {
				t.Fatalf("could not unmarshal response body: %v", err)
			}

			if got.Volume != tt.expectedVolume {
				t.Errorf("expected volume %.2f, got %.2f", tt.expectedVolume, got.Volume)
			}
		})
	}
}
