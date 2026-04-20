package app

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"

	"github.com/haukened/mirrorselect/internal/llog"
)

// contains checks if a given string is present in a slice of strings.
// It returns true if the string is found, and false otherwise.
//
// Parameters:
//
//	slice []string - the slice of strings to search within
//	s string - the string to search for
//
// Returns:
//
//	bool - true if the string is found in the slice, false otherwise
func contains(slice []string, s string) bool {
	for _, item := range slice {
		if item == s {
			return true
		}
	}
	return false
}

// getDistribCodename reads the /etc/lab-release file to find and return the
// distribution codename. It looks for a line that starts with "DISTRIB_CODENAME="
// and returns the value after the equals sign. If the file cannot be opened or
// read, or if the codename is not found, it returns an error.
func getDistribCodename() (string, error) {
	osInfo, err := getOSReleaseInfo()
	if err == nil {
		if osInfo.VersionCodename != "" {
			return osInfo.VersionCodename, nil
		}
		if osInfo.UbuntuCodename != "" {
			return osInfo.UbuntuCodename, nil
		}
	} else if !errors.Is(err, os.ErrNotExist) {
		return "", err
	}

	if codename, err := getLSBReleaseCodename(); err == nil && codename != "" {
		return codename, nil
	} else if err != nil && !errors.Is(err, os.ErrNotExist) {
		return "", err
	}

	output, err := exec.Command("lsb_release", "-cs").Output()
	if err == nil {
		codename := strings.TrimSpace(string(output))
		if codename != "" {
			return codename, nil
		}
	}

	return "", nil
}

func ensureUbuntuHost() error {
	if runtime.GOOS != "linux" {
		return fmt.Errorf("unsupported operating system %q (Ubuntu Linux required)", runtime.GOOS)
	}

	osInfo, err := getOSReleaseInfo()
	if err != nil {
		return fmt.Errorf("failed to detect operating system: %w", err)
	}

	if isUbuntuDistribution(osInfo) {
		return nil
	}

	if osInfo.ID == "" {
		return errors.New("unsupported distribution (Ubuntu required)")
	}

	return fmt.Errorf("unsupported distribution %q (Ubuntu required)", osInfo.ID)
}

type osReleaseInfo struct {
	ID              string
	IDLike          []string
	VersionCodename string
	UbuntuCodename  string
}

func getOSReleaseInfo() (osReleaseInfo, error) {
	for _, path := range []string{"/etc/os-release", "/usr/lib/os-release"} {
		data, err := parseKeyValueFile(path)
		if err != nil {
			if errors.Is(err, os.ErrNotExist) {
				continue
			}
			return osReleaseInfo{}, err
		}

		idLike := strings.Fields(strings.ToLower(data["ID_LIKE"]))

		return osReleaseInfo{
			ID:              strings.ToLower(data["ID"]),
			IDLike:          idLike,
			VersionCodename: strings.TrimSpace(data["VERSION_CODENAME"]),
			UbuntuCodename:  strings.TrimSpace(data["UBUNTU_CODENAME"]),
		}, nil
	}

	return osReleaseInfo{}, os.ErrNotExist
}

func isUbuntuDistribution(osInfo osReleaseInfo) bool {
	if osInfo.ID == "ubuntu" {
		return true
	}

	for _, id := range osInfo.IDLike {
		if id == "ubuntu" {
			return true
		}
	}

	return false
}

func getLSBReleaseCodename() (string, error) {
	file, err := os.Open("/etc/lsb-release")
	if err != nil {
		return "", err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "DISTRIB_CODENAME=") {
			return strings.TrimPrefix(line, "DISTRIB_CODENAME="), nil
		}
	}

	if err := scanner.Err(); err != nil {
		return "", err
	}

	return "", nil
}

func parseKeyValueFile(path string) (map[string]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	entries := make(map[string]string)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}

		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])
		value = strings.Trim(value, "\"'")
		entries[key] = value
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return entries, nil
}

type geoIP struct {
	IP          string `json:"ip"`
	ASO         string `json:"aso"`
	ASN         string `json:"asn"`
	Continent   string `json:"continent"`
	CountryCode string `json:"cc"`
	CountryName string `json:"country"`
	City        string `json:"city"`
	PostalCode  string `json:"postal"`
	Latitude    string `json:"latitude"`
	Longitude   string `json:"longitude"`
	Timezone    string `json:"tz"`
}

