package app

import (
	"errors"
	"fmt"
	"os"
	"runtime"
	"sort"
	"strings"

	"github.com/biter777/countries"
	"github.com/haukened/mirrorselect/internal/llog"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var FinalMirrors []Mirror

type cliConfig struct {
	Arch      string
	Country   string
	Max       int
	Protocol  string
	Release   string
	Timeout   int
	Verbosity string
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
		Version: "dev",
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
			return after()
		},
		SilenceUsage:  true,
		SilenceErrors: true,
	}

	flags := cmd.Flags()
	flags.StringP("arch", "a", runtime.GOARCH, "Architecture to select mirrors for (amd64, i386, arm64, armhf, ppc64el, riscv64, s390x)")
	flags.StringP("country", "c", "", "Country to select mirrors from (ISO 3166-1 alpha-2 country code)")
	flags.IntP("max", "m", 5, "Maximum number of mirrors to test (if available)")
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
	cfg.Country = country
	cfg.Max = v.GetInt("max")
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

func after() error {
	if len(FinalMirrors) == 0 {
		llog.Info("Flag options resulted in no mirrors being selected")
	} else {
		for i, mirror := range FinalMirrors {
			fmt.Printf("%d. %s %s\n", i+1, humanizeTransferSpeed(mirror.Size, mirror.Time), mirror.URL)
		}
	}
	return nil
}

func before(cfg *cliConfig) error {
	// set the logging level
	err := llog.SetLogLevel(cfg.Verbosity)
	if err != nil {
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

func run(cfg *cliConfig) error {
	// get the mirrors
	mirrors, err := getMirrors(cfg.Country, cfg.Protocol)
	if err != nil {
		return err
	}

	// test the mirrors for latency
	fmt.Fprintf(os.Stderr, "Testing TCP latency for %d mirrors\n", len(mirrors))
	timeout := cfg.Timeout
	for idx := range len(mirrors) {
		// grab the pointer to the mirror so it can self-update
		mirror := &mirrors[idx]
		// test the latency and validity of the mirror
		mirror.TestLatency(timeout, cfg.Release)
	}

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
		mirror.TestDownload(cfg.Release)
	}

	sort.Sort(ByTransferSpeed(mirrors))
	FinalMirrors = mirrors

	return nil
}
