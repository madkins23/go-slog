package language

import (
	"fmt"
)

func ExampleSetup() {
	// Setup() tries all -language flag values to find one that works.
	if Setup() == nil { // No error, language configured.
		// Format a number using the language "printer".
		fmt.Println(Printer().Sprintf("%d", 1_234_567_890))
		// Output: 1,234,567,890
	}
}
