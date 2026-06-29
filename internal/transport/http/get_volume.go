package handler

import (
	"net/http"
)

// GetSales gets all sales
func (h *SaleHandler) GetTotalVolume(w http.ResponseWriter, r *http.Request) {
	volume, err := h.service.GetTotalVolume(r.Context())
	if err != nil {
		EncodeError(w, http.StatusInternalServerError, err)
		return
	}
	Encode(w, http.StatusOK, GetTotalVolumeResponse{Volume: volume})
}

type GetTotalVolumeResponse struct {
	Volume float64 `json:"volume"`
}
