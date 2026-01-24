package scalar

import (
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
)

func TestEnsureFileURL(t *testing.T) {
	tmpDir := t.TempDir()
	origWD, err := os.Getwd()
	if err != nil {
		t.Fatalf("getwd: %v", err)
	}
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("chdir: %v", err)
	}
	t.Cleanup(func() {
		_ = os.Chdir(origWD)
	})

	absPath := filepath.Join(tmpDir, "a.txt")
	tests := []struct {
		name string
		in   string
		want string
	}{
		{name: "file abs", in: "file://" + absPath, want: "file://" + absPath},
		{name: "file rel", in: "file://b.txt", want: "file://" + filepath.Join(tmpDir, "b.txt")},
		{name: "abs", in: absPath, want: "file://" + absPath},
		{name: "rel", in: "c.txt", want: "file://" + filepath.Join(tmpDir, "c.txt")},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, err := ensureFileURL(tc.in)
			if err != nil {
				t.Fatalf("ensureFileURL error: %v", err)
			}
			if got != tc.want {
				t.Fatalf("ensureFileURL = %q, want %q", got, tc.want)
			}
		})
	}
}

func TestFetchContentFromURL(t *testing.T) {
	tests := []struct {
		name string
		body string
	}{
		{name: "empty", body: ""},
		{name: "text", body: "hello"},
		{name: "multi", body: "line1\nline2"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				_, _ = w.Write([]byte(tc.body))
			}))
			defer srv.Close()

			got, err := fetchContentFromURL(srv.URL)
			if err != nil {
				t.Fatalf("fetchContentFromURL error: %v", err)
			}
			if got != tc.body {
				t.Fatalf("fetchContentFromURL = %q, want %q", got, tc.body)
			}
		})
	}
}

func TestReadFileFromURL(t *testing.T) {
	tmpDir := t.TempDir()
	path := filepath.Join(tmpDir, "data.txt")
	want := []byte("content")
	if err := os.WriteFile(path, want, 0o600); err != nil {
		t.Fatalf("writefile: %v", err)
	}

	tests := []struct {
		name    string
		in      string
		want    []byte
		wantErr bool
	}{
		{name: "file", in: "file://" + path, want: want},
		{name: "scheme", in: "http://example.com", wantErr: true},
		{name: "bad", in: "http://[::1", wantErr: true},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, err := readFileFromURL(tc.in)
			if tc.wantErr {
				if err == nil {
					t.Fatalf("readFileFromURL error: expected")
				}
				return
			}
			if err != nil {
				t.Fatalf("readFileFromURL error: %v", err)
			}
			if string(got) != string(tc.want) {
				t.Fatalf("readFileFromURL = %q, want %q", got, tc.want)
			}
		})
	}
}
