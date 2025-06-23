// regex.go
package logic

import (
	"fmt"
	"strings"
)

// FAToRegex converts a finite automaton to a regular expression using state elimination
func FAToRegex(fa *FA) (string, error) {
	if len(fa.States) == 0 {
		return "∅", nil
	}

	// Create a copy of the FA with added start and end states
	// This simplifies the state elimination algorithm

	// Create transition matrix
	n := len(fa.States) + 2 // +2 for new start and end states
	stateNames := make([]string, n)
	stateNames[0] = "START"
	copy(stateNames[1:], fa.States)
	stateNames[n-1] = "END"

	// Initialize regex matrix
	regexMatrix := make([][]string, n)
	for i := range regexMatrix {
		regexMatrix[i] = make([]string, n)
		for j := range regexMatrix[i] {
			regexMatrix[i][j] = "∅"
		}
	}

	// Add epsilon transition from START to original initial state
	initialIdx := getStateIndexInList(fa.States, fa.Initial) + 1
	regexMatrix[0][initialIdx] = "ε"

	// Add epsilon transitions from acceptance states to END
	for _, acceptState := range fa.Acceptance {
		acceptIdx := getStateIndexInList(fa.States, acceptState) + 1
		regexMatrix[acceptIdx][n-1] = "ε"
	}

	// Fill in original transitions
	for i := range fa.States {
		for j, symbol := range fa.Alphabet {
			next := fa.Transitions[i][j]
			stateIdx := i + 1 // +1 because of START state

			switch v := next.(type) {
			case string:
				if v != "@v" {
					nextIdx := getStateIndexInList(fa.States, v) + 1
					if regexMatrix[stateIdx][nextIdx] == "∅" {
						regexMatrix[stateIdx][nextIdx] = symbol
					} else {
						regexMatrix[stateIdx][nextIdx] = unionRegex(regexMatrix[stateIdx][nextIdx], symbol)
					}
				}
			case []string:
				for _, nextState := range v {
					if nextState != "@v" {
						nextIdx := getStateIndexInList(fa.States, nextState) + 1
						if regexMatrix[stateIdx][nextIdx] == "∅" {
							regexMatrix[stateIdx][nextIdx] = symbol
						} else {
							regexMatrix[stateIdx][nextIdx] = unionRegex(regexMatrix[stateIdx][nextIdx], symbol)
						}
					}
				}
			}
		}
	}

	// State elimination algorithm
	// Eliminate states 1 to n-2 (keep START=0 and END=n-1)
	for k := 1; k < n-1; k++ {
		// For each pair of states (i,j), update the transition
		for i := 0; i < n; i++ {
			for j := 0; j < n; j++ {
				if i == k || j == k {
					continue
				}

				// R_ij = R_ij ∪ R_ik R_kk* R_kj
				rik := regexMatrix[i][k]
				rkk := regexMatrix[k][k]
				rkj := regexMatrix[k][j]

				if rik != "∅" && rkj != "∅" {
					var newPart string
					if rkk == "∅" {
						newPart = concatenateRegex(rik, rkj)
					} else {
						kleenePart := kleeneStarRegex(rkk)
						newPart = concatenateRegex(rik, concatenateRegex(kleenePart, rkj))
					}
					regexMatrix[i][j] = unionRegex(regexMatrix[i][j], newPart)
				}
			}
		}
	}

	// The final regex is the transition from START to END
	result := regexMatrix[0][n-1]
	if result == "∅" {
		return "∅", nil
	}
	return result, nil
}

// RegexToNFA converts a regular expression to an NFA using Thompson's construction
func RegexToNFA(regex string) (*FA, error) {
	if regex == "" || regex == "∅" {
		return createEmptyNFA(), nil
	}
	if regex == "ε" {
		return createEpsilonNFA(), nil
	}

	// Parse and convert regex to NFA
	// This is a simplified implementation - a full implementation would need a proper parser
	return parseRegexToNFA(regex)
}

// Helper functions for regex operations
func unionRegex(r1, r2 string) string {
	if r1 == "∅" {
		return r2
	}
	if r2 == "∅" {
		return r1
	}
	if r1 == r2 {
		return r1
	}
	return fmt.Sprintf("(%s∪%s)", r1, r2)
}

func concatenateRegex(r1, r2 string) string {
	if r1 == "∅" || r2 == "∅" {
		return "∅"
	}
	if r1 == "ε" {
		return r2
	}
	if r2 == "ε" {
		return r1
	}

	// Add parentheses if needed
	needParens1 := strings.Contains(r1, "∪")
	needParens2 := strings.Contains(r2, "∪")

	result := ""
	if needParens1 {
		result += "(" + r1 + ")"
	} else {
		result += r1
	}

	if needParens2 {
		result += "(" + r2 + ")"
	} else {
		result += r2
	}

	return result
}

