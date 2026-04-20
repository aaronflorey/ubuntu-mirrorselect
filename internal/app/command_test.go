package app

import (
	"bytes"
	"errors"
	"net/url"
	"strings"
	"testing"
)

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

func TestValidateConfigOutput(t *testing.T) {
	t.Parallel()

	base := cliConfig{
		Arch:        "amd64",
		Country:     "US",
		Interactive: false,
		Max:         5,
		Protocol:    "https",
		Release:     "questing",
		Timeout:     500,
		Verbosity:   "WARN",
	}

	valid := base
	valid.Output = "json"
	if err := validateConfig(&valid); err != nil {
		t.Fatalf("validateConfig(valid json output) returned error: %v", err)
	}

	invalid := base
	invalid.Output = "yaml"
	if err := validateConfig(&invalid); err == nil {
		t.Fatal("validateConfig(invalid output) expected error, got nil")
	}

	interactiveJSON := base
	interactiveJSON.Output = "json"
	interactiveJSON.Interactive = true
	if err := validateConfig(&interactiveJSON); err == nil {
		t.Fatal("validateConfig(interactive json output) expected error, got nil")
	}

	applyJSON := base
	applyJSON.Output = "json"
	applyJSON.Apply = true
	if err := validateConfig(&applyJSON); err == nil {
		t.Fatal("validateConfig(apply json output) expected error, got nil")
	}
}

func TestLatencyWorkerCount(t *testing.T) {
	t.Parallel()

	if got := latencyWorkerCount(0); got != 0 {
		t.Fatalf("latencyWorkerCount(0)=%d, want 0", got)
	}

	if got := latencyWorkerCount(3); got != 3 {
		t.Fatalf("latencyWorkerCount(3)=%d, want 3", got)
	}

	got := latencyWorkerCount(100)
	if got < 1 || got > 32 {
		t.Fatalf("latencyWorkerCount(100)=%d, want between 1 and 32", got)
	}
}

func TestSelectMirrorInteractively(t *testing.T) {
	t.Parallel()

	mirrors := []Mirror{
		{URL: mustURL(t, "https://mirror1.example/ubuntu/"), Latency: 100, Size: 1_048_576, Time: 1},
		{URL: mustURL(t, "https://mirror2.example/ubuntu/"), Latency: 80, Size: 2_097_152, Time: 1},
	}

	t.Run("select valid option", func(t *testing.T) {
		t.Parallel()

		input := strings.NewReader("2\n")
		var output bytes.Buffer

		selected, err := selectMirrorInteractively(mirrors, input, &output)
		if err != nil {
			t.Fatalf("selectMirrorInteractively returned error: %v", err)
		}
		if selected.URL.String() != "https://mirror2.example/ubuntu/" {
			t.Fatalf("selected URL=%q, want %q", selected.URL.String(), "https://mirror2.example/ubuntu/")
		}
	})

	t.Run("retry after invalid option", func(t *testing.T) {
		t.Parallel()

		input := strings.NewReader("9\n1\n")
		var output bytes.Buffer

		selected, err := selectMirrorInteractively(mirrors, input, &output)
		if err != nil {
			t.Fatalf("selectMirrorInteractively returned error: %v", err)
		}
		if selected.URL.String() != "https://mirror1.example/ubuntu/" {
			t.Fatalf("selected URL=%q, want %q", selected.URL.String(), "https://mirror1.example/ubuntu/")
		}
		if !strings.Contains(output.String(), "Invalid selection") {
			t.Fatalf("expected output to contain invalid selection message, got %q", output.String())
		}
	})

	t.Run("cancel selection", func(t *testing.T) {
		t.Parallel()

		input := strings.NewReader("q\n")
		var output bytes.Buffer

		_, err := selectMirrorInteractively(mirrors, input, &output)
		if err == nil {
			t.Fatal("selectMirrorInteractively expected cancel error, got nil")
		}
	})
}

func TestSudoReexecExitCode(t *testing.T) {
	t.Parallel()

	if code, ok := SudoReexecExitCode(&sudoReexecResult{exitCode: 7}); !ok || code != 7 {
		t.Fatalf("SudoReexecExitCode returned (%d, %t), want (7, true)", code, ok)
	}

	if code, ok := SudoReexecExitCode(errors.New("other error")); ok || code != 0 {
		t.Fatalf("SudoReexecExitCode returned (%d, %t), want (0, false)", code, ok)
	}
}

func mustURL(t *testing.T, raw string) *url.URL {
	t.Helper()

	parsed, err := url.Parse(raw)
	if err != nil {
		t.Fatalf("url.Parse(%q) returned error: %v", raw, err)
	}

	return parsed
}
