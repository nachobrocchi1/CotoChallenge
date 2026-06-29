package handler

import (
	"CotoChallenge/internal/service"
)

type SaleHandler struct {
	service service.SaleService
}

func NewSaleHandler(service service.SaleService) (*SaleHandler, error) {
	return &SaleHandler{
		service: service,
	}, nil
}
