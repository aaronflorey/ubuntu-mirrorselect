package app

import (
	"errors"
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const (
	aptSourceFormatDeb822 = "deb822"
	aptSourceFormatLegacy = "legacy"

	ubuntuArchiveKeyringPath = "/usr/share/keyrings/ubuntu-archive-keyring.gpg"
)

var defaultAptSourceTargets = []aptSourceTarget{
	{Path: "/etc/apt/sources.list.d/ubuntu.sources", Format: aptSourceFormatDeb822},
	{Path: "/etc/apt/sources.list", Format: aptSourceFormatLegacy},
}

type aptSourceTarget struct {
	Path   string
	Format string
}

type aptApplyResult struct {
	MirrorURI    string
	UpdatedFiles []string
	BackupFiles  []string
}

func applyMirrorToAPT(mirrorURL string, release string) (aptApplyResult, error) {
	if os.Geteuid() != 0 {
		return aptApplyResult{}, errors.New("--apply requires root privileges; run with sudo")
	}

	return applyMirrorToAPTTargets(mirrorURL, release, defaultAptSourceTargets)
}

func applyMirrorToAPTTargets(
	mirrorURL string,
	release string,
	targets []aptSourceTarget,
) (aptApplyResult, error) {
	release = strings.TrimSpace(release)
	if release == "" {
		return aptApplyResult{}, errors.New("release codename is required to apply APT sources")
	}

	mirrorURI, err := normalizeAptMirrorURI(mirrorURL)
	if err != nil {
		return aptApplyResult{}, err
	}

	result := aptApplyResult{MirrorURI: mirrorURI}

	for _, target := range targets {
		info, err := os.Stat(target.Path)
		if err != nil {
			if errors.Is(err, os.ErrNotExist) {
				continue
			}
			return aptApplyResult{}, fmt.Errorf("failed to inspect %s: %w", target.Path, err)
		}

		backupPath, err := createBackup(target.Path, info.Mode().Perm())
		if err != nil {
			return aptApplyResult{}, err
		}

		var content string
		switch target.Format {
		case aptSourceFormatDeb822:
			content = renderDeb822Sources(mirrorURI, release)
		case aptSourceFormatLegacy:
			content = renderLegacySourcesList(mirrorURI, release)
		default:
			return aptApplyResult{}, fmt.Errorf("unsupported APT source format %q", target.Format)
		}

		if err := writeFileAtomic(target.Path, []byte(content), info.Mode().Perm()); err != nil {
			return aptApplyResult{}, err
		}

		result.UpdatedFiles = append(result.UpdatedFiles, target.Path)
		result.BackupFiles = append(result.BackupFiles, backupPath)
	}

	if len(result.UpdatedFiles) == 0 {
		return aptApplyResult{}, errors.New("no APT source files found to update")
	}

	return result, nil
}

func normalizeAptMirrorURI(mirrorURL string) (string, error) {
	parsed, err := url.Parse(strings.TrimSpace(mirrorURL))
	if err != nil {
		return "", fmt.Errorf("invalid mirror URL: %w", err)
	}

	if parsed.Scheme != "http" && parsed.Scheme != "https" {
		return "", fmt.Errorf("unsupported mirror URL scheme %q", parsed.Scheme)
	}

	if parsed.Host == "" {
		return "", errors.New("mirror URL host cannot be empty")
	}

	trimmedPath := strings.TrimRight(parsed.Path, "/")
	if !strings.HasSuffix(trimmedPath, "/ubuntu") {
		if trimmedPath == "" {
			trimmedPath = "/ubuntu"
		} else {
			trimmedPath = trimmedPath + "/ubuntu"
		}
	}

	parsed.Path = trimmedPath
	parsed.RawQuery = ""
	parsed.Fragment = ""

	return parsed.String(), nil
}

func renderDeb822Sources(mirrorURI string, release string) string {
	release = strings.TrimSpace(release)
	return fmt.Sprintf(
		"# Managed by mirrorselect; original file backed up with .mirrorselect.bak timestamp\n"+
			"Types: deb deb-src\n"+
			"URIs: %s\n"+
			"Suites: %s %s-updates %s-backports %s-security\n"+
			"Components: main restricted universe multiverse\n"+
			"Signed-By: %s\n",
		mirrorURI,
		release,
		release,
		release,
		release,
		ubuntuArchiveKeyringPath,
	)
}

func renderLegacySourcesList(mirrorURI string, release string) string {
	release = strings.TrimSpace(release)
	suites := []string{release, release + "-updates", release + "-backports", release + "-security"}

	var builder strings.Builder
	builder.WriteString("# Managed by mirrorselect; original file backed up with .mirrorselect.bak timestamp\n")
	for _, suite := range suites {
		builder.WriteString(fmt.Sprintf("deb %s %s main restricted universe multiverse\n", mirrorURI, suite))
	}
	for _, suite := range suites {
		builder.WriteString(fmt.Sprintf("deb-src %s %s main restricted universe multiverse\n", mirrorURI, suite))
	}

	return builder.String()
}

func createBackup(path string, perm os.FileMode) (string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return "", fmt.Errorf("failed to read %s for backup: %w", path, err)
	}

	backupPath := fmt.Sprintf("%s.mirrorselect.bak.%s", path, time.Now().Format("20060102T150405"))
	if err := os.WriteFile(backupPath, data, perm); err != nil {
		return "", fmt.Errorf("failed to write backup %s: %w", backupPath, err)
	}

	return backupPath, nil
}

func writeFileAtomic(path string, data []byte, perm os.FileMode) error {
	dir := filepath.Dir(path)
	tmpFile, err := os.CreateTemp(dir, ".mirrorselect-tmp-")
	if err != nil {
		return fmt.Errorf("failed to create temp file for %s: %w", path, err)
	}

	tmpPath := tmpFile.Name()
	cleanup := true
	defer func() {
		if cleanup {
			_ = os.Remove(tmpPath)
		}
	}()

	if _, err := tmpFile.Write(data); err != nil {
		_ = tmpFile.Close()
		return fmt.Errorf("failed to write temp file for %s: %w", path, err)
	}

	if err := tmpFile.Chmod(perm); err != nil {
		_ = tmpFile.Close()
		return fmt.Errorf("failed to set file mode for %s: %w", path, err)
	}

	if err := tmpFile.Close(); err != nil {
		return fmt.Errorf("failed to close temp file for %s: %w", path, err)
	}

	if err := os.Rename(tmpPath, path); err != nil {
		return fmt.Errorf("failed to replace %s: %w", path, err)
	}

	cleanup = false
	return nil
}
