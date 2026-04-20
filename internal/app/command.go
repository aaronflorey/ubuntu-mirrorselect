package app

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"runtime"
	"sort"
	"strconv"
	"strings"
	"sync"

	"github.com/biter777/countries"
	"github.com/haukened/mirrorselect/internal/llog"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	version      = "dev"
	FinalMirrors []Mirror
)

type sudoReexecResult struct {
	exitCode int
}

func (r *sudoReexecResult) Error() string {
	return "sudo re-exec completed"
}

func SudoReexecExitCode(err error) (int, bool) {
	var result *sudoReexecResult
	if errors.As(err, &result) {
		return result.exitCode, true
	}

	return 0, false
}

type cliConfig struct {
	Arch        string
	Apply       bool
	AssumeYes   bool
	Country     string
	Interactive bool
	Max         int
	Output      string
	Protocol    string
	Release     string
	Timeout     int
	Verbosity   string
}

func NewRootCmd() *cobra.Command {
	v := viper.New()
	v.SetEnvPrefix("MIRRORSELECT")
	v.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
	v.AutomaticEnv()

	cfg := cliConfig{}

	cmd := &cobra.Command{
		Use:     "mirrorselect",
		Short:   "Select the fastest Ubuntu mirrors",
		Long:    "MirrorSelect discovers Ubuntu archive mirrors and ranks them by latency and download speed.",
		Version: version,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			if err := readConfig(v, &cfg); err != nil {
				return err
			}
			return before(&cfg)
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := run(&cfg); err != nil {
				return err
			}
			return after(&cfg)
		},
		SilenceUsage:  true,
		SilenceErrors: true,
	}

	flags := cmd.Flags()
	flags.StringP("arch", "a", runtime.GOARCH, "Architecture to select mirrors for (amd64, i386, arm64, armhf, ppc64el, riscv64, s390x)")
	flags.Bool("apply", false, "Apply the selected mirror to APT source files")
	flags.Bool("yes", false, "Skip confirmation prompt when using --apply")
	flags.StringP("country", "c", "", "Country to select mirrors from (ISO 3166-1 alpha-2 country code)")
	flags.BoolP("interactive", "i", false, "Interactively select a mirror from ranked results")
	flags.IntP("max", "m", 5, "Maximum number of mirrors to test (if available)")
	flags.StringP("output", "o", "text", "Output format (text, json)")
	flags.StringP("protocol", "p", "any", "Protocol to select mirrors for (http, https, any)")
	flags.StringP("release", "r", "", "Release to select mirrors for")
	flags.IntP("timeout", "t", 500, "Timeout for testing mirrors in milliseconds")
	flags.StringP("verbosity", "v", "WARN", "Set the log verbosity level (DEBUG, INFO, WARN, ERROR)")

	if err := v.BindPFlags(flags); err != nil {
		panic(err)
	}

	return cmd
}

func readConfig(v *viper.Viper, cfg *cliConfig) error {
	country := normalizeCountryCode(v.GetString("country"))

	cfg.Arch = strings.ToLower(strings.TrimSpace(v.GetString("arch")))
	cfg.Apply = v.GetBool("apply")
	cfg.AssumeYes = v.GetBool("yes")
	cfg.Country = country
	cfg.Interactive = v.GetBool("interactive")
	cfg.Max = v.GetInt("max")
	cfg.Output = strings.ToLower(strings.TrimSpace(v.GetString("output")))
	cfg.Protocol = strings.ToLower(strings.TrimSpace(v.GetString("protocol")))
	cfg.Release = strings.TrimSpace(v.GetString("release"))
	cfg.Timeout = v.GetInt("timeout")
	cfg.Verbosity = strings.ToUpper(strings.TrimSpace(v.GetString("verbosity")))

	return validateConfig(cfg)
}

