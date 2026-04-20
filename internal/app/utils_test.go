package app

import (
	"encoding/json"
	"testing"
)

func TestIsUbuntuDistribution(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		info osReleaseInfo
		want bool
	}{
		{name: "id ubuntu", info: osReleaseInfo{ID: "ubuntu"}, want: true},
		{name: "id-like ubuntu", info: osReleaseInfo{ID: "debian", IDLike: []string{"debian", "ubuntu"}}, want: true},
		{name: "non-ubuntu distro", info: osReleaseInfo{ID: "fedora", IDLike: []string{"rhel"}}, want: false},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			got := isUbuntuDistribution(tc.info)
			if got != tc.want {
				t.Fatalf("isUbuntuDistribution(%+v)=%t, want %t", tc.info, got, tc.want)
			}
		})
	}
}

func TestNormalizeStringField(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		rawJSON string
		want    string
		wantErr bool
	}{
		{name: "number", rawJSON: `13335`, want: "13335"},
		{name: "string", rawJSON: `"AS13335"`, want: "AS13335"},
		{name: "null", rawJSON: `null`, want: ""},
		{name: "bool", rawJSON: `true`, want: "true"},
		{name: "invalid object", rawJSON: `{}`, wantErr: true},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			got, err := normalizeStringField(json.RawMessage(tc.rawJSON))
			if tc.wantErr {
				if err == nil {
					t.Fatalf("normalizeStringField(%s) expected error, got nil", tc.rawJSON)
				}
				return
			}

			if err != nil {
				t.Fatalf("normalizeStringField(%s) returned error: %v", tc.rawJSON, err)
			}

			if got != tc.want {
				t.Fatalf("normalizeStringField(%s)=%q, want %q", tc.rawJSON, got, tc.want)
			}
		})
	}
}

func TestGeoIPUnmarshalASNNumber(t *testing.T) {
	t.Parallel()

	data := []byte(`{
		"ip":"203.0.113.1",
		"aso":"Example ASN",
		"asn":13335,
		"continent":"North America",
		"cc":"US",
		"country":"United States",
		"latitude":40.7128,
		"longitude":"-74.0060"
	}`)

	var geo geoIP
	if err := json.Unmarshal(data, &geo); err != nil {
		t.Fatalf("json.Unmarshal returned error: %v", err)
	}

	if geo.ASN != "13335" {
		t.Fatalf("geo.ASN=%q, want %q", geo.ASN, "13335")
	}
	if geo.CountryCode != "US" {
		t.Fatalf("geo.CountryCode=%q, want %q", geo.CountryCode, "US")
	}
	if geo.Latitude != "40.7128" {
		t.Fatalf("geo.Latitude=%q, want %q", geo.Latitude, "40.7128")
	}
}

func TestGeoIPUnmarshalRejectsObjectField(t *testing.T) {
	t.Parallel()

	data := []byte(`{"asn":{}}`)

	var geo geoIP
	if err := json.Unmarshal(data, &geo); err == nil {
		t.Fatal("json.Unmarshal expected error for object scalar field, got nil")
	}
}
