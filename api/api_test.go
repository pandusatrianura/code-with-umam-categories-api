package api

import "testing"

func TestNewAPIServer(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		addr string
	}{
		{name: "empty", addr: ""},
		{name: "with-port", addr: "127.0.0.1:8080"},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			srv := NewAPIServer(tt.addr)
			if srv == nil {
				t.Fatalf("expected server, got nil")
			}
			if srv.addr != tt.addr {
				t.Fatalf("expected addr %q, got %q", tt.addr, srv.addr)
			}
		})
	}
}

func TestServerRun(t *testing.T) {
	tests := []struct {
		name string
		addr string
	}{
		{name: "missing-port", addr: "127.0.0.1"},
		{name: "invalid", addr: "://"},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			srv := NewAPIServer(tt.addr)
			if srv == nil {
				t.Fatalf("expected server, got nil")
			}

			err := srv.Run()
			if err == nil {
				t.Fatalf("expected error for addr %q, got nil", tt.addr)
			}
		})
	}
}
