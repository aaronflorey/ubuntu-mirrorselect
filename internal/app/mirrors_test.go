package app

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestDownloadProbePath(t *testing.T) {
	t.Parallel()

	got := downloadProbePath("noble", "amd64")
	want := "dists/noble/main/binary-amd64/Packages.gz"
	if got != want {
		t.Fatalf("downloadProbePath()=%q, want %q", got, want)
	}
}

func TestMirrorTestDownloadUsesPackagesIndex(t *testing.T) {
	t.Parallel()

	const expectedPath = "/dists/noble/main/binary-amd64/Packages.gz"
	payload := []byte("package-index-payload")

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != expectedPath {
			http.NotFound(w, r)
			return
		}

		_, _ = fmt.Fprint(w, string(payload))
	}))
	defer server.Close()

	mirror, ok := NewMirror(server.URL + "/")
	if !ok {
		t.Fatalf("NewMirror(%q) returned false", server.URL+"/")
	}

	mirror.TestDownload("noble", "amd64")

	if !mirror.Valid {
		t.Fatal("TestDownload marked mirror invalid")
	}
	if mirror.Size != int64(len(payload)) {
		t.Fatalf("mirror.Size=%d, want %d", mirror.Size, len(payload))
	}
	if mirror.Time <= 0 {
		t.Fatalf("mirror.Time=%f, want > 0", mirror.Time)
	}
}