func validateConfig(cfg *cliConfig) error {
	allowedArchs := []string{"amd64", "i386", "arm64", "armhf", "ppc64el", "riscv64", "s390x"}
	if !contains(allowedArchs, cfg.Arch) {
		return fmt.Errorf("invalid architecture: %s", cfg.Arch)
	}

	if cfg.Country != "" {
		if len(cfg.Country) != 2 {
			return fmt.Errorf("invalid country code: %s", cfg.Country)
		}
		if countries.ByName(cfg.Country) == countries.Unknown {
			return fmt.Errorf("unknown country code: %s", cfg.Country)
		}
	}

	if cfg.Max < 1 {
		return fmt.Errorf("invalid max value: %d (must be >= 1)", cfg.Max)
	}

	allowedOutputs := []string{"text", "json"}
	if !contains(allowedOutputs, cfg.Output) {
		return fmt.Errorf("invalid output format: %s", cfg.Output)
	}

	if cfg.Interactive && cfg.Output != "text" {
		return errors.New("interactive mode is only supported with --output text")
	}

	if cfg.Apply && cfg.Output != "text" {
		return errors.New("apply mode is only supported with --output text")
	}

	if cfg.Timeout < 1 {
		return fmt.Errorf("invalid timeout value: %d (must be >= 1)", cfg.Timeout)
	}

	allowedProtocols := []string{"http", "https", "any"}
	if !contains(allowedProtocols, cfg.Protocol) {
		return fmt.Errorf("invalid protocol: %s", cfg.Protocol)
	}

	allowedLevels := []string{"DEBUG", "INFO", "WARN", "ERROR"}
	if !contains(allowedLevels, cfg.Verbosity) {
		return fmt.Errorf("invalid log level: %s", cfg.Verbosity)
	}

	return nil
}

func normalizeCountryCode(code string) string {
	normalized := strings.ToUpper(strings.TrimSpace(code))
	if normalized == "UK" {
		return "GB"
	}
	return normalized
}

func after(cfg *cliConfig) error {
	if cfg.Output == "json" {
		results := make([]map[string]any, 0, len(FinalMirrors))
		for i, mirror := range FinalMirrors {
			results = append(results, map[string]any{
				"rank":        i + 1,
				"url":         mirror.URL.String(),
				"latency_ms":  mirror.Latency,
				"speed_bps":   mirror.bps(),
				"speed_human": humanizeTransferSpeed(mirror.Size, mirror.Time),
			})
		}

		encoder := json.NewEncoder(os.Stdout)
		encoder.SetIndent("", "  ")
		return encoder.Encode(results)
	}

	if len(FinalMirrors) == 0 {
		if cfg.Apply {
			return errors.New("no mirrors available to apply")
		}
		llog.Info("Flag options resulted in no mirrors being selected")
		return nil
	}

	var selected Mirror
	selectedChosen := false

	if cfg.Interactive {
		choice, err := selectMirrorInteractively(FinalMirrors, os.Stdin, os.Stdout)
		if err != nil {
			return err
		}
		selected = choice
		selectedChosen = true
	}

	if cfg.Apply {
		if !selectedChosen {
			selected = FinalMirrors[0]
			selectedChosen = true
			fmt.Fprintf(os.Stdout, "Selected top-ranked mirror: %s\n", selected.URL)
		}

		if !cfg.AssumeYes {
			confirmed, err := confirmApplySelection(os.Stdin, os.Stdout)
			if err != nil {
				return err
			}
			if !confirmed {
				return errors.New("mirror apply cancelled")
			}
		}

		result, err := applyMirrorToAPT(selected.URL.String(), cfg.Release)
		if err != nil {
			return err
		}

		fmt.Fprintf(os.Stdout, "Applied mirror %s to APT source files.\n", result.MirrorURI)
		for i := range len(result.UpdatedFiles) {
			fmt.Fprintf(
				os.Stdout,
				"- Updated %s (backup: %s)\n",
				result.UpdatedFiles[i],
				result.BackupFiles[i],
			)
		}

		return nil
	}

	for i, mirror := range FinalMirrors {
		fmt.Printf("%d. %s %s\n", i+1, humanizeTransferSpeed(mirror.Size, mirror.Time), mirror.URL)
	}

	return nil
}

func confirmApplySelection(in io.Reader, out io.Writer) (bool, error) {
	reader := bufio.NewReader(in)
	fmt.Fprint(out, "Apply this mirror to APT source files? [y/N]: ")

	line, err := reader.ReadString('\n')
	if err != nil && !errors.Is(err, io.EOF) {
		return false, err
	}

	answer := strings.ToLower(strings.TrimSpace(line))
	return answer == "y" || answer == "yes", nil
}

