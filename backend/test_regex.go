package main

import (
	"encoding/json"
	"fmt"

	"github.com/yuuhikaze/rgxr/logic"
)

func main() {
	testCases := []string{
		"a",
		"ab", 
		"a∪b",
		"a*",
		"a+",
		"(a∪b)*",
		"ε",
		"∅",
		// Additional complex test cases
		"(0∪1)*000(0∪1)*",
		"(((00)*(11))∪01)*",
		"∅*",
		"a(abb)*∪b",
		"a+∪(ab)+",
		"(a∪b+)a+b+",
	}

	for _, regex := range testCases {
		fmt.Printf("\n=== Testing regex: %s ===\n", regex)
		
		fa, err := logic.RegexToNFA(regex)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			continue
		}

		// Print the FA structure
		faJSON, err := json.MarshalIndent(fa, "", "  ")
		if err != nil {
			fmt.Printf("Error marshaling FA: %v\n", err)
			continue
		}
		
		fmt.Printf("Generated NFA:\n%s\n", string(faJSON))
	}
}
