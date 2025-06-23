package handlers

import (
	"github.com/gorilla/mux"
	"github.com/yuuhikaze/rgxr/storage"
	"net/http"
)

func GetTeXHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	uuid := vars["uuid"]

	tex, err := storage.GetTeX(uuid)
	if err != nil {
		http.Error(w, "TeX file not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.Write([]byte(tex))
}

func GetSVGHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	uuid := vars["uuid"]

	svg, err := storage.GetSVG(uuid)
	if err != nil {
		http.Error(w, "SVG file not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "image/svg+xml")
	w.Write([]byte(svg))
}
