package ticktick

import "fmt"

// Priority constants matching TickTick API values.
const (
	PriorityNone   = 0
	PriorityLow    = 1
	PriorityMedium = 3
	PriorityHigh   = 5
)

// ParsePriority converts a human-readable priority string to TickTick API value.
func ParsePriority(s string) (int, error) {
	switch s {
	case "none", "0":
		return PriorityNone, nil
	case "low", "1":
		return PriorityLow, nil
	case "medium", "med", "3":
		return PriorityMedium, nil
	case "high", "5":
		return PriorityHigh, nil
	default:
		return 0, fmt.Errorf("invalid priority: %s (use none, low, medium, high)", s)
	}
}

// PriorityString converts a TickTick API priority value to human-readable string.
func PriorityString(p int) string {
	switch p {
	case PriorityNone:
		return "none"
	case PriorityLow:
		return "low"
	case PriorityMedium:
		return "medium"
	case PriorityHigh:
		return "high"
	default:
		return fmt.Sprintf("unknown(%d)", p)
	}
}
