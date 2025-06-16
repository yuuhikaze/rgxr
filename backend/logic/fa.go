package logic

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
)

// FA represents a finite automaton.
type FA struct {
	Alphabet    []string        `json:"alphabet"`
	States      []string        `json:"states"`
	Initial     string          `json:"initial"`
	Acceptance  []string        `json:"acceptance"`
	Transitions [][]interface{} `json:"transitions"` // 2D array; each cell string or []string
}

// ParseFAFromJSON parses a FA from raw JSON bytes.
func ParseFAFromJSON(data []byte) (*FA, error) {
	var fa FA
	err := json.Unmarshal(data, &fa)
	if err != nil {
		return nil, err
	}
	if len(fa.States) == 0 || len(fa.Alphabet) == 0 {
		return nil, errors.New("invalid FA: empty states or alphabet")
	}
	return &fa, nil
}

// Contains returns true if slice contains val.
func Contains(slice []string, val string) bool {
	for _, s := range slice {
		if s == val {
			return true
		}
	}
	return false
}

// Union creates a new FA representing the union of two FAs.
// Note: for brevity, this example assumes the alphabets are identical.
// It builds a product FA where acceptance states are union of both.
func Union(fa1, fa2 *FA) (*FA, error) {
	if len(fa1.Alphabet) != len(fa2.Alphabet) {
		return nil, errors.New("alphabets differ")
	}
	for i := range fa1.Alphabet {
		if fa1.Alphabet[i] != fa2.Alphabet[i] {
			return nil, errors.New("alphabets differ")
		}
	}

	// Build new states as concatenation "fa1State|fa2State"
	newStates := []string{}
	stateIndexMap := make(map[string]int) // maps combined state name to its index

	for _, s1 := range fa1.States {
		for _, s2 := range fa2.States {
			newState := s1 + "|" + s2
			stateIndexMap[newState] = len(newStates)
			newStates = append(newStates, newState)
		}
	}

	// Initial state is combined initial states
	newInitial := fa1.Initial + "|" + fa2.Initial

	// New acceptance states: any combined state where one part is acceptance in original
	newAcceptance := []string{}
	for _, state := range newStates {
		parts := strings.Split(state, "|")
		if len(parts) == 2 && (Contains(fa1.Acceptance, parts[0]) || Contains(fa2.Acceptance, parts[1])) {
			newAcceptance = append(newAcceptance, state)
		}
	}

	// Build transitions: for each new state, for each symbol, get next state from both FAs
	newTransitions := make([][]interface{}, len(newStates))
	
	for i, state := range newStates {
		parts := strings.Split(state, "|")
		if len(parts) != 2 {
			continue // Skip malformed states
		}
		
		row := make([]interface{}, len(fa1.Alphabet))
		
		for j := range fa1.Alphabet {
			// Find next states in fa1 and fa2
			next1 := getNextState(fa1, parts[0], j)
			next2 := getNextState(fa2, parts[1], j)
			
			// Combine next states
			combinedNext := combineNextStates(next1, next2)
			
			// Convert to proper combined state names
			if combinedNext == "@v" {
				row[j] = "@v"
			} else if nextStates, ok := combinedNext.([]string); ok {
				// Multiple next states
				var validNextStates []string
				for _, ns := range nextStates {
					if Contains(newStates, ns) {
						validNextStates = append(validNextStates, ns)
					}
				}
				if len(validNextStates) == 0 {
					row[j] = "@v"
				} else if len(validNextStates) == 1 {
					row[j] = validNextStates[0]
				} else {
					row[j] = validNextStates
				}
			} else if nextState, ok := combinedNext.(string); ok {
				if Contains(newStates, nextState) {
					row[j] = nextState
				} else {
					row[j] = "@v"
				}
			} else {
				row[j] = "@v"
			}
		}
		newTransitions[i] = row
	}

	return &FA{
		Alphabet:    fa1.Alphabet,
		States:      newStates,
		Initial:     newInitial,
		Acceptance:  newAcceptance,
		Transitions: newTransitions,
	}, nil
}

