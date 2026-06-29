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

func TestGetVolumeByCenter(t *testing.T) {
	tests := []struct {
		name           string
		mockSetup      func(m *mocks.MockSaleService)
		expectedStatus int
		expectedBody   map[string]float64 // center -> volume for easy lookup
		expectError    bool
	}{
		{
			name: "returns volume by center successfully",
			mockSetup: func(m *mocks.MockSaleService) {
				m.EXPECT().
					GetVolumeByCenter(gomock.Any()).
					Return(map[string]float64{
						"C1": 45700.00,
						"C2": 38274.00,
						"C3": 37500.00,
						"C4": 55974.00,
					}, nil)
			},
			expectedStatus: http.StatusOK,
			expectedBody: map[string]float64{
				"C1": 45700.00,
				"C2": 38274.00,
				"C3": 37500.00,
				"C4": 55974.00,
			},
		},
		{
			name: "returns empty slice when no sales exist",
			mockSetup: func(m *mocks.MockSaleService) {
				m.EXPECT().
					GetVolumeByCenter(gomock.Any()).
					Return(map[string]float64{}, nil)
			},
			expectedStatus: http.StatusOK,
			expectedBody:   map[string]float64{},
		},
		{
			name: "returns 500 when service fails",
			mockSetup: func(m *mocks.MockSaleService) {
				m.EXPECT().
					GetVolumeByCenter(gomock.Any()).
					Return(nil, errors.New("internal error"))
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

			req := httptest.NewRequest(http.MethodGet, "/api/v1/sales/volume/centers", nil)
			w := httptest.NewRecorder()

			handler.GetVolumeByCenter(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, w.Code)
			}

			if tt.expectError {
				return
			}

			// unmarshal into slice then convert to map for order-independent comparison
			// since ranging over a map in the handler produces non-deterministic order
			var got []GetVolumeByCenterResponse
			if err := json.Unmarshal(w.Body.Bytes(), &got); err != nil {
				t.Fatalf("could not unmarshal response body: %v", err)
			}

			if len(got) != len(tt.expectedBody) {
				t.Errorf("expected %d centers, got %d", len(tt.expectedBody), len(got))
				return
			}

			for _, item := range got {
				expectedVolume, ok := tt.expectedBody[item.Center]
				if !ok {
					t.Errorf("unexpected center in response: %s", item.Center)
					continue
				}
				if item.Volume != expectedVolume {
					t.Errorf("center %s: expected volume %.2f, got %.2f", item.Center, expectedVolume, item.Volume)
				}
			}
		})
	}
}