func selectMirrorInteractively(mirrors []Mirror, in io.Reader, out io.Writer) (Mirror, error) {
	if len(mirrors) == 0 {
		return Mirror{}, errors.New("no mirrors available to select")
	}

	fmt.Fprintln(out, "Available mirrors:")
	for i, mirror := range mirrors {
		fmt.Fprintf(
			out,
			"%d) %s | %d ms | %s\n",
			i+1,
			humanizeTransferSpeed(mirror.Size, mirror.Time),
			mirror.Latency,
			mirror.URL,
		)
	}

	reader := bufio.NewReader(in)
	for {
		fmt.Fprintf(out, "Select a mirror [1-%d] (or q to cancel): ", len(mirrors))

		line, err := reader.ReadString('\n')
		if err != nil && !errors.Is(err, io.EOF) {
			return Mirror{}, err
		}

		choice := strings.TrimSpace(line)
		switch strings.ToLower(choice) {
		case "q", "quit", "exit":
			return Mirror{}, errors.New("mirror selection cancelled")
		}

		index, convErr := strconv.Atoi(choice)
		if convErr != nil || index < 1 || index > len(mirrors) {
			if errors.Is(err, io.EOF) {
				return Mirror{}, errors.New("invalid mirror selection")
			}
			fmt.Fprintf(out, "Invalid selection. Enter a number from 1 to %d.\n", len(mirrors))
			continue
		}

		selected := mirrors[index-1]
		fmt.Fprintf(out, "Selected mirror: %s\n", selected.URL)
		return selected, nil
	}
}

func before(cfg *cliConfig) error {
	if err := maybeReexecWithSudo(cfg); err != nil {
		return err
	}

	// set the logging level
	err := llog.SetLogLevel(cfg.Verbosity)
	if err != nil {
		return err
	}

	if err := ensureUbuntuHost(); err != nil {
		return err
	}

	// get the distribution codename
	if cfg.Release == "" {
		codename, err := getDistribCodename()
		if err != nil {
			return err
		}
		if codename == "" {
			return errors.New("failed to detect distribution codename")
		}
		llog.Infof("Detected distribution codename %s", codename)
		cfg.Release = codename
	}
	// ensure we have a country code
	if cfg.Country == "" {
		// get the country code from the geoIP
		geo, err := getGeoIP()
		if err != nil {
			llog.Error("Unable to auto-detect country code, please specify one manually using --country")
			return err
		}
		llog.Infof("Using public IP address %s", geo.IP)
		llog.Infof("Detected country %s (%s)", geo.CountryName, geo.CountryCode)
		cfg.Country = normalizeCountryCode(geo.CountryCode)
	}
	return nil
}

func maybeReexecWithSudo(cfg *cliConfig) error {
	if !cfg.Apply || os.Geteuid() == 0 {
		return nil
	}

	executablePath, err := os.Executable()
	if err != nil {
		return fmt.Errorf("failed to resolve executable path for sudo re-exec: %w", err)
	}

	args := append([]string{executablePath}, os.Args[1:]...)
	cmd := exec.Command("sudo", args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err = cmd.Run()
	if err == nil {
		return &sudoReexecResult{exitCode: 0}
	}

	var exitErr *exec.ExitError
	if errors.As(err, &exitErr) {
		return &sudoReexecResult{exitCode: exitErr.ExitCode()}
	}

	return fmt.Errorf("failed to invoke sudo: %w", err)
}

func run(cfg *cliConfig) error {
	// get the mirrors
	mirrors, err := getMirrors(cfg.Country, cfg.Protocol)
	if err != nil {
		return err
	}

	// test the mirrors for latency
	fmt.Fprintf(os.Stderr, "Testing TCP latency for %d mirrors\n", len(mirrors))
	testLatencyInParallel(mirrors, cfg.Timeout, cfg.Release)

	// filter out the invalid mirrors
	mirrors = filterInvalidMirrors(mirrors)

	// get the top N mirrors
	mirrors = TopNByLatency(mirrors, cfg.Max)

	// then test the mirrors for download speed
	fmt.Fprintf(os.Stderr, "Testing download speed for %d mirrors\n", len(mirrors))
	for idx := range len(mirrors) {
		// grab the pointer to the mirror so it can self-update
		mirror := &mirrors[idx]
		// test the download speed of the mirror
		mirror.TestDownload(cfg.Release, cfg.Arch)
	}

	sort.Sort(ByTransferSpeed(mirrors))
	FinalMirrors = mirrors

	return nil
}

func testLatencyInParallel(mirrors []Mirror, timeout int, release string) {
	workers := latencyWorkerCount(len(mirrors))
	if workers == 0 {
		return
	}

	jobs := make(chan int)
	var wg sync.WaitGroup

	for range workers {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for idx := range jobs {
				mirrors[idx].TestLatency(timeout, release)
			}
		}()
	}

	for idx := range len(mirrors) {
		jobs <- idx
	}
	close(jobs)
	wg.Wait()
}

func latencyWorkerCount(total int) int {
	if total <= 0 {
		return 0
	}

	workers := runtime.NumCPU() * 4
	if workers < 1 {
		workers = 1
	}
	if workers > 32 {
		workers = 32
	}
	if workers > total {
		workers = total
	}

	return workers
}
