package handler

import (
	"CotoChallenge/internal/domain"
	"CotoChallenge/internal/service/mocks"
	"bytes"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"go.uber.org/mock/gomock"
)

func TestSaleHandler_CreateSale(t *testing.T) {
	type args struct {
		body []byte
	}
	tests := []struct {
		name       string
		args       args
		setupMock  func(mockService *mocks.MockSaleService)
		wantStatus int
		wantInBody string
	}{
		{
			name: "Success - Valid payload routes correctly to service layer",
			args: args{
				body: []byte(`{"vehicle":"sedan","center":"Buenos Aires"}`),
			},
			setupMock: func(mockService *mocks.MockSaleService) {
				mockService.EXPECT().
					CreateSale(gomock.Any(), domain.VehicleType("sedan"), "Buenos Aires").
					Times(1).
					Return(nil)
			},
			wantStatus: http.StatusCreated,
			wantInBody: "", // EncodeNoContent so no response body
		},
		{
			name: "Failure - Malformed JSON triggers Bad Request",
			args: args{
				body: []byte(`{"vehicle":"sedan",`),
			},
			setupMock:  func(mockService *mocks.MockSaleService) {}, // service shouldn't be called
			wantStatus: http.StatusBadRequest,
			wantInBody: "unexpected EOF",
		},
		{
			name: "Failure - DTO Validation error returns structured message",
			args: args{
				body: []byte(`{"vehicle":"invalid_type","center":""}`),
			},
			setupMock:  func(mockService *mocks.MockSaleService) {},
			wantStatus: http.StatusBadRequest,
			wantInBody: "validation failed:",
		},
		{
			name: "Failure - Service layer infrastructure failure bubbles up as 500",
			args: args{
				body: []byte(`{"vehicle":"sport","center":"Rosario"}`),
			},
			setupMock: func(mockService *mocks.MockSaleService) {
				mockService.EXPECT().
					CreateSale(gomock.Any(), domain.VehicleType("sport"), "Rosario").
					Times(1).
					Return(errors.New("internal cluster error"))
			},
			wantStatus: http.StatusInternalServerError,
			wantInBody: "internal cluster error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockService := mocks.NewMockSaleService(ctrl)
			tt.setupMock(mockService)

			h := &SaleHandler{
				service: mockService,
			}

			req := httptest.NewRequest(http.MethodPost, "/api/v1/sales", bytes.NewBuffer(tt.args.body))
			rec := httptest.NewRecorder()

			h.CreateSale(rec, req)

			if rec.Code != tt.wantStatus {
				t.Errorf("CreateSale() status = %d, want %d", rec.Code, tt.wantStatus)
			}

			if tt.wantInBody != "" {
				bodyStr := rec.Body.String()
				if !strings.Contains(bodyStr, tt.wantInBody) {
					t.Errorf("CreateSale() body = %q, expected to contain %q", bodyStr, tt.wantInBody)
				}
			}
		})
	}
}
