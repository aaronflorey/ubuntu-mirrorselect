package app

import "testing"

func TestParseCountryCode(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		input      string
		wantCode   string
		wantLookup bool
	}{
		{name: "simple country", input: "United States", wantCode: "US", wantLookup: true},
		{name: "comma qualified", input: "Korea, Republic of", wantCode: "KR", wantLookup: true},
		{name: "alias fallback", input: "Viet Nam", wantCode: "VN", wantLookup: true},
		{name: "unknown country", input: "Atlantis", wantCode: "", wantLookup: false},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			gotCode, gotLookup := parseCountryCode(tc.input)
			if gotLookup != tc.wantLookup {
				t.Fatalf("parseCountryCode(%q) lookup=%t, want %t", tc.input, gotLookup, tc.wantLookup)
			}
			if gotCode != tc.wantCode {
				t.Fatalf("parseCountryCode(%q)=%q, want %q", tc.input, gotCode, tc.wantCode)
			}
		})
	}
}
