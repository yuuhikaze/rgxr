package logic

import (
	"fmt"
	"log"
	"slices"
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

// PerformBoolean applies a boolean operation to multiple FAs and returns the resulting FA.
func PerformBoolean(fas []*FA, mode BooleanMode) (*FA, error) {
	if len(fas) < 2 {
		return nil, fmt.Errorf("need at least two FAs for %s", mode)
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

	// Acceptance states - FIXED LOGIC
	newAcceptance := []string{}
	for _, state := range newStates {
		parts := strings.Split(state, "|")
		if len(parts) == len(fas) {
			switch mode {
			case Intersection:
				// For intersection: accept if ALL component states are accepting
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
				// For union: accept if ANY component state is accepting
				any := false
				for i, part := range parts {
					if Contains(fas[i].Acceptance, part) {
						any = true
						break
					}
				}
				if any {
					newAcceptance = append(newAcceptance, state)
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

// NPerformBoolean applies a non-deterministic boolean operation to multiple FAs and returns the resulting FA.
func NPerformBoolean(fas []*FA, mode BooleanMode) (*FA, error) {
	if len(fas) < 2 {
		return nil, fmt.Errorf("need at least two FAs for %s", mode)
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

	var newStates []string
	var newInitial string
	var newAcceptance []string
	var newTransitions [][]any

	switch mode {
	case Union:
		// Create a new initial state
		newInitial = "S"
		newStates = []string{newInitial}

		// Add all states from all FAs without renaming
		for _, fa := range fas {
			newStates = append(newStates, fa.States...)
		}

		// Find epsilon symbol index or add it if it doesn't exist
		epsilonIdx := -1
		for i, symbol := range baseAlphabet {
			if symbol == "@e" {
				epsilonIdx = i
				break
			}
		}

		// If epsilon doesn't exist, add it to the alphabet
		if epsilonIdx == -1 {
			baseAlphabet = append(baseAlphabet, "@e")
			epsilonIdx = len(baseAlphabet) - 1
		}

		// Build transitions
		newTransitions = make([][]any, len(newStates))

		// Initial state transitions: epsilon to all FA initial states
		initialRow := make([]any, len(baseAlphabet))
		for j := range baseAlphabet {
			if j == epsilonIdx {
				// Epsilon transitions to all FA initial states
				epsilonTargets := make([]string, len(fas))
				for faIdx, fa := range fas {
					epsilonTargets[faIdx] = fa.Initial
				}
				initialRow[j] = epsilonTargets
			} else {
				initialRow[j] = "@v"
			}
		}
		newTransitions[0] = initialRow

		// Copy transitions from all FAs
		stateIndex := 1
		for _, fa := range fas {
			for i := range fa.States {
				row := make([]any, len(baseAlphabet))
				for j := range baseAlphabet {
					if j < len(fa.Transitions[i]) {
						originalNext := fa.Transitions[i][j]
						row[j] = originalNext
					} else {
						// New epsilon column for existing FAs
						row[j] = "@v"
					}
				}
				newTransitions[stateIndex] = row
				stateIndex++
			}
		}

		// Acceptance states are all acceptance states from all FAs
		for _, fa := range fas {
			newAcceptance = append(newAcceptance, fa.Acceptance...)
		}
	}

	log.Print(&FA{
		Alphabet:    baseAlphabet,
		States:      newStates,
		Initial:     newInitial,
		Acceptance:  newAcceptance,
		Transitions: newTransitions,
	})

	return &FA{
		Alphabet:    baseAlphabet,
		States:      newStates,
		Initial:     newInitial,
		Acceptance:  newAcceptance,
		Transitions: newTransitions,
	}, nil
}

// Concatenation creates an NFA representing the concatenation of multiple NFAs
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

	// Collect all states from all FAs without renaming
	var newStates []string
	for _, fa := range fas {
		newStates = append(newStates, fa.States...)
	}

	// Find epsilon symbol index or add it if it doesn't exist
	epsilonIdx := -1
	for i, symbol := range baseAlphabet {
		if symbol == "@e" {
			epsilonIdx = i
			break
		}
	}

	// If epsilon doesn't exist, add it to the alphabet
	if epsilonIdx == -1 {
		baseAlphabet = append(baseAlphabet, "@e")
		epsilonIdx = len(baseAlphabet) - 1
	}

	// Build transitions
	newTransitions := make([][]any, len(newStates))
	stateIndex := 0

	for faIdx, fa := range fas {
		for i, state := range fa.States {
			row := make([]any, len(baseAlphabet))

			for j := range baseAlphabet {
				if j < len(fa.Transitions[i]) {
					originalNext := fa.Transitions[i][j]
					row[j] = originalNext
				} else {
					// New epsilon column for existing FAs
					row[j] = "@v"
				}
			}

			// If this is an accepting state and not the last FA,
			// add epsilon transition to next FA's initial state
			isAccepting := Contains(fa.Acceptance, state)
			isLastFA := faIdx == len(fas)-1

			if isAccepting && !isLastFA && epsilonIdx >= 0 {
				nextFAInitial := fas[faIdx+1].Initial

				// Handle existing epsilon transitions
				if row[epsilonIdx] == "@v" {
					row[epsilonIdx] = nextFAInitial
				} else {
					// Combine with existing epsilon transitions
					switch v := row[epsilonIdx].(type) {
					case string:
						if v != nextFAInitial {
							row[epsilonIdx] = []string{v, nextFAInitial}
						}
					case []string:
						found := slices.Contains(v, nextFAInitial)
						if !found {
							row[epsilonIdx] = append(v, nextFAInitial)
						}
					}
				}
			}

			newTransitions[stateIndex] = row
			stateIndex++
		}
	}

	// Initial state is the first FA's initial state
	newInitial := fas[0].Initial

	// Acceptance states are only from the last FA
	lastFAIdx := len(fas) - 1
	newAcceptance := append([]string{}, fas[lastFAIdx].Acceptance...)

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
	// Find epsilon symbol index
	epsilonIdx := -1
	for i, symbol := range nfa.Alphabet {
		if symbol == "@e" {
			epsilonIdx = i
			break
		}
	}

	// Helper function to compute epsilon closure
	epsilonClosure := func(states []string) []string {
		closure := make(map[string]bool)
		stack := make([]string, len(states))
		copy(stack, states)

		for _, s := range states {
			closure[s] = true
		}

		for len(stack) > 0 {
			current := stack[len(stack)-1]
			stack = stack[:len(stack)-1]

			// Find epsilon transitions if epsilon symbol exists
			if epsilonIdx >= 0 {
				stateIdx := -1
				for i, s := range nfa.States {
					if s == current {
						stateIdx = i
						break
					}
				}

				if stateIdx >= 0 && stateIdx < len(nfa.Transitions) &&
					epsilonIdx < len(nfa.Transitions[stateIdx]) {
					next := nfa.Transitions[stateIdx][epsilonIdx]

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

	// Create DFA alphabet without epsilon
	dfaAlphabet := make([]string, 0)
	for i, symbol := range nfa.Alphabet {
		if i != epsilonIdx {
			dfaAlphabet = append(dfaAlphabet, symbol)
		}
	}

	// Build DFA states and transitions
	dfaTransitions := [][]any{}
	queue := [][]string{initialClosure}
	processed := make(map[string]bool)
	needsTrashState := false

	for len(queue) > 0 {
		currentSet := queue[0]
		queue = queue[1:]

		currentKey := strings.Join(currentSet, ",")
		if processed[currentKey] {
			continue
		}
		processed[currentKey] = true

		row := make([]any, len(dfaAlphabet))

		for dfaSymbolIdx, symbol := range dfaAlphabet {
			// Find original symbol index in NFA
			nfaSymbolIdx := -1
			for i, nfaSymbol := range nfa.Alphabet {
				if nfaSymbol == symbol {
					nfaSymbolIdx = i
					break
				}
			}

			if nfaSymbolIdx == -1 {
				row[dfaSymbolIdx] = "@t"
				needsTrashState = true
				continue
			}

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
					nfaSymbolIdx < len(nfa.Transitions[stateIdx]) {
					next := nfa.Transitions[stateIdx][nfaSymbolIdx]

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

			if len(nextStates) == 0 {
				row[dfaSymbolIdx] = "@t"
				needsTrashState = true
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

				row[dfaSymbolIdx] = nextKey
			}
		}
		dfaTransitions = append(dfaTransitions, row)
	}

	// Convert state sets to state names and build final DFA
	dfaStateNames := make([]string, len(dfaStates))
	for i := range dfaStates {
		dfaStateNames[i] = fmt.Sprintf("q%d", i)
	}

	// Add trash state if needed
	if needsTrashState {
		dfaStateNames = append(dfaStateNames, "@t")
		// Create transitions for trash state - all transitions go to itself
		trashRow := make([]any, len(dfaAlphabet))
		for j := range trashRow {
			trashRow[j] = "@t"
		}
		dfaTransitions = append(dfaTransitions, trashRow)
	}

	// Update transitions to use new state names
	finalTransitions := make([][]any, len(dfaTransitions))
	for i, row := range dfaTransitions {
		finalRow := make([]any, len(row))
		for j, transition := range row {
			if transition == "@t" {
				finalRow[j] = "@t"
			} else if stateKey, ok := transition.(string); ok {
				if idx, exists := stateToIndex[stateKey]; exists {
					finalRow[j] = dfaStateNames[idx]
				} else {
					finalRow[j] = "@t"
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
		Alphabet:    dfaAlphabet,
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

// MinimizeDFA minimizes a DFA using Hopcroft's algorithm
func MinimizeDFA(dfa *FA) (*FA, error) {
	n := len(dfa.States)
	if n <= 1 {
		return dfa, nil
	}

	log.Print(dfa)
	log.Print("==============")

	// First, remove inaccessible states
	accessible := make(map[string]bool)
	queue := []string{dfa.Initial}
	accessible[dfa.Initial] = true

	// BFS to find all accessible states
	for len(queue) > 0 {
		current := queue[0]
		queue = queue[1:]

		currentIdx := getStateIndex(dfa, current)
		if currentIdx == -1 || currentIdx >= len(dfa.Transitions) {
			continue
		}

		for _, transition := range dfa.Transitions[currentIdx] {
			switch v := transition.(type) {
			case string:
				if v != "@v" && !accessible[v] {
					accessible[v] = true
					queue = append(queue, v)
				}
			case []string:
				for _, state := range v {
					if state != "@v" && !accessible[state] {
						accessible[state] = true
						queue = append(queue, state)
					}
				}
			}
		}
	}

	// Create new DFA with only accessible states
	accessibleStates := make([]string, 0)
	oldToAccessible := make(map[int]int)
	
	for i, state := range dfa.States {
		if accessible[state] {
			oldToAccessible[i] = len(accessibleStates)
			accessibleStates = append(accessibleStates, state)
		}
	}

	// Build transitions for accessible states only
	accessibleTransitions := make([][]any, len(accessibleStates))
	for i, state := range accessibleStates {
		oldIdx := getStateIndex(dfa, state)
		if oldIdx == -1 || oldIdx >= len(dfa.Transitions) {
			continue
		}
		
		row := make([]any, len(dfa.Alphabet))
		for j, transition := range dfa.Transitions[oldIdx] {
			switch v := transition.(type) {
			case string:
				if v == "@v" || !accessible[v] {
					row[j] = "@v"
				} else {
					row[j] = v
				}
			case []string:
				validStates := make([]string, 0)
				for _, state := range v {
					if state != "@v" && accessible[state] {
						validStates = append(validStates, state)
					}
				}
				if len(validStates) == 0 {
					row[j] = "@v"
				} else if len(validStates) == 1 {
					row[j] = validStates[0]
				} else {
					row[j] = validStates
				}
			default:
				row[j] = "@v"
			}
		}
		accessibleTransitions[i] = row
	}

	// Build accessible acceptance states
	accessibleAcceptance := make([]string, 0)
	for _, state := range dfa.Acceptance {
		if accessible[state] {
			accessibleAcceptance = append(accessibleAcceptance, state)
		}
	}

	// Create temporary DFA with only accessible states
	accessibleDFA := &FA{
		Alphabet:    dfa.Alphabet,
		States:      accessibleStates,
		Initial:     dfa.Initial,
		Acceptance:  accessibleAcceptance,
		Transitions: accessibleTransitions,
	}

	// Now proceed with Hopcroft's algorithm on accessible states
	n = len(accessibleDFA.States)
	if n <= 1 {
		return accessibleDFA, nil
	}

	// Partition states into accepting and non-accepting
	accepting := make([]int, 0)
	nonAccepting := make([]int, 0)

	for i, state := range accessibleDFA.States {
		if Contains(accessibleDFA.Acceptance, state) {
			accepting = append(accepting, i)
		} else {
			nonAccepting = append(nonAccepting, i)
		}
	}

	// Initial partition
	partition := make([][]int, 0)
	if len(nonAccepting) > 0 {
		partition = append(partition, nonAccepting)
	}
	if len(accepting) > 0 {
		partition = append(partition, accepting)
	}

	// Create work list with all alphabet symbols for each partition
	workList := make([][]int, 0)
	for _, block := range partition {
		for range accessibleDFA.Alphabet {
			workList = append(workList, block)
		}
	}

	// Hopcroft's algorithm main loop
	for len(workList) > 0 {
		// Pop from work list
		splitter := workList[0]
		workList = workList[1:]

		for symbolIdx := range accessibleDFA.Alphabet {
			// Find states that transition to splitter on this symbol
			predecessors := make([]int, 0)
			for i, state := range accessibleDFA.States {
				next := getNextState(accessibleDFA, state, symbolIdx)
				nextIdx := getStateIndex(accessibleDFA, next)
				if nextIdx != -1 {
					if slices.Contains(splitter, nextIdx) {
						predecessors = append(predecessors, i)
					}
				}
			}

			if len(predecessors) == 0 {
				continue
			}

			// Split blocks that intersect with predecessors
			newPartition := make([][]int, 0)
			for _, block := range partition {
				intersect := make([]int, 0)
				difference := make([]int, 0)

				for _, state := range block {
					found := slices.Contains(predecessors, state)
					if found {
						intersect = append(intersect, state)
					} else {
						difference = append(difference, state)
					}
				}

				if len(intersect) > 0 && len(difference) > 0 {
					// Block was split
					newPartition = append(newPartition, intersect)
					newPartition = append(newPartition, difference)

					// Update work list
					newWorkList := make([][]int, 0)
					for _, workItem := range workList {
						if slicesEqual(workItem, block) {
							// Replace with smaller block
							if len(intersect) <= len(difference) {
								newWorkList = append(newWorkList, intersect)
							} else {
								newWorkList = append(newWorkList, difference)
							}
						} else {
							newWorkList = append(newWorkList, workItem)
						}
					}
					workList = newWorkList

					// Add new block to work list for all symbols
					var smallerBlock []int
					if len(intersect) <= len(difference) {
						smallerBlock = intersect
					} else {
						smallerBlock = difference
					}
					for range accessibleDFA.Alphabet {
						workList = append(workList, smallerBlock)
					}
				} else {
					newPartition = append(newPartition, block)
				}
			}
			partition = newPartition
		}
	}

	// Build minimized DFA from partition
	newStates := make([]string, len(partition))
	oldToNew := make(map[int]int)

	for i, block := range partition {
		newStates[i] = fmt.Sprintf("q%d", i)
		for _, oldIdx := range block {
			oldToNew[oldIdx] = i
		}
	}

	// Build transitions
	newTransitions := make([][]any, len(newStates))
	for i, block := range partition {
		representative := block[0]
		row := make([]any, len(accessibleDFA.Alphabet))

		for j := range accessibleDFA.Alphabet {
			next := getNextState(accessibleDFA, accessibleDFA.States[representative], j)
			nextIdx := getStateIndex(accessibleDFA, next)

			if nextIdx == -1 {
				row[j] = "@v"
			} else {
				row[j] = newStates[oldToNew[nextIdx]]
			}
		}
		newTransitions[i] = row
	}

	// Determine initial and acceptance states
	initialIdx := getStateIndex(accessibleDFA, accessibleDFA.Initial)
	newInitial := newStates[oldToNew[initialIdx]]

	newAcceptance := []string{}
	for i, block := range partition {
		for _, oldIdx := range block {
			if Contains(accessibleDFA.Acceptance, accessibleDFA.States[oldIdx]) {
				newAcceptance = append(newAcceptance, newStates[i])
				break
			}
		}
	}

	log.Print(&FA{
		Alphabet:    accessibleDFA.Alphabet,
		States:      newStates,
		Initial:     newInitial,
		Acceptance:  newAcceptance,
		Transitions: newTransitions,
	})

	return &FA{
		Alphabet:    accessibleDFA.Alphabet,
		States:      newStates,
		Initial:     newInitial,
		Acceptance:  newAcceptance,
		Transitions: newTransitions,
	}, nil
}

// slicesEqual checks if two int slices are equal
func slicesEqual(a, b []int) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
