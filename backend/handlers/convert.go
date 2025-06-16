package handlers

import (
	"encoding/json"
	"net/http"
	"github.com/yuuhikaze/rgxr/logic"
)

func ConvertHandler(w http.ResponseWriter, r *http.Request) {
	var fa logic.FA
	json.NewDecoder(r.Body).Decode(&fa)

	dot := logic.ToDot(fa)
	tex, err := logic.DotToTex(dot)
	if err != nil {
		http.Error(w, "dot2tex error: "+err.Error(), 500)
		return
	}

	w.Write([]byte(tex))
}
