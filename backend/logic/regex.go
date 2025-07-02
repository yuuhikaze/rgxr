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

	// Parse and convert regex to NFA using Thompson's construction
	parser := &RegexParser{
		input:        regex,
		runes:        []rune(regex),
		pos:          0,
		stateCounter: 0,
	}
	fragment, err := parser.parseExpression()
	if err != nil {
		return nil, err
	}
	
	if parser.pos < len(parser.runes) {
		return nil, fmt.Errorf("unexpected character at position %d", parser.pos)
	}

	return fragment.toFA(), nil
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
	start       string
	end         string
	states      map[string]bool
	transitions map[string]map[string][]string
	alphabet    map[string]bool
}

// Convert NFAFragment to FA
func (f *NFAFragment) toFA() *FA {
	// Handle empty language case
	if f.end == "" {
		return &FA{
			Alphabet:    []string{},
			States:      []string{f.start},
			Initial:     f.start,
			Acceptance:  []string{}, // No accepting states for empty language
			Transitions: [][]any{{"@v"}},
		}
	}

	// Build alphabet slice (regular symbols first, then epsilon)
	alphabetSlice := make([]string, 0, len(f.alphabet))
	for symbol := range f.alphabet {
		if symbol != "@e" {
			alphabetSlice = append(alphabetSlice, symbol)
		}
	}
	// Add epsilon if it exists in transitions
	hasEpsilon := false
	for symbol := range f.alphabet {
		if symbol == "@e" {
			hasEpsilon = true
			break
		}
	}
	if hasEpsilon {
		alphabetSlice = append(alphabetSlice, "@e")
	}

	// Build states slice in sorted order for consistency
	statesSlice := make([]string, 0, len(f.states))
	for state := range f.states {
		statesSlice = append(statesSlice, state)
	}
	// Simple sort by state name
	for i := 0; i < len(statesSlice); i++ {
		for j := i + 1; j < len(statesSlice); j++ {
			if statesSlice[i] > statesSlice[j] {
				statesSlice[i], statesSlice[j] = statesSlice[j], statesSlice[i]
			}
		}
	}

	// Build transitions matrix
	transitions := make([][]any, len(statesSlice))
	for i, state := range statesSlice {
		row := make([]any, len(alphabetSlice))
		for j, symbol := range alphabetSlice {
			if stateTransitions, exists := f.transitions[state]; exists {
				if nextStates, exists := stateTransitions[symbol]; exists && len(nextStates) > 0 {
					if len(nextStates) == 1 {
						row[j] = nextStates[0]
					} else {
						row[j] = nextStates
					}
				} else {
					row[j] = "@v"
				}
			} else {
				row[j] = "@v"
			}
		}
		transitions[i] = row
	}

	return &FA{
		Alphabet:    alphabetSlice,
		States:      statesSlice,
		Initial:     f.start,
		Acceptance:  []string{f.end},
		Transitions: transitions,
	}
}

// RegexParser implements a recursive descent parser for regular expressions
type RegexParser struct {
	input        string
	runes        []rune
	pos          int
	stateCounter int
}

// Generate unique state names
func (p *RegexParser) newState() string {
	state := fmt.Sprintf("q%d", p.stateCounter)
	p.stateCounter++
	return state
}

// Peek at current character without advancing position
func (p *RegexParser) peek() rune {
	if p.pos >= len(p.runes) {
		return 0
	}
	return p.runes[p.pos]
}

// Advance position and return current character
func (p *RegexParser) advance() rune {
	if p.pos >= len(p.runes) {
		return 0
	}
	ch := p.runes[p.pos]
	p.pos++
	return ch
}

// Parse expression (handles union with lowest precedence)
func (p *RegexParser) parseExpression() (*NFAFragment, error) {
	left, err := p.parseSequence()
	if err != nil {
		return nil, err
	}

	for p.peek() == '∪' {
		p.advance() // consume '∪'
		right, err := p.parseSequence()
		if err != nil {
			return nil, err
		}
		left = p.union(left, right)
	}

	return left, nil
}

