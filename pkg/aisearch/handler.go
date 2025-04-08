package aisearch

import (
	"encoding/json"
	"net/http"
)

type DiagnosisHandler struct {
	service DeepSeekService
}

func NewDiagnosisHandler(s DeepSeekService) *DiagnosisHandler {
	return &DiagnosisHandler{s}
}

func (h *DiagnosisHandler) SuggestHandler(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Symptoms string `json:"symptoms"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.Symptoms == "" {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	specialist, err := h.service.SuggestSpecialist(req.Symptoms)
	if err != nil {
		http.Error(w, "Failed to get suggestion: "+err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{
		"specialist": specialist,
	})
}
