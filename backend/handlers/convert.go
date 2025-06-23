package handlers

import (
	"encoding/json"
	"github.com/google/uuid"
	"github.com/yuuhikaze/rgxr/logic"
	"regexp"
	"github.com/yuuhikaze/rgxr/storage"
	"net/http"
)

type ConvertRequest struct {
	FA   *logic.FA `json:"fa,omitempty"`
	UUID string    `json:"uuid,omitempty"`
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
	// Fix pipe symbols in TikZ code
	tex = fixPipeSymbols(tex)

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

func fixPipeSymbols(tex string) string {
	// Regex to match node identifiers in parentheses: (q1|q2|...|qn)
	// This captures content within parentheses that contains pipe symbols
	nodeIDRegex := regexp.MustCompile(`\(([^)]*\|[^)]*)\)`)

	// Regex to match node labels in braces: {q1|q2|...|qn}
	// This captures content within braces that contains pipe symbols
	nodeLabelRegex := regexp.MustCompile(`\{([^}]*\|[^}]*)\}`)

	// Fix node identifiers: remove pipes from (q1|q2) -> (q1q2)
	tex = nodeIDRegex.ReplaceAllStringFunc(tex, func(match string) string {
		// Remove the outer parentheses, replace pipes, then add parentheses back
		content := match[1 : len(match)-1] // Remove ( and )
		content = regexp.MustCompile(`\|`).ReplaceAllString(content, "")
		return "(" + content + ")"
	})

	// Fix node labels: replace pipes with $|$ in {q1|q2} -> {q1$|$q2}
	tex = nodeLabelRegex.ReplaceAllStringFunc(tex, func(match string) string {
		// Remove the outer braces, replace pipes with $|$, then add braces back
		content := match[1 : len(match)-1] // Remove { and }
		content = regexp.MustCompile(`\|`).ReplaceAllString(content, "$|$")
		return "{" + content + "}"
	})

	return tex
}
