package app

import "testing"

func TestNormalizeCountryCode(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		input string
		want  string
	}{
		{name: "upper stays upper", input: "US", want: "US"},
		{name: "lower converts to upper", input: "us", want: "US"},
		{name: "maps UK to GB", input: "uk", want: "GB"},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			got := normalizeCountryCode(tc.input)
			if got != tc.want {
				t.Fatalf("normalizeCountryCode(%q)=%q, want %q", tc.input, got, tc.want)
			}
		})
	}
}
