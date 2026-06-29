package handler

import (
	"CotoChallenge/internal/service"
	"fmt"
	"net/http"
)

type SaleHandler struct {
	service service.SaleService
}

func NewSaleHandler(service service.SaleService) (*SaleHandler, error) {
	return &SaleHandler{
		service: service,
	}, nil
}

// GetSales gets all sales
func (h *SaleHandler) GetSales(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Obteniendo ventas")
}

// GetSalesByCenter gets sales by center
func (h *SaleHandler) GetSalesByCenter(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Obteniendo ventas por centro")
}

// GetSalesPercentegeByCenterOverTotalSales gets the percentage of sales by center over total sales
func (h *SaleHandler) GetSalesPercentegeByCenterOverTotalSales(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Obteniendo porcentaje de ventas por centro sobre total de ventas")
}
