module embedtutorial

go 1.24.1

require (
	embedtutorial/shared v0.0.0-00010101000000-000000000000
	github.com/spf13/cobra v1.8.1
	gonum.org/v1/gonum v0.15.1
	gonum.org/v1/plot v0.14.0
)

// Use the shared embedder package
replace embedtutorial/shared => ../../shared