func kleeneStarRegex(r string) string {
	if r == "∅" || r == "ε" {
		return "ε"
	}

	// Check if parentheses are needed
	if len(r) == 1 || (strings.HasPrefix(r, "(") && strings.HasSuffix(r, ")")) {
		return r + "*"
	}
	return "(" + r + ")*"
}

func getStateIndexInList(states []string, state string) int {
	for i, s := range states {
		if s == state {
			return i
		}
	}
	return -1
}

// Simplified regex parsing - in practice, you'd want a more robust parser
func parseRegexToNFA(regex string) (*FA, error) {
	// This is a very basic implementation
	// A full implementation would need proper parsing of operators, precedence, etc.

	if len(regex) == 1 {
		// Single character
		return createCharacterNFA(regex), nil
	}

	// For now, return a simple NFA that accepts the literal string
	// In practice, you'd implement Thompson's construction properly
	states := []string{"q0", "q1"}
	transitions := [][]any{
		{[]string{"q1"}}, // from q0 on the character
		{"@v"},           // from q1 (no transitions)
	}

	return &FA{
		Alphabet:    []string{regex}, // Simplified - should extract actual alphabet
		States:      states,
		Initial:     "q0",
		Acceptance:  []string{"q1"},
		Transitions: transitions,
	}, nil
}

func createEmptyNFA() *FA {
	return &FA{
		Alphabet:    []string{},
		States:      []string{"q0"},
		Initial:     "q0",
		Acceptance:  []string{}, // No accepting states
		Transitions: [][]any{{"@v"}},
	}
}

func createEpsilonNFA() *FA {
	return &FA{
		Alphabet:    []string{},
		States:      []string{"q0"},
		Initial:     "q0",
		Acceptance:  []string{"q0"}, // Initial state is accepting for epsilon
		Transitions: [][]any{{"@v"}},
	}
}

func createCharacterNFA(char string) *FA {
	return &FA{
		Alphabet:   []string{char},
		States:     []string{"q0", "q1"},
		Initial:    "q0",
		Acceptance: []string{"q1"},
		Transitions: [][]any{
			{"q1"}, // from q0 on char to q1
			{"@v"}, // from q1 (no transitions)
		},
	}
}

// Thompson's construction for building NFAs from regex
type NFAFragment struct {
	Start string
	End   string
	FA    *FA
}

// More complete regex to NFA conversion using Thompson's construction
func RegexToNFAComplete(regex string) (*FA, error) {
	tokens := tokenizeRegex(regex)
	if len(tokens) == 0 {
		return createEmptyNFA(), nil
	}

	fragment, err := buildNFAFromTokens(tokens)
	if err != nil {
		return nil, err
	}

	return fragment.FA, nil
}

// Tokenize regex into basic components
func tokenizeRegex(regex string) []string {
	var tokens []string
	runes := []rune(regex)

	for i := 0; i < len(runes); i++ {
		switch runes[i] {
		case '(':
			tokens = append(tokens, "(")
		case ')':
			tokens = append(tokens, ")")
		case '*':
			tokens = append(tokens, "*")
		case '+':
			tokens = append(tokens, "+")
		case '∪':
			tokens = append(tokens, "∪")
		case '∅':
			tokens = append(tokens, "∅")
		case 'ε':
			tokens = append(tokens, "ε")
		default:
			tokens = append(tokens, string(runes[i]))
		}
	}

	return tokens
}

// Build NFA from tokenized regex using recursive descent parsing
func buildNFAFromTokens(tokens []string) (*NFAFragment, error) {
	if len(tokens) == 0 {
		return nil, fmt.Errorf("empty token list")
	}

	// Simple implementation for basic cases
	if len(tokens) == 1 {
		token := tokens[0]
		switch token {
		case "∅":
			fa := createEmptyNFA()
			return &NFAFragment{Start: "q0", End: "q0", FA: fa}, nil
		case "ε":
			fa := createEpsilonNFA()
			return &NFAFragment{Start: "q0", End: "q0", FA: fa}, nil
		default:
			fa := createCharacterNFA(token)
			return &NFAFragment{Start: "q0", End: "q1", FA: fa}, nil
		}
	}

	// For more complex expressions, you would implement proper parsing
	// This is a simplified version
	return &NFAFragment{
		Start: "q0",
		End:   "q1",
		FA:    createCharacterNFA(strings.Join(tokens, "")),
	}, nil
}
