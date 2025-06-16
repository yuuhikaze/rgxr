package handlers

import (
	"encoding/json"
	"net/http"
	"github.com/yuuhikaze/rgxr/storage"
)

type SaveImageRequest struct {
	UUID string `json:"uuid"`
	SVG  string `json:"svg"` // raw SVG content
}

func SaveImageHandler(w http.ResponseWriter, r *http.Request) {
	var req SaveImageRequest
	json.NewDecoder(r.Body).Decode(&req)

	err := storage.SaveSVG(req.UUID, req.SVG)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	w.WriteHeader(http.StatusCreated)
}