// Parse sequence (handles concatenation)
func (p *RegexParser) parseSequence() (*NFAFragment, error) {
	fragments := []*NFAFragment{}

	for p.pos < len(p.runes) && p.peek() != ')' && p.peek() != '∪' {
		fragment, err := p.parseFactor()
		if err != nil {
			return nil, err
		}
		fragments = append(fragments, fragment)
	}

	if len(fragments) == 0 {
		return p.epsilon(), nil
	}

	result := fragments[0]
	for i := 1; i < len(fragments); i++ {
		result = p.concatenate(result, fragments[i])
	}

	return result, nil
}

// Parse factor (handles Kleene star and plus)
func (p *RegexParser) parseFactor() (*NFAFragment, error) {
	base, err := p.parseAtom()
	if err != nil {
		return nil, err
	}

	for p.peek() == '*' || p.peek() == '∗' || p.peek() == '+' {
		op := p.advance()
		if op == '*' || op == '∗' {
			base = p.kleeneStar(base)
		} else if op == '+' {
			base = p.kleenePlus(base)
		}
	}

	return base, nil
}

// Parse atom (basic elements)
func (p *RegexParser) parseAtom() (*NFAFragment, error) {
	if p.pos >= len(p.runes) {
		return nil, fmt.Errorf("unexpected end of input")
	}

	ch := p.runes[p.pos]

	if ch == '(' {
		p.pos++
		expr, err := p.parseExpression()
		if err != nil {
			return nil, err
		}
		if p.pos >= len(p.runes) || p.runes[p.pos] != ')' {
			return nil, fmt.Errorf("expected closing parenthesis")
		}
		p.pos++
		return expr, nil
	}

	if ch == 'ε' {
		p.pos++
		return p.epsilon(), nil
	}

	if ch == '∅' {
		p.pos++
		return p.empty(), nil
	}

	// Regular character
	p.pos++
	return p.character(string(ch)), nil
}

// Thompson construction primitives

// Create epsilon NFA fragment
func (p *RegexParser) epsilon() *NFAFragment {
	start := p.newState()
	end := p.newState()

	states := map[string]bool{start: true, end: true}
	transitions := map[string]map[string][]string{
		start: {"@e": {end}},
		end:   {},
	}
	alphabet := map[string]bool{"@e": true}

	return &NFAFragment{
		start:       start,
		end:         end,
		states:      states,
		transitions: transitions,
		alphabet:    alphabet,
	}
}

// Create empty language NFA fragment
func (p *RegexParser) empty() *NFAFragment {
	start := p.newState()

	states := map[string]bool{start: true}
	transitions := map[string]map[string][]string{start: {}}
	alphabet := map[string]bool{}

	return &NFAFragment{
		start:       start,
		end:         "", // No accepting state for empty language
		states:      states,
		transitions: transitions,
		alphabet:    alphabet,
	}
}

// Create character NFA fragment
func (p *RegexParser) character(char string) *NFAFragment {
	start := p.newState()
	end := p.newState()

	states := map[string]bool{start: true, end: true}
	transitions := map[string]map[string][]string{
		start: {char: {end}},
		end:   {},
	}
	alphabet := map[string]bool{char: true}

	return &NFAFragment{
		start:       start,
		end:         end,
		states:      states,
		transitions: transitions,
		alphabet:    alphabet,
	}
}