// combineNextStates combines the next states from two FAs into combined state names
func combineNextStates(next1, next2 interface{}) interface{} {
	// Handle case where either is "@v" (void/error state)
	if next1 == "@v" && next2 == "@v" {
		return "@v"
	}
	
	// Convert to string slices for easier handling
	states1 := interfaceToStateSlice(next1)
	states2 := interfaceToStateSlice(next2)
	
	// If either is empty (void), treat as empty
	if len(states1) == 0 && len(states2) == 0 {
		return "@v"
	}
	
	// Handle case where one is void
	if len(states1) == 0 {
		states1 = []string{"@v"}
	}
	if len(states2) == 0 {
		states2 = []string{"@v"}
	}
	
	// Create all combinations
	var combined []string
	for _, s1 := range states1 {
		for _, s2 := range states2 {
			combinedState := s1 + "|" + s2
			combined = append(combined, combinedState)
		}
	}
	
	// Remove duplicates
	combined = uniqueStrings(combined)
	
	if len(combined) == 0 {
		return "@v"
	} else if len(combined) == 1 {
		return combined[0]
	} else {
		return combined
	}
}

// interfaceToStateSlice converts interface{} to []string
func interfaceToStateSlice(next interface{}) []string {
	if next == "@v" || next == nil {
		return []string{}
	}
	
	switch v := next.(type) {
	case string:
		if v == "@v" {
			return []string{}
		}
		return []string{v}
	case []string:
		var result []string
		for _, s := range v {
			if s != "@v" {
				result = append(result, s)
			}
		}
		return result
	case []interface{}:
		var result []string
		for _, item := range v {
			if s, ok := item.(string); ok && s != "@v" {
				result = append(result, s)
			}
		}
		return result
	default:
		return []string{}
	}
}

// getNextState returns the next state(s) for fa at state index and symbol index.
func getNextState(fa *FA, state string, symbolIdx int) interface{} {
	// Find row index for state
	rowIdx := -1
	for i, s := range fa.States {
		if s == state {
			rowIdx = i
			break
		}
	}
	if rowIdx == -1 {
		return "@v"
	}
	if rowIdx >= len(fa.Transitions) || symbolIdx >= len(fa.Transitions[rowIdx]) {
		return "@v"
	}
	return fa.Transitions[rowIdx][symbolIdx]
}

// uniqueStrings removes duplicates.
func uniqueStrings(arr []string) []string {
	set := make(map[string]bool)
	res := []string{}
	for _, v := range arr {
		if !set[v] {
			set[v] = true
			res = append(res, v)
		}
	}
	return res
}

// UnionFAs loads FAs by UUIDs via the PostgREST API and unions them all (left-associative), returns result.
func UnionFAs(uuids []string) (*FA, error) {
	if len(uuids) < 2 {
		return nil, fmt.Errorf("need at least two FAs to union")
	}

	// Load from PostgREST API
	var automata []*FA
	for _, uuid := range uuids {
		url := fmt.Sprintf("http://postgrest:3000/finite_automatas?id=eq.%s", uuid)
		resp, err := http.Get(url)
		if err != nil {
			return nil, fmt.Errorf("error fetching FA from PostgREST: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			return nil, fmt.Errorf("error fetching FA: %s", resp.Status)
		}

		// Read the JSON response
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("error reading response body: %v", err)
		}

		var faArray []struct {
			ID          string          `json:"id"`
			Description *string         `json:"description"`
			Tuple       FA              `json:"tuple"`
			Render      string          `json:"render"`
			CreatedAt   string          `json:"created_at"`
		}
		
		if err := json.Unmarshal(body, &faArray); err != nil {
			return nil, fmt.Errorf("error unmarshalling response JSON: %v", err)
		}

		if len(faArray) == 0 {
			return nil, fmt.Errorf("no FA found for UUID %s", uuid)
		}

		// Extract the FA from the tuple field
		fa := faArray[0].Tuple
		automata = append(automata, &fa)
	}

	// Perform union pairwise
	result := automata[0]
	for i := 1; i < len(automata); i++ {
		merged, err := Union(result, automata[i])
		if err != nil {
			return nil, err
		}
		result = merged
	}

	return result, nil
}