func (g *geoIP) UnmarshalJSON(data []byte) error {
	type geoIPAlias struct {
		IP          json.RawMessage `json:"ip"`
		ASO         json.RawMessage `json:"aso"`
		ASN         json.RawMessage `json:"asn"`
		Continent   json.RawMessage `json:"continent"`
		CountryCode json.RawMessage `json:"cc"`
		CountryName json.RawMessage `json:"country"`
		City        json.RawMessage `json:"city"`
		PostalCode  json.RawMessage `json:"postal"`
		Latitude    json.RawMessage `json:"latitude"`
		Longitude   json.RawMessage `json:"longitude"`
		Timezone    json.RawMessage `json:"tz"`
	}

	var alias geoIPAlias
	if err := json.Unmarshal(data, &alias); err != nil {
		return err
	}

	ip, err := normalizeStringField(alias.IP)
	if err != nil {
		return err
	}

	aso, err := normalizeStringField(alias.ASO)
	if err != nil {
		return err
	}

	asn, err := normalizeStringField(alias.ASN)
	if err != nil {
		return err
	}

	continent, err := normalizeStringField(alias.Continent)
	if err != nil {
		return err
	}

	countryCode, err := normalizeStringField(alias.CountryCode)
	if err != nil {
		return err
	}

	countryName, err := normalizeStringField(alias.CountryName)
	if err != nil {
		return err
	}

	city, err := normalizeStringField(alias.City)
	if err != nil {
		return err
	}

	postalCode, err := normalizeStringField(alias.PostalCode)
	if err != nil {
		return err
	}

	latitude, err := normalizeStringField(alias.Latitude)
	if err != nil {
		return err
	}

	longitude, err := normalizeStringField(alias.Longitude)
	if err != nil {
		return err
	}

	timezone, err := normalizeStringField(alias.Timezone)
	if err != nil {
		return err
	}

	*g = geoIP{
		IP:          ip,
		ASO:         aso,
		ASN:         asn,
		Continent:   continent,
		CountryCode: countryCode,
		CountryName: countryName,
		City:        city,
		PostalCode:  postalCode,
		Latitude:    latitude,
		Longitude:   longitude,
		Timezone:    timezone,
	}

	return nil
}

func normalizeStringField(raw json.RawMessage) (string, error) {
	trimmed := strings.TrimSpace(string(raw))
	if trimmed == "" || trimmed == "null" {
		return "", nil
	}

	var asString string
	if err := json.Unmarshal(raw, &asString); err == nil {
		return asString, nil
	}

	var asNumber json.Number
	if err := json.Unmarshal(raw, &asNumber); err == nil {
		return asNumber.String(), nil
	}

	var asBool bool
	if err := json.Unmarshal(raw, &asBool); err == nil {
		return strconv.FormatBool(asBool), nil
	}

	return "", fmt.Errorf("invalid scalar value: %s", trimmed)
}

// getGeoIP fetches the geoIP information from https://ident.me/json and parses it into a geoIP struct
func getGeoIP() (*geoIP, error) {
	llog.Debug("Fetching geoIP data")
	resp, err := http.Get("https://ident.me/json")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	llog.Debugf("Response status: %s", resp.Status)
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch geoIP data: %s", resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var geo geoIP
	if err := json.Unmarshal(body, &geo); err != nil {
		return nil, err
	}
	llog.Debugf("GeoIP data: %+v", geo)
	return &geo, nil
}

func humanizeTransferSpeed(bytes int64, seconds float64) string {
	bits := bytes * 8
	if seconds == 0 {
		return "0 b/s"
	}
	speed := float64(bits) / seconds
	units := []string{"b/s", "Kbps", "Mbps", "Gbps", "Tbps"}
	for _, unit := range units {
		if speed < 1024 {
			return fmt.Sprintf("%4.2f %s", speed, unit)
		}
		speed /= 1024
	}
	return fmt.Sprintf("%.2f %s", speed, "Tbps")
}
