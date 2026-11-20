package errors_test

import (
	"fmt"

	"github.com/jgolang/errors/codes"
)

// duplicateCodeStrings returns each code string declared more than once.
// A consumer runs this over its own codes in a test.
func duplicateCodeStrings(declared ...codes.Coder) []string {
	seen := make(map[string]bool, len(declared))
	var duplicates []string
	for _, code := range declared {
		if seen[code.Str()] {
			duplicates = append(duplicates, code.Str())
		}
		seen[code.Str()] = true
	}
	return duplicates
}

func Example_codeUniqueness() {
	titleRequired := codes.New("PLN-001", "title is required")
	taskNotFound := codes.New("PLN-002", "task not found")
	reused := codes.New("PLN-001", "project not found")

	fmt.Println(duplicateCodeStrings(titleRequired, taskNotFound))
	fmt.Println(duplicateCodeStrings(titleRequired, taskNotFound, reused))
	// Output:
	// []
	// [PLN-001]
}
