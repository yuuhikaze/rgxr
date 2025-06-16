package logic

import (
	"fmt"
	"log"
	"strings"
)

// ToDot generates a DOT language string representing the FA.
func ToDot(fa FA) string {
	var b strings.Builder

	b.WriteString("digraph FA {\n")
	b.WriteString("  rankdir=LR;\n") // Left to right
	b.WriteString("  start [style=invis];\n")

	// Define nodes
	if len(fa.Acceptance) > 0 {
		b.WriteString("  node [shape=doublecircle];")
		for _, acc := range fa.Acceptance {
			b.WriteString(fmt.Sprintf(" \"%s\"", acc))
		}
		b.WriteString(";\n")
	}
	b.WriteString("  node [shape=circle];\n")

	// Transitions: for each state, for each symbol, add edges
	b.WriteString(fmt.Sprintf("  start -> \"%s\";\n", fa.Initial))
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
	log.Println(b.String()) // DEBUGLOG
	return b.String()
}
