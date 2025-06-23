package logic

import (
	"encoding/json"
	"errors"
	"slices"
)

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
	return slices.Contains(slice, val)
}

// getNextState returns the next state(s) for fa at state index and symbol index.
func getNextState(fa *FA, state string, symbolIdx int) any {
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

// interfaceToStateSlice converts any to []string
func interfaceToStateSlice(next any) []string {
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
	case []any:
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

func getStateIndex(fa *FA, state any) int {
	if stateStr, ok := state.(string); ok && stateStr != "@v" {
		for i, s := range fa.States {
			if s == stateStr {
				return i
			}
		}
	}
	return -1
}
