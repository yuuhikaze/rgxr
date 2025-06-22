package handlers

import (
    "encoding/json"
    "net/http"
    "github.com/google/uuid"
    "github.com/yuuhikaze/rgxr/logic"
    "github.com/yuuhikaze/rgxr/storage"
)

type ConvertRequest struct {
    FA *logic.FA `json:"fa,omitempty"`
    UUID string `json:"uuid,omitempty"`
}

type ConvertResponse struct {
    ID  string `json:"id"`
    SVG string `json:"svg"`
    TeX string `json:"tex"`
    DOT string `json:"dot"`
}

func ConvertHandler(w http.ResponseWriter, r *http.Request) {
    var req ConvertRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, "Invalid JSON: "+err.Error(), http.StatusBadRequest)
        return
    }

    var fa *logic.FA
    
    // Get FA either from request body or by loading from UUID
    if req.FA != nil {
        fa = req.FA
    } else if req.UUID != "" {
        loadedFA, err := loadFAFromAPI(req.UUID)
        if err != nil {
            http.Error(w, "Error loading FA: "+err.Error(), http.StatusInternalServerError)
            return
        }
        fa = loadedFA
    } else {
        http.Error(w, "Must provide either FA or UUID", http.StatusBadRequest)
        return
    }

    // Generate unique ID for this render
    id := uuid.New().String()

    // Convert FA to DOT
    dot := logic.ToDot(*fa)

    // Convert DOT to TikZ
    tex, err := logic.DotToTex(dot)
    if err != nil {
        http.Error(w, "dot2tex error: "+err.Error(), http.StatusInternalServerError)
        return
    }

    // Save TeX file
    if err := storage.SaveTeX(id, tex); err != nil {
        http.Error(w, "Failed to save TeX: "+err.Error(), http.StatusInternalServerError)
        return
    }

    // Convert TikZ to SVG
    svg, err := logic.TikZToSVG(tex)
    if err != nil {
        http.Error(w, "TikZ to SVG conversion error: "+err.Error(), http.StatusInternalServerError)
        return
    }

    // Save SVG file
    if err := storage.SaveSVG(id, svg); err != nil {
        http.Error(w, "Failed to save SVG: "+err.Error(), http.StatusInternalServerError)
        return
    }

    // Return response
    resp := ConvertResponse{
        ID:  id,
        SVG: svg,
        TeX: tex,
        DOT: dot,
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(resp)
}
