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
			if i >= len(fa.Transitions) || j >= len(fa.Transitions[i]) {
				continue
			}
			
			dest := fa.Transitions[i][j]
			
			// Handle different types of destinations
			switch v := dest.(type) {
			case string:
				if v != "@v" && v != "" {
					// Escape special characters in labels
					label := escapeLabel(symbol)
					b.WriteString(fmt.Sprintf("  \"%s\" -> \"%s\" [label=\"%s\"];\n", from, v, label))
				}
			case []string:
				for _, d := range v {
					if d != "@v" && d != "" {
						label := escapeLabel(symbol)
						b.WriteString(fmt.Sprintf("  \"%s\" -> \"%s\" [label=\"%s\"];\n", from, d, label))
					}
				}
			case []any:
				// Handle []any type which can contain strings
				for _, item := range v {
					if s, ok := item.(string); ok && s != "@v" && s != "" {
						label := escapeLabel(symbol)
						b.WriteString(fmt.Sprintf("  \"%s\" -> \"%s\" [label=\"%s\"];\n", from, s, label))
					}
				}
			}
		}
	}

	b.WriteString("}\n")
	return b.String()
}

// escapeLabel escapes special characters in DOT labels
func escapeLabel(label string) string {
	// Handle epsilon symbol
	// if label == "@e" {
	// 	return "$\\varepsilon$"
	// }
	
	// Escape quotes and backslashes
	label = strings.ReplaceAll(label, "\\", "\\\\")
	label = strings.ReplaceAll(label, "\"", "\\\"")
	
	return label
}
