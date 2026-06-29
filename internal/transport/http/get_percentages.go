package handler

import (
	"net/http"
)

// GetSalesPercentegeByCenterOverTotalSales gets the percentage of sales by center over total sales
func (h *SaleHandler) GetSalesPercentegeByCenterOverTotalSales(w http.ResponseWriter, r *http.Request) {
	percentages, err := h.service.GetSalesPercentegeByCenterOverTotalSales(r.Context())
	if err != nil {
		EncodeError(w, http.StatusInternalServerError, err)
		return
	}

	Encode(w, http.StatusOK, percentages)
}
