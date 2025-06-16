package handlers

import (
	"encoding/json"
	"net/http"
	"github.com/yuuhikaze/rgxr/logic"
)

type UnionRequest struct {
	UUIDs []string `json:"uuids"`
}

func UnionHandler(w http.ResponseWriter, r *http.Request) {
	var req UnionRequest
	json.NewDecoder(r.Body).Decode(&req)

	unionedFA, err := logic.UnionFAs(req.UUIDs)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	json.NewEncoder(w).Encode(unionedFA)
}
