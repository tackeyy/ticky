package ticktick

import "testing"

func TestParsePriority_ValidInputs(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  int
	}{
		{"none keyword", "none", PriorityNone},
		{"none numeric", "0", PriorityNone},
		{"low keyword", "low", PriorityLow},
		{"low numeric", "1", PriorityLow},
		{"medium keyword", "medium", PriorityMedium},
		{"medium abbreviation", "med", PriorityMedium},
		{"medium numeric", "3", PriorityMedium},
		{"high keyword", "high", PriorityHigh},
		{"high numeric", "5", PriorityHigh},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParsePriority(tt.input)
			if err != nil {
				t.Fatalf("ParsePriority(%q) returned unexpected error: %v", tt.input, err)
			}
			if got != tt.want {
				t.Errorf("ParsePriority(%q) = %d, want %d", tt.input, got, tt.want)
			}
		})
	}
}

func TestParsePriority_InvalidInputs(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{"empty string", ""},
		{"unsupported word", "urgent"},
		{"gap between low and medium", "2"},
		{"gap between medium and high", "4"},
		{"uppercase MEDIUM", "MEDIUM"},
		{"uppercase HIGH", "HIGH"},
		{"mixed case Low", "Low"},
		{"negative number", "-1"},
		{"large number", "100"},
		{"whitespace", " "},
		{"leading space", " low"},
		{"trailing space", "low "},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParsePriority(tt.input)
			if err == nil {
				t.Errorf("ParsePriority(%q) = %d, want error", tt.input, got)
			}
		})
	}
}

func TestPriorityString(t *testing.T) {
	tests := []struct {
		name  string
		input int
		want  string
	}{
		{"none", PriorityNone, "none"},
		{"low", PriorityLow, "low"},
		{"medium", PriorityMedium, "medium"},
		{"high", PriorityHigh, "high"},
		{"unknown 2", 2, "unknown(2)"},
		{"unknown 4", 4, "unknown(4)"},
		{"unknown 99", 99, "unknown(99)"},
		{"negative", -1, "unknown(-1)"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := PriorityString(tt.input)
			if got != tt.want {
				t.Errorf("PriorityString(%d) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}

func TestParsePriority_RoundTrip(t *testing.T) {
	// Parse a canonical name, convert back to string, verify it matches.
	tests := []struct {
		name      string
		input     string
		wantLabel string
	}{
		{"none round trip", "none", "none"},
		{"low round trip", "low", "low"},
		{"medium round trip", "medium", "medium"},
		{"high round trip", "high", "high"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			val, err := ParsePriority(tt.input)
			if err != nil {
				t.Fatalf("ParsePriority(%q) returned unexpected error: %v", tt.input, err)
			}
			got := PriorityString(val)
			if got != tt.wantLabel {
				t.Errorf("PriorityString(ParsePriority(%q)) = %q, want %q", tt.input, got, tt.wantLabel)
			}
		})
	}
}

func TestPriorityConstants(t *testing.T) {
	tests := []struct {
		name string
		got  int
		want int
	}{
		{"PriorityNone", PriorityNone, 0},
		{"PriorityLow", PriorityLow, 1},
		{"PriorityMedium", PriorityMedium, 3},
		{"PriorityHigh", PriorityHigh, 5},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.got != tt.want {
				t.Errorf("%s = %d, want %d", tt.name, tt.got, tt.want)
			}
		})
	}
}
