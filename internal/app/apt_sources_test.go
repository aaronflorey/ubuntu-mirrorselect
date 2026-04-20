package app

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestNormalizeAptMirrorURI(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		input   string
		want    string
		wantErr bool
	}{
		{name: "keeps ubuntu path", input: "https://mirror.example/ubuntu/", want: "https://mirror.example/ubuntu"},
		{name: "adds ubuntu path", input: "https://mirror.example", want: "https://mirror.example/ubuntu"},
		{name: "rejects unsupported scheme", input: "ftp://mirror.example/ubuntu", wantErr: true},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			got, err := normalizeAptMirrorURI(tc.input)
			if tc.wantErr {
				if err == nil {
					t.Fatalf("normalizeAptMirrorURI(%q) expected error, got nil", tc.input)
				}
				return
			}

			if err != nil {
				t.Fatalf("normalizeAptMirrorURI(%q) returned error: %v", tc.input, err)
			}
			if got != tc.want {
				t.Fatalf("normalizeAptMirrorURI(%q)=%q, want %q", tc.input, got, tc.want)
			}
		})
	}
}

func TestRenderDeb822SourcesIncludesSecurityAndTypes(t *testing.T) {
	t.Parallel()

	content := renderDeb822Sources("https://mirror.example/ubuntu", "noble")

	wants := []string{
		"Types: deb deb-src",
		"URIs: https://mirror.example/ubuntu",
		"Suites: noble noble-updates noble-backports noble-security",
		"Components: main restricted universe multiverse",
		"Signed-By: /usr/share/keyrings/ubuntu-archive-keyring.gpg",
	}

	for _, want := range wants {
		if !strings.Contains(content, want) {
			t.Fatalf("renderDeb822Sources missing %q\ncontent:\n%s", want, content)
		}
	}
}

func TestRenderLegacySourcesListIncludesSecurityAndTypes(t *testing.T) {
	t.Parallel()

	content := renderLegacySourcesList("https://mirror.example/ubuntu", "noble")

	wants := []string{
		"deb https://mirror.example/ubuntu noble main restricted universe multiverse",
		"deb https://mirror.example/ubuntu noble-security main restricted universe multiverse",
		"deb-src https://mirror.example/ubuntu noble main restricted universe multiverse",
		"deb-src https://mirror.example/ubuntu noble-security main restricted universe multiverse",
	}

	for _, want := range wants {
		if !strings.Contains(content, want) {
			t.Fatalf("renderLegacySourcesList missing %q\ncontent:\n%s", want, content)
		}
	}
}

func TestApplyMirrorToAPTTargetsCreatesBackupsAndWrites(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	deb822Path := filepath.Join(dir, "ubuntu.sources")
	legacyPath := filepath.Join(dir, "sources.list")

	if err := os.WriteFile(deb822Path, []byte("deb822 old\n"), 0644); err != nil {
		t.Fatalf("failed to write deb822 fixture: %v", err)
	}
	if err := os.WriteFile(legacyPath, []byte("legacy old\n"), 0644); err != nil {
		t.Fatalf("failed to write legacy fixture: %v", err)
	}

	targets := []aptSourceTarget{
		{Path: deb822Path, Format: aptSourceFormatDeb822},
		{Path: legacyPath, Format: aptSourceFormatLegacy},
	}

	result, err := applyMirrorToAPTTargets("https://mirror.example/ubuntu/", "noble", targets)
	if err != nil {
		t.Fatalf("applyMirrorToAPTTargets returned error: %v", err)
	}

	if len(result.UpdatedFiles) != 1 || len(result.BackupFiles) != 1 {
		t.Fatalf(
			"applyMirrorToAPTTargets updated=%d backups=%d, want 1 and 1",
			len(result.UpdatedFiles),
			len(result.BackupFiles),
		)
	}

	if result.UpdatedFiles[0] != deb822Path {
		t.Fatalf("applyMirrorToAPTTargets updated %q, want %q", result.UpdatedFiles[0], deb822Path)
	}

	for _, backupPath := range result.BackupFiles {
		if _, err := os.Stat(backupPath); err != nil {
			t.Fatalf("expected backup file %s to exist: %v", backupPath, err)
		}
		if filepath.Base(filepath.Dir(backupPath)) != mirrorselectBackupDirName {
			t.Fatalf("expected backup file %s to be stored in %s", backupPath, mirrorselectBackupDirName)
		}
	}

	deb822Content, err := os.ReadFile(deb822Path)
	if err != nil {
		t.Fatalf("failed to read updated deb822 file: %v", err)
	}
	if !strings.Contains(string(deb822Content), "noble-security") {
		t.Fatalf("updated deb822 file missing security suite:\n%s", string(deb822Content))
	}

	legacyContent, err := os.ReadFile(legacyPath)
	if err != nil {
		t.Fatalf("failed to read updated legacy file: %v", err)
	}
	if string(legacyContent) != "legacy old\n" {
		t.Fatalf("legacy file should remain unchanged when deb822 is present:\n%s", string(legacyContent))
	}
}

func TestApplyMirrorToAPTTargetsFallsBackToLegacy(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	legacyPath := filepath.Join(dir, "sources.list")

	if err := os.WriteFile(legacyPath, []byte("legacy old\n"), 0644); err != nil {
		t.Fatalf("failed to write legacy fixture: %v", err)
	}

	targets := []aptSourceTarget{
		{Path: filepath.Join(dir, "ubuntu.sources"), Format: aptSourceFormatDeb822},
		{Path: legacyPath, Format: aptSourceFormatLegacy},
	}

	result, err := applyMirrorToAPTTargets("https://mirror.example/ubuntu/", "noble", targets)
	if err != nil {
		t.Fatalf("applyMirrorToAPTTargets returned error: %v", err)
	}

	if len(result.UpdatedFiles) != 1 || result.UpdatedFiles[0] != legacyPath {
		t.Fatalf("applyMirrorToAPTTargets updated %v, want [%s]", result.UpdatedFiles, legacyPath)
	}

	legacyContent, err := os.ReadFile(legacyPath)
	if err != nil {
		t.Fatalf("failed to read updated legacy file: %v", err)
	}
	if !strings.Contains(string(legacyContent), "deb-src https://mirror.example/ubuntu noble-security") {
		t.Fatalf("updated legacy file missing deb-src security entry:\n%s", string(legacyContent))
	}
}
