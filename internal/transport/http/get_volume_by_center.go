package handler

import (
	"net/http"
)

// GetSalesByCenter gets sales by center
func (h *SaleHandler) GetVolumeByCenter(w http.ResponseWriter, r *http.Request) {
	volumeMap, err := h.service.GetVolumeByCenter(r.Context())
	if err != nil {
		EncodeError(w, http.StatusInternalServerError, err)
		return
	}

	response := make([]GetVolumeByCenterResponse, 0, len(volumeMap))
	for center, volume := range volumeMap {
		response = append(response, GetVolumeByCenterResponse{
			Center: center,
			Volume: volume,
		})
	}
	Encode(w, http.StatusOK, response)
}

type GetVolumeByCenterResponse struct {
	Center string  `json:"center"`
	Volume float64 `json:"volume"`
}
