package handler

import (
	"CotoChallenge/internal/domain"
	"net/http"
)

const maxReadBytes = 1024 * 1024 // 1 Mb

type CreateSaleRequest struct {
	Vehicle domain.VehicleType `json:"vehicle"`
	Center  string             `json:"center"`
}

// CreateSale creates a new sale
func (h *SaleHandler) CreateSale(w http.ResponseWriter, r *http.Request) {
	request, err := Decode[CreateSaleRequest](r.Body, maxReadBytes)
	if err != nil {
		EncodeError(w, http.StatusBadRequest, err)
		return
	}

	if err := request.Validate(); err != nil {
		EncodeError(w, http.StatusBadRequest, err)
		return
	}

	if err := h.service.CreateSale(r.Context(), request.Vehicle, request.Center); err != nil {
		EncodeError(w, http.StatusInternalServerError, err)
		return
	}

	EncodeNoContent(w, http.StatusCreated)
}

func (r CreateSaleRequest) Validate() error {
	ve := ValidationError{Fields: make(map[string]string)}
	if _, ok := domain.Prices[r.Vehicle]; !ok {
		ve.Fields["vehicle"] = "must be one of: sedan, suv, offroad, sport"
	}
	if r.Center == "" {
		ve.Fields["center"] = "required"
	}
	if len(ve.Fields) > 0 {
		return ve
	}
	return nil
}
