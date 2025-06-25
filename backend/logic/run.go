package logic

import (
	"strings"
)

// RunString runs a string through an FA and returns whether it's accepted and the path taken
func RunString(fa *FA, input string) (bool, []string) {
	// Handle epsilon symbol
	epsilonIdx := -1
	for i, symbol := range fa.Alphabet {
		if symbol == "@e" {
			epsilonIdx = i
			break
		}
	}

	// Convert string to runes for proper character handling
	runes := []rune(input)
	
	// Start from initial state (considering epsilon closure for NFAs)
	currentStates := []string{fa.Initial}
	if epsilonIdx >= 0 {
		currentStates = computeEpsilonClosure(fa, currentStates, epsilonIdx)
	}
	
	path := []string{strings.Join(currentStates, ",")}

	// Process each character
	for _, char := range runes {
		// Find symbol index
		symbolIdx := -1
		for i, symbol := range fa.Alphabet {
			if symbol == string(char) {
				symbolIdx = i
				break
			}
		}

		if symbolIdx == -1 {
			// Symbol not in alphabet
			return false, path
		}

		// Compute next states
		nextStates := make(map[string]bool)
		
		for _, state := range currentStates {
			stateIdx := -1
			for i, s := range fa.States {
				if s == state {
					stateIdx = i
					break
				}
			}

			if stateIdx >= 0 && stateIdx < len(fa.Transitions) &&
				symbolIdx < len(fa.Transitions[stateIdx]) {
				next := fa.Transitions[stateIdx][symbolIdx]

				switch v := next.(type) {
				case string:
					if v != "" && v != "@v" {
						nextStates[v] = true
					}
				case []string:
					for _, ns := range v {
						if ns != "" && ns != "@v" {
							nextStates[ns] = true
						}
					}
				case []any:
					for _, item := range v {
						if s, ok := item.(string); ok && s != "" && s != "@v" {
							nextStates[s] = true
						}
					}
				}
			}
		}

		// Convert map to slice
		currentStates = []string{}
		for state := range nextStates {
			currentStates = append(currentStates, state)
		}

		// Apply epsilon closure if needed
		if epsilonIdx >= 0 {
			currentStates = computeEpsilonClosure(fa, currentStates, epsilonIdx)
		}

		// Add to path
		if len(currentStates) == 0 {
			path = append(path, "∅")
			return false, path
		}
		path = append(path, strings.Join(currentStates, ","))
	}

	// Check if any current state is accepting
	for _, state := range currentStates {
		if Contains(fa.Acceptance, state) {
			return true, path
		}
	}

	return false, path
}

// computeEpsilonClosure computes the epsilon closure of a set of states
func computeEpsilonClosure(fa *FA, states []string, epsilonIdx int) []string {
	closure := make(map[string]bool)
	stack := make([]string, len(states))
	copy(stack, states)

	for _, s := range states {
		closure[s] = true
	}

	for len(stack) > 0 {
		current := stack[len(stack)-1]
		stack = stack[:len(stack)-1]

		// Find state index
		stateIdx := -1
		for i, s := range fa.States {
			if s == current {
				stateIdx = i
				break
			}
		}

		if stateIdx >= 0 && stateIdx < len(fa.Transitions) &&
			epsilonIdx < len(fa.Transitions[stateIdx]) {
			next := fa.Transitions[stateIdx][epsilonIdx]

			switch v := next.(type) {
			case string:
				if v != "" && v != "@v" && !closure[v] {
					closure[v] = true
					stack = append(stack, v)
				}
			case []string:
				for _, ns := range v {
					if ns != "" && ns != "@v" && !closure[ns] {
						closure[ns] = true
						stack = append(stack, ns)
					}
				}
			case []any:
				for _, item := range v {
					if s, ok := item.(string); ok && s != "" && s != "@v" && !closure[s] {
						closure[s] = true
						stack = append(stack, s)
					}
				}
			}
		}
	}

	result := make([]string, 0, len(closure))
	for state := range closure {
		result = append(result, state)
	}
	return result
}