// Union operation
func (p *RegexParser) union(left, right *NFAFragment) *NFAFragment {
	start := p.newState()
	end := p.newState()

	// Combine states
	states := map[string]bool{start: true, end: true}
	for state := range left.states {
		states[state] = true
	}
	for state := range right.states {
		states[state] = true
	}

	// Combine alphabet
	alphabet := map[string]bool{"@e": true}
	for symbol := range left.alphabet {
		alphabet[symbol] = true
	}
	for symbol := range right.alphabet {
		alphabet[symbol] = true
	}

	// Combine transitions
	transitions := map[string]map[string][]string{
		start: {"@e": {left.start, right.start}},
		end:   {},
	}

	// Copy left transitions
	for state, stateTransitions := range left.transitions {
		if transitions[state] == nil {
			transitions[state] = make(map[string][]string)
		}
		for symbol, nextStates := range stateTransitions {
			transitions[state][symbol] = append(transitions[state][symbol], nextStates...)
		}
	}

	// Copy right transitions
	for state, stateTransitions := range right.transitions {
		if transitions[state] == nil {
			transitions[state] = make(map[string][]string)
		}
		for symbol, nextStates := range stateTransitions {
			transitions[state][symbol] = append(transitions[state][symbol], nextStates...)
		}
	}

	// Add epsilon transitions from old end states to new end state
	if transitions[left.end] == nil {
		transitions[left.end] = make(map[string][]string)
	}
	transitions[left.end]["@e"] = append(transitions[left.end]["@e"], end)

	if transitions[right.end] == nil {
		transitions[right.end] = make(map[string][]string)
	}
	transitions[right.end]["@e"] = append(transitions[right.end]["@e"], end)

	return &NFAFragment{
		start:       start,
		end:         end,
		states:      states,
		transitions: transitions,
		alphabet:    alphabet,
	}
}

// Concatenation operation
func (p *RegexParser) concatenate(left, right *NFAFragment) *NFAFragment {
	// Combine states
	states := map[string]bool{}
	for state := range left.states {
		states[state] = true
	}
	for state := range right.states {
		states[state] = true
	}

	// Combine alphabet
	alphabet := map[string]bool{"@e": true}
	for symbol := range left.alphabet {
		alphabet[symbol] = true
	}
	for symbol := range right.alphabet {
		alphabet[symbol] = true
	}

	// Combine transitions
	transitions := map[string]map[string][]string{}

	// Copy left transitions
	for state, stateTransitions := range left.transitions {
		if transitions[state] == nil {
			transitions[state] = make(map[string][]string)
		}
		for symbol, nextStates := range stateTransitions {
			transitions[state][symbol] = append(transitions[state][symbol], nextStates...)
		}
	}

	// Copy right transitions
	for state, stateTransitions := range right.transitions {
		if transitions[state] == nil {
			transitions[state] = make(map[string][]string)
		}
		for symbol, nextStates := range stateTransitions {
			transitions[state][symbol] = append(transitions[state][symbol], nextStates...)
		}
	}

	// Add epsilon transition from left end to right start
	if transitions[left.end] == nil {
		transitions[left.end] = make(map[string][]string)
	}
	transitions[left.end]["@e"] = append(transitions[left.end]["@e"], right.start)

	return &NFAFragment{
		start:       left.start,
		end:         right.end,
		states:      states,
		transitions: transitions,
		alphabet:    alphabet,
	}
}

// Kleene star operation
func (p *RegexParser) kleeneStar(fragment *NFAFragment) *NFAFragment {
	start := p.newState()
	end := p.newState()

	// Combine states
	states := map[string]bool{start: true, end: true}
	for state := range fragment.states {
		states[state] = true
	}

	// Combine alphabet
	alphabet := map[string]bool{"@e": true}
	for symbol := range fragment.alphabet {
		alphabet[symbol] = true
	}

	// Copy transitions
	transitions := map[string]map[string][]string{
		start: {"@e": {fragment.start, end}}, // Can go to fragment or skip it
		end:   {},
	}

	// Copy fragment transitions
	for state, stateTransitions := range fragment.transitions {
		if transitions[state] == nil {
			transitions[state] = make(map[string][]string)
		}
		for symbol, nextStates := range stateTransitions {
			transitions[state][symbol] = append(transitions[state][symbol], nextStates...)
		}
	}

	// Add epsilon transition from fragment end back to fragment start and to new end
	if transitions[fragment.end] == nil {
		transitions[fragment.end] = make(map[string][]string)
	}
	transitions[fragment.end]["@e"] = append(transitions[fragment.end]["@e"], fragment.start, end)

	return &NFAFragment{
		start:       start,
		end:         end,
		states:      states,
		transitions: transitions,
		alphabet:    alphabet,
	}
}

// Kleene plus operation (a+ = aa*)
func (p *RegexParser) kleenePlus(fragment *NFAFragment) *NFAFragment {
	star := p.kleeneStar(fragment)
	return p.concatenate(fragment, star)
}
