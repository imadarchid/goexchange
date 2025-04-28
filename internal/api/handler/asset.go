package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func (h *Handler) GetAssets(w http.ResponseWriter, r *http.Request) {
	assets, err := h.Queries.GetAllAssets(r.Context())
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ErrorResponse{
			Code:    "ASSET_RETRIEVAL_FAILED",
			Message: "Failed to retrieve assets",
		})
		fmt.Print(err)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(assets)
}
