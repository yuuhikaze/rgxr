package handlers

import (
	"encoding/json"
	"github.com/yuuhikaze/rgxr/logic"
	"net/http"
)

// RunStringRequest represents a request to run a string through an FA
type RunStringRequest struct {
	UUID   string `json:"uuid"`
	String string `json:"string"`
}

// RunStringResponse represents the result of running a string through an FA
type RunStringResponse struct {
	Accepted bool     `json:"accepted"`
	Path     []string `json:"path"`
}

// RunStringHandler runs a string through an FA and returns whether it's accepted
func RunStringHandler(w http.ResponseWriter, r *http.Request) {
	var req RunStringRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	if req.UUID == "" || req.String == "" {
		http.Error(w, "Missing uuid or string parameter", http.StatusBadRequest)
		return
	}

	fa, err := loadFAFromAPI(req.UUID)
	if err != nil {
		http.Error(w, "Error loading FA: "+err.Error(), http.StatusInternalServerError)
		return
	}

	accepted, path := logic.RunString(fa, req.String)

	resp := RunStringResponse{
		Accepted: accepted,
		Path:     path,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}
