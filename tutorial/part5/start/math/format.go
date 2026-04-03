package math

import "fmt"

// Transaction represents a financial transaction for anomaly detection.
// Fields are ordered by semantic weight: label and category for identification,
// then the three fields that FormatTransaction uses in sentence order.
type Transaction struct {
	Label    string  `json:"label"`
	Category string  `json:"category"`
	Amount   float64 `json:"amount"`
	Merchant string  `json:"merchant"`
	Time     string  `json:"time"`
}

// FormatTransaction builds a synthetic sentence that gives the embedding model
// full context: amount, merchant, and time of day together.
//
// Field order matters — semantically heavy fields come first (amount, merchant)
// so the model weights them more strongly. Modifiers (time) come after.
//
// Example output: "Transaction: 45.00 USD at Whole Foods Market at 11:00 AM"
func FormatTransaction(t Transaction) string {
	return fmt.Sprintf("Transaction: %.2f USD at %s at %s", t.Amount, t.Merchant, t.Time)
}

// FormatRaw returns only the dollar amount as a decimal string.
// Without merchant or time context, the model embeds a number, not a transaction.
//
// Example output: "45.00"
func FormatRaw(t Transaction) string {
	return fmt.Sprintf("%.2f", t.Amount)
}
