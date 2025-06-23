// extended_handlers.go
package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/yuuhikaze/rgxr/logic"
	"io"
	"net/http"
)

// IntersectionRequest represents request for FA intersection
type IntersectionRequest struct {
	UUIDs []string `json:"uuids"`
}

// IntersectionHandler handles intersection of multiple FAs
func IntersectionHandler(w http.ResponseWriter, r *http.Request) {
	var req IntersectionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	if len(req.UUIDs) < 2 {
		http.Error(w, "Need at least two FAs for intersection", http.StatusBadRequest)
		return
	}

	// Load FAs from PostgREST API
	var automata []*logic.FA
	for _, uuid := range req.UUIDs {
		fa, err := loadFAFromAPI(uuid)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error loading FA %s: %v", uuid, err), http.StatusInternalServerError)
			return
		}
		automata = append(automata, fa)
	}

	// Perform intersection
	result, err := logic.Intersection(automata)
	if err != nil {
		http.Error(w, "Intersection error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

// NFAToDFAHandler converts NFA to DFA
func NFAToDFAHandler(w http.ResponseWriter, r *http.Request) {
	uuid := r.URL.Query().Get("uuid")
	if uuid == "" {
		http.Error(w, "Missing uuid parameter", http.StatusBadRequest)
		return
	}

	nfa, err := loadFAFromAPI(uuid)
	if err != nil {
		http.Error(w, "Error loading NFA: "+err.Error(), http.StatusInternalServerError)
		return
	}

	dfa, err := logic.NFAToDFA(nfa)
	if err != nil {
		http.Error(w, "NFA to DFA conversion error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(dfa)
}

// FAToRegexHandler converts FA to regular expression
func FAToRegexHandler(w http.ResponseWriter, r *http.Request) {
	uuid := r.URL.Query().Get("uuid")
	if uuid == "" {
		http.Error(w, "Missing uuid parameter", http.StatusBadRequest)
		return
	}

	fa, err := loadFAFromAPI(uuid)
	if err != nil {
		http.Error(w, "Error loading FA: "+err.Error(), http.StatusInternalServerError)
		return
	}

	regex, err := logic.FAToRegex(fa)
	if err != nil {
		http.Error(w, "FA to regex conversion error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.Write([]byte(regex))
}

// RegexToNFARequest represents request for regex to NFA conversion
type RegexToNFARequest struct {
	Regex string `json:"regex"`
}

// RegexToNFAHandler converts regular expression to NFA
func RegexToNFAHandler(w http.ResponseWriter, r *http.Request) {
	var req RegexToNFARequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	if req.Regex == "" {
		http.Error(w, "Missing regex field", http.StatusBadRequest)
		return
	}

	nfa, err := logic.RegexToNFAComplete(req.Regex)
	if err != nil {
		http.Error(w, "Regex to NFA conversion error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(nfa)
}

// MinimizeDFAHandler minimizes a DFA
func MinimizeDFAHandler(w http.ResponseWriter, r *http.Request) {
	uuid := r.URL.Query().Get("uuid")
	if uuid == "" {
		http.Error(w, "Missing uuid parameter", http.StatusBadRequest)
		return
	}

	dfa, err := loadFAFromAPI(uuid)
	if err != nil {
		http.Error(w, "Error loading DFA: "+err.Error(), http.StatusInternalServerError)
		return
	}

	minimized, err := logic.MinimizeDFA(dfa)
	if err != nil {
		http.Error(w, "DFA minimization error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(minimized)
}

// ComplementHandler returns the complement of an FA
func ComplementHandler(w http.ResponseWriter, r *http.Request) {
	uuid := r.URL.Query().Get("uuid")
	if uuid == "" {
		http.Error(w, "Missing uuid parameter", http.StatusBadRequest)
		return
	}

	fa, err := loadFAFromAPI(uuid)
	if err != nil {
		http.Error(w, "Error loading FA: "+err.Error(), http.StatusInternalServerError)
		return
	}

	complement := logic.Complement(fa)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(complement)
}

// ConcatenationRequest represents request for FA concatenation
type ConcatenationRequest struct {
	UUIDs []string `json:"uuids"`
}

// ConcatenationHandler handles concatenation of multiple FAs
func ConcatenationHandler(w http.ResponseWriter, r *http.Request) {
	var req ConcatenationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	if len(req.UUIDs) == 0 {
		http.Error(w, "Need at least one FA for concatenation", http.StatusBadRequest)
		return
	}

	// Load FAs from PostgREST API
	var automata []*logic.FA
	for _, uuid := range req.UUIDs {
		fa, err := loadFAFromAPI(uuid)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error loading FA %s: %v", uuid, err), http.StatusInternalServerError)
			return
		}
		automata = append(automata, fa)
	}

	// Perform concatenation
	result, err := logic.Concatenation(automata)
	if err != nil {
		http.Error(w, "Concatenation error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

// Helper function to load FA from PostgREST API
func loadFAFromAPI(uuid string) (*logic.FA, error) {
	url := fmt.Sprintf("http://postgrest:3000/finite_automatas?id=eq.%s", uuid)
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("error fetching FA from PostgREST: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("error fetching FA: %s", resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %v", err)
	}

	var faArray []struct {
		ID          string   `json:"id"`
		Description *string  `json:"description"`
		Tuple       logic.FA `json:"tuple"`
		Render      string   `json:"render"`
		CreatedAt   string   `json:"created_at"`
	}

	if err := json.Unmarshal(body, &faArray); err != nil {
		return nil, fmt.Errorf("error unmarshalling response JSON: %v", err)
	}

	if len(faArray) == 0 {
		return nil, fmt.Errorf("no FA found for UUID %s", uuid)
	}

	// Extract the FA from the tuple field
	fa := faArray[0].Tuple
	return &fa, nil
}
