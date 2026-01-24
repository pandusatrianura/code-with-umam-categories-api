package scalar

import (
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestEnsureFileURL(t *testing.T) {
	tempDir := t.TempDir()
	originalDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("getwd: %v", err)
	}
	if err := os.Chdir(tempDir); err != nil {
		t.Fatalf("chdir: %v", err)
	}
	t.Cleanup(func() {
		_ = os.Chdir(originalDir)
	})

	absPath := filepath.Join(tempDir, "abs.txt")

	cases := []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "file-abs",
			input: "file://" + absPath,
			want:  "file://" + absPath,
		},
		{
			name:  "abs",
			input: absPath,
			want:  "file://" + absPath,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := ensureFileURL(tc.input)
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
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		_, _ = w.Write([]byte("ok"))
	}))
	t.Cleanup(server.Close)

	cases := []struct {
		name    string
		url     string
		want    string
		wantErr bool
	}{
		{
			name: "ok",
			url:  server.URL,
			want: "ok",
		},
		{
			name:    "bad-url",
			url:     "http://%",
			wantErr: true,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := fetchContentFromURL(tc.url)
			if tc.wantErr {
				if err == nil {
					t.Fatalf("expected error")
				}
				return
			}
			if err != nil {
				t.Fatalf("fetchContentFromURL error: %v", err)
			}
			if got != tc.want {
				t.Fatalf("fetchContentFromURL = %q, want %q", got, tc.want)
			}
		})
	}
}

func TestReadFileFromURL(t *testing.T) {
	tempDir := t.TempDir()
	filePath := filepath.Join(tempDir, "file.txt")
	content := []byte("data")
	if err := os.WriteFile(filePath, content, 0644); err != nil {
		t.Fatalf("write file: %v", err)
	}

	cases := []struct {
		name    string
		url     string
		want    []byte
		wantErr bool
		match   string
	}{
		{
			name: "ok",
			url:  "file://" + filePath,
			want: content,
		},
		{
			name:    "scheme",
			url:     "http://example.com/file.txt",
			wantErr: true,
			match:   "unsupported URL scheme",
		},
		{
			name:    "parse",
			url:     "file://%",
			wantErr: true,
			match:   "error parsing URL",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := readFileFromURL(tc.url)
			if tc.wantErr {
				if err == nil {
					t.Fatalf("expected error")
				}
				if tc.match != "" && !strings.Contains(err.Error(), tc.match) {
					t.Fatalf("error = %q, want %q", err.Error(), tc.match)
				}
				return
			}
			if err != nil {
				t.Fatalf("readFileFromURL error: %v", err)
			}
			if string(got) != string(tc.want) {
				t.Fatalf("readFileFromURL = %q, want %q", string(got), string(tc.want))
			}
		})
	}
}
