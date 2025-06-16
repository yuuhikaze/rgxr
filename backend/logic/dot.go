package logic

import (
	"fmt"
	"strings"
)

// ToDot generates a DOT language string representing the FA.
func ToDot(fa FA) string {
	var b strings.Builder

	b.WriteString("digraph FA {\n")
	b.WriteString("  rankdir=LR;\n") // Left to right

	// Invisible start node to initial
	b.WriteString("  start [shape=point];\n")
	b.WriteString(fmt.Sprintf("  start -> \"%s\";\n", fa.Initial))

	// Define nodes
	for _, state := range fa.States {
		shape := "circle"
		peripheries := 1
		if Contains(fa.Acceptance, state) {
			peripheries = 2 // double circle for acceptance
		}
		b.WriteString(fmt.Sprintf("  \"%s\" [shape=%s peripheries=%d];\n", state, shape, peripheries))
	}

	// Transitions: for each state, for each symbol, add edges
	for i, from := range fa.States {
		for j, symbol := range fa.Alphabet {
			dest := fa.Transitions[i][j]
			if dest == "@v" {
				continue
			}
			switch v := dest.(type) {
			case string:
				b.WriteString(fmt.Sprintf("  \"%s\" -> \"%s\" [label=\"%s\"];\n", from, v, symbol))
			case []string:
				for _, d := range v {
					b.WriteString(fmt.Sprintf("  \"%s\" -> \"%s\" [label=\"%s\"];\n", from, d, symbol))
				}
			}
		}
	}

	b.WriteString("}\n")
	return b.String()
}
