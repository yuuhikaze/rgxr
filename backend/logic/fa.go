package logic

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"os/exec"
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
	stateMap := map[string]bool{}

	for _, s1 := range fa1.States {
		for _, s2 := range fa2.States {
			newState := s1 + "|" + s2
			stateMap[newState] = true
			newStates = append(newStates, newState)
		}
	}

	// Initial state is combined initial states
	newInitial := fa1.Initial + "|" + fa2.Initial

	// New acceptance states: any combined state where one part is acceptance in original
	newAcceptance := []string{}
	for state := range stateMap {
		parts := strings.Split(state, "|")
		if Contains(fa1.Acceptance, parts[0]) || Contains(fa2.Acceptance, parts[1]) {
			newAcceptance = append(newAcceptance, state)
		}
	}

	// Build transitions: for each new state, for each symbol, get next state from both FAs
	newTransitions := make([][]interface{}, len(newStates))
	for i, state := range newStates {
		parts := strings.Split(state, "|")
		row := make([]interface{}, len(fa1.Alphabet))
		for j := range fa1.Alphabet {
			// Find index of symbol in alphabet (should be j)
			// Find next states in fa1 and fa2
			next1 := getNextState(fa1, parts[0], j)
			next2 := getNextState(fa2, parts[1], j)
			if next1 == "@v" && next2 == "@v" {
				row[j] = "@v"
			} else {
				// Combine next states; if either is multiple, flatten to array
				nextStates := []string{}
				if ns, ok := next1.([]string); ok {
					nextStates = append(nextStates, ns...)
				} else if s, ok := next1.(string); ok && s != "@v" {
					nextStates = append(nextStates, s)
				}
				if ns, ok := next2.([]string); ok {
					nextStates = append(nextStates, ns...)
				} else if s, ok := next2.(string); ok && s != "@v" {
					nextStates = append(nextStates, s)
				}
				// Remove duplicates
				nextStates = uniqueStrings(nextStates)
				if len(nextStates) == 0 {
					row[j] = "@v"
				} else if len(nextStates) == 1 {
					row[j] = nextStates[0]
				} else {
					row[j] = nextStates
				}
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
	if symbolIdx >= len(fa.Transitions[rowIdx]) {
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

// UnionFAs loads FAs by UUIDs, unions them all (left-associative), returns result
func UnionFAs(uuids []string) (*FA, error) {
	if len(uuids) < 2 {
		return nil, fmt.Errorf("need at least two FAs to union")
	}

	// Load from DB or file system — here we simulate via JSON in files for simplicity.
	var automata []*FA

	for _, uuid := range uuids {
		cmd := exec.Command("psql", "-U", "authenticator", "-d", "postgres", "-c",
			fmt.Sprintf("SELECT tuple FROM api.finite_automatas WHERE id='%s';", uuid))
		out, err := cmd.Output()
		if err != nil {
			return nil, err
		}

		// Extract JSON from psql output (this is a simplification — better to query through Go pgx or pgxpool)
		jsonStart := bytes.IndexByte(out, '{')
		jsonEnd := bytes.LastIndexByte(out, '}') + 1
		jsonData := out[jsonStart:jsonEnd]

		fa, err := ParseFAFromJSON(jsonData)
		if err != nil {
			return nil, err
		}
		automata = append(automata, fa)
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
