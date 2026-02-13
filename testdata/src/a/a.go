package a

import "log"

func Example() {
	// This should trigger a warning - no arguments
	log.Print() // want "log call should have at least one argument"

	// Valid log calls
	log.Println("valid log message")
	log.Printf("formatted message: %s", "value")
	
	// Add your own test cases here based on your custom rules
}
