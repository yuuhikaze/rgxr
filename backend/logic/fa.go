package logic

import (
	"fmt"
	"sort"
	"strings"
)

type BooleanMode string

const (
	Union        BooleanMode = "union"
	Intersection BooleanMode = "intersection"
	// Difference   BooleanMode = "difference" // planned support on the near future
)

// FA represents a finite automaton.
type FA struct {
	Alphabet    []string `json:"alphabet"`
	States      []string `json:"states"`
	Initial     string   `json:"initial"`
	Acceptance  []string `json:"acceptance"`
	Transitions [][]any  `json:"transitions"` // 2D array; each cell string or []string
}

// PerformBoolean creates a new FA representing the boolean operation result of multiple FAs.
// Similar to Union but acceptance states require ALL original FAs to accept.
func PerformBoolean(fas []*FA, mode BooleanMode) (*FA, error) {
	if len(fas) < 2 {
		return nil, fmt.Errorf("need at least two FAs for intersection")
	}

	// Verify all alphabets are identical
	baseAlphabet := fas[0].Alphabet
	for i := 1; i < len(fas); i++ {
		if len(fas[i].Alphabet) != len(baseAlphabet) {
			return nil, fmt.Errorf("alphabets differ")
		}
		for j := range baseAlphabet {
			if fas[i].Alphabet[j] != baseAlphabet[j] {
				return nil, fmt.Errorf("alphabets differ")
			}
		}
	}

	// Build product states
	newStates := []string{}
	stateIndexMap := make(map[string]int)

	// Generate all combinations of states
	var generateStates func([]string, int)
	generateStates = func(current []string, faIndex int) {
		if faIndex == len(fas) {
			newState := strings.Join(current, "|")
			stateIndexMap[newState] = len(newStates)
			newStates = append(newStates, newState)
			return
		}
		for _, state := range fas[faIndex].States {
			generateStates(append(current, state), faIndex+1)
		}
	}
	generateStates([]string{}, 0)

	// Initial state is combination of all initial states
	initialParts := make([]string, len(fas))
	for i, fa := range fas {
		initialParts[i] = fa.Initial
	}
	newInitial := strings.Join(initialParts, "|")

	// Acceptance states
	newAcceptance := []string{}
	for _, state := range newStates {
		parts := strings.Split(state, "|")
		if len(parts) == len(fas) {
			switch mode {
			case Intersection:
				all := true
				for i, part := range parts {
					if !Contains(fas[i].Acceptance, part) {
						all = false
						break
					}
				}
				if all {
					newAcceptance = append(newAcceptance, state)
				}
			case Union:
				for i, part := range parts {
					if Contains(fas[i].Acceptance, part) {
						newAcceptance = append(newAcceptance, state)
						break
					}
				}
			}
		}
	}

	// Build transitions
	newTransitions := make([][]any, len(newStates))
	for i, state := range newStates {
		parts := strings.Split(state, "|")
		if len(parts) != len(fas) {
			continue
		}

		row := make([]any, len(baseAlphabet))
		for j := range baseAlphabet {
			// Get next states from all FAs
			nextParts := make([]any, len(fas))
			for k, part := range parts {
				nextParts[k] = getNextState(fas[k], part, j)
			}

			// Combine next states
			combinedNext := combineBooleanStates(nextParts)
			if combinedNext == "@v" {
				row[j] = "@v"
			} else if nextStates, ok := combinedNext.([]string); ok {
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
		Alphabet:    baseAlphabet,
		States:      newStates,
		Initial:     newInitial,
		Acceptance:  newAcceptance,
		Transitions: newTransitions,
	}, nil
}

// combineBooleanStates combines next states for intersection
func combineBooleanStates(nextParts []any) any {
	// Convert all parts to state slices
	stateSets := make([][]string, len(nextParts))
	for i, part := range nextParts {
		stateSets[i] = interfaceToStateSlice(part)
		if len(stateSets[i]) == 0 {
			return "@v" // If any FA has no next state, boolean operation results in empty
		}
	}

	// Generate all combinations
	var result []string
	var generateCombinations func([]string, int)
	generateCombinations = func(current []string, index int) {
		if index == len(stateSets) {
			result = append(result, strings.Join(current, "|"))
			return
		}
		for _, state := range stateSets[index] {
			generateCombinations(append(current, state), index+1)
		}
	}
	generateCombinations([]string{}, 0)

	if len(result) == 0 {
		return "@v"
	} else if len(result) == 1 {
		return result[0]
	} else {
		return result
	}
}

// NFAToDFA converts an NFA to DFA using subset construction
func NFAToDFA(nfa *FA) (*FA, error) {
	// Helper function to compute epsilon closure
	epsilonClosure := func(states []string) []string {
		closure := make(map[string]bool)
		stack := make([]string, len(states))
		copy(stack, states)

		for _, s := range states {
			closure[s] = true
		}

		for len(stack) > 0 {
			// current := stack[len(stack)-1]
			stack = stack[:len(stack)-1]

			// Find epsilon transitions (assuming epsilon is represented as empty string or special symbol)
			// This implementation assumes no explicit epsilon transitions in your format
			// If you have epsilon transitions, you'd need to handle them here
		}

		result := make([]string, 0, len(closure))
		for state := range closure {
			result = append(result, state)
		}
		sort.Strings(result)
		return result
	}

	// Start with initial state closure
	initialClosure := epsilonClosure([]string{nfa.Initial})

	// DFA states will be sets of NFA states
	dfaStates := [][]string{initialClosure}
	stateToIndex := make(map[string]int)
	stateToIndex[strings.Join(initialClosure, ",")] = 0

	// Build DFA states and transitions
	dfaTransitions := [][]any{}
	queue := [][]string{initialClosure}
	processed := make(map[string]bool)

	for len(queue) > 0 {
		currentSet := queue[0]
		queue = queue[1:]

		currentKey := strings.Join(currentSet, ",")
		if processed[currentKey] {
			continue
		}
		processed[currentKey] = true

		row := make([]any, len(nfa.Alphabet))

		for symbolIdx := range nfa.Alphabet {
			// Compute next state set
			nextStates := make(map[string]bool)

			for _, state := range currentSet {
				stateIdx := -1
				for i, s := range nfa.States {
					if s == state {
						stateIdx = i
						break
					}
				}

				if stateIdx >= 0 && stateIdx < len(nfa.Transitions) &&
					symbolIdx < len(nfa.Transitions[stateIdx]) {
					next := nfa.Transitions[stateIdx][symbolIdx]

					switch v := next.(type) {
					case string:
						if v != "@v" {
							nextStates[v] = true
						}
					case []string:
						for _, ns := range v {
							if ns != "@v" {
								nextStates[ns] = true
							}
						}
					case []any:
						for _, item := range v {
							if s, ok := item.(string); ok && s != "@v" {
								nextStates[s] = true
							}
						}
					}
				}
			}

			if len(nextStates) == 0 {
				row[symbolIdx] = "@v"
			} else {
				nextStateSlice := make([]string, 0, len(nextStates))
				for state := range nextStates {
					nextStateSlice = append(nextStateSlice, state)
				}
				sort.Strings(nextStateSlice)

				nextClosure := epsilonClosure(nextStateSlice)
				nextKey := strings.Join(nextClosure, ",")

				if _, exists := stateToIndex[nextKey]; !exists {
					stateToIndex[nextKey] = len(dfaStates)
					dfaStates = append(dfaStates, nextClosure)
					queue = append(queue, nextClosure)
				}

				row[symbolIdx] = nextKey
			}
		}
		dfaTransitions = append(dfaTransitions, row)
	}

	// Convert state sets to state names and build final DFA
	dfaStateNames := make([]string, len(dfaStates))
	for i := range dfaStates {
		dfaStateNames[i] = fmt.Sprintf("q%d", i)
	}

	// Update transitions to use new state names
	finalTransitions := make([][]any, len(dfaTransitions))
	for i, row := range dfaTransitions {
		finalRow := make([]any, len(row))
		for j, transition := range row {
			if transition == "@v" {
				finalRow[j] = "@v"
			} else if stateKey, ok := transition.(string); ok {
				if idx, exists := stateToIndex[stateKey]; exists {
					finalRow[j] = dfaStateNames[idx]
				} else {
					finalRow[j] = "@v"
				}
			}
		}
		finalTransitions[i] = finalRow
	}

	// Determine acceptance states
	dfaAcceptance := []string{}
	for i, stateSet := range dfaStates {
		for _, nfaState := range stateSet {
			if Contains(nfa.Acceptance, nfaState) {
				dfaAcceptance = append(dfaAcceptance, dfaStateNames[i])
				break
			}
		}
	}

	return &FA{
		Alphabet:    nfa.Alphabet,
		States:      dfaStateNames,
		Initial:     dfaStateNames[0],
		Acceptance:  dfaAcceptance,
		Transitions: finalTransitions,
	}, nil
}

// Complement returns the complement of an FA (flip acceptance states)
func Complement(fa *FA) *FA {
	newAcceptance := []string{}

	for _, state := range fa.States {
		if !Contains(fa.Acceptance, state) {
			newAcceptance = append(newAcceptance, state)
		}
	}

	// Create a copy of the FA with flipped acceptance states
	newTransitions := make([][]any, len(fa.Transitions))
	for i, row := range fa.Transitions {
		newRow := make([]any, len(row))
		copy(newRow, row)
		newTransitions[i] = newRow
	}

	return &FA{
		Alphabet:    append([]string{}, fa.Alphabet...),
		States:      append([]string{}, fa.States...),
		Initial:     fa.Initial,
		Acceptance:  newAcceptance,
		Transitions: newTransitions,
	}
}

// Concatenation creates an FA representing the concatenation of multiple FAs
func Concatenation(fas []*FA) (*FA, error) {
	if len(fas) == 0 {
		return nil, fmt.Errorf("need at least one FA for concatenation")
	}
	if len(fas) == 1 {
		return fas[0], nil
	}

	// Verify all alphabets are identical
	baseAlphabet := fas[0].Alphabet
	for i := 1; i < len(fas); i++ {
		if len(fas[i].Alphabet) != len(baseAlphabet) {
			return nil, fmt.Errorf("alphabets differ")
		}
		for j := range baseAlphabet {
			if fas[i].Alphabet[j] != baseAlphabet[j] {
				return nil, fmt.Errorf("alphabets differ")
			}
		}
	}

	// Rename states to avoid conflicts
	allStates := []string{}
	stateMapping := make([]map[string]string, len(fas))

	for faIdx, fa := range fas {
		stateMapping[faIdx] = make(map[string]string)
		for _, state := range fa.States {
			newName := fmt.Sprintf("fa%d_%s", faIdx, state)
			stateMapping[faIdx][state] = newName
			allStates = append(allStates, newName)
		}
	}

	// Build concatenated transitions
	allTransitions := make([][]any, len(allStates))
	stateIndex := 0

	for faIdx, fa := range fas {
		for i, state := range fa.States {
			row := make([]any, len(baseAlphabet))

			for j := range baseAlphabet {
				originalNext := fa.Transitions[i][j]

				// If this is an accepting state and not the last FA,
				// add epsilon transitions to next FA's initial state
				isAccepting := Contains(fa.Acceptance, state)
				isLastFA := faIdx == len(fas)-1

				switch v := originalNext.(type) {
				case string:
					if v == "@v" {
						if isAccepting && !isLastFA {
							// Add transition to next FA's initial state
							nextFAInitial := stateMapping[faIdx+1][fas[faIdx+1].Initial]
							row[j] = nextFAInitial
						} else {
							row[j] = "@v"
						}
					} else {
						row[j] = stateMapping[faIdx][v]
					}
				case []string:
					newNext := []string{}
					for _, next := range v {
						if next != "@v" {
							newNext = append(newNext, stateMapping[faIdx][next])
						}
					}
					if isAccepting && !isLastFA {
						nextFAInitial := stateMapping[faIdx+1][fas[faIdx+1].Initial]
						newNext = append(newNext, nextFAInitial)
					}

					if len(newNext) == 0 {
						row[j] = "@v"
					} else if len(newNext) == 1 {
						row[j] = newNext[0]
					} else {
						row[j] = uniqueStrings(newNext)
					}
				default:
					if isAccepting && !isLastFA {
						nextFAInitial := stateMapping[faIdx+1][fas[faIdx+1].Initial]
						row[j] = nextFAInitial
					} else {
						row[j] = "@v"
					}
				}
			}
			allTransitions[stateIndex] = row
			stateIndex++
		}
	}

	// Initial state is the first FA's initial state
	newInitial := stateMapping[0][fas[0].Initial]

	// Acceptance states are only from the last FA
	lastFAIdx := len(fas) - 1
	newAcceptance := []string{}
	for _, acceptState := range fas[lastFAIdx].Acceptance {
		newAcceptance = append(newAcceptance, stateMapping[lastFAIdx][acceptState])
	}

	return &FA{
		Alphabet:    baseAlphabet,
		States:      allStates,
		Initial:     newInitial,
		Acceptance:  newAcceptance,
		Transitions: allTransitions,
	}, nil
}

// MinimizeDFA minimizes a DFA using table-filling algorithm
func MinimizeDFA(dfa *FA) (*FA, error) {
	n := len(dfa.States)
	if n <= 1 {
		return dfa, nil
	}

	// Create distinguishability table
	distinguishable := make([][]bool, n)
	for i := range distinguishable {
		distinguishable[i] = make([]bool, n)
	}

	// Mark pairs where one is accepting and other is not
	for i := range n {
		for j := i + 1; j < n; j++ {
			state1 := dfa.States[i]
			state2 := dfa.States[j]

			isAccepting1 := Contains(dfa.Acceptance, state1)
			isAccepting2 := Contains(dfa.Acceptance, state2)

			if isAccepting1 != isAccepting2 {
				distinguishable[i][j] = true
				distinguishable[j][i] = true
			}
		}
	}

	// Iteratively mark distinguishable pairs
	changed := true
	for changed {
		changed = false
		for i := range n {
			for j := i + 1; j < n; j++ {
				if !distinguishable[i][j] {
					// Check if states i and j are distinguishable
					for symbolIdx := range dfa.Alphabet {
						next1 := getNextState(dfa, dfa.States[i], symbolIdx)
						next2 := getNextState(dfa, dfa.States[j], symbolIdx)

						// Convert to state indices
						idx1 := getStateIndex(dfa, next1)
						idx2 := getStateIndex(dfa, next2)

						if idx1 != -1 && idx2 != -1 && idx1 != idx2 &&
							distinguishable[min(idx1, idx2)][max(idx1, idx2)] {
							distinguishable[i][j] = true
							distinguishable[j][i] = true
							changed = true
							break
						}
					}
				}
			}
		}
	}

	// Find equivalent classes
	processed := make([]bool, n)
	equivalenceClasses := [][]int{}

	for i := range n {
		if processed[i] {
			continue
		}

		class := []int{i}
		processed[i] = true

		for j := i + 1; j < n; j++ {
			if !processed[j] && !distinguishable[i][j] {
				class = append(class, j)
				processed[j] = true
			}
		}

		equivalenceClasses = append(equivalenceClasses, class)
	}

	// Build minimized DFA
	newStates := make([]string, len(equivalenceClasses))
	oldToNew := make(map[int]int)

	for i, class := range equivalenceClasses {
		newStates[i] = fmt.Sprintf("q%d", i)
		for _, oldIdx := range class {
			oldToNew[oldIdx] = i
		}
	}

	// Build transitions for minimized DFA
	newTransitions := make([][]any, len(newStates))
	for i, class := range equivalenceClasses {
		representative := class[0] // Use first state as representative
		row := make([]any, len(dfa.Alphabet))

		for j := range dfa.Alphabet {
			next := getNextState(dfa, dfa.States[representative], j)
			nextIdx := getStateIndex(dfa, next)

			if nextIdx == -1 {
				row[j] = "@v"
			} else {
				row[j] = newStates[oldToNew[nextIdx]]
			}
		}
		newTransitions[i] = row
	}

	// Determine new initial and acceptance states
	initialIdx := getStateIndex(dfa, dfa.Initial)
	newInitial := newStates[oldToNew[initialIdx]]

	newAcceptance := []string{}
	for i, class := range equivalenceClasses {
		for _, oldIdx := range class {
			if Contains(dfa.Acceptance, dfa.States[oldIdx]) {
				newAcceptance = append(newAcceptance, newStates[i])
				break
			}
		}
	}

	return &FA{
		Alphabet:    dfa.Alphabet,
		States:      newStates,
		Initial:     newInitial,
		Acceptance:  newAcceptance,
		Transitions: newTransitions,
	}, nil
}
