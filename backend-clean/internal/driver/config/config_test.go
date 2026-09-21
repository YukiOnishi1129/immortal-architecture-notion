package config_test

import (
	"os"
	"testing"

	"immortal-architecture-cqrs/backend/internal/driver/config"
)

func TestLoad(t *testing.T) {
	tests := []struct {
		name       string
		envVars    map[string]string
		wantErr    bool
		wantPort   int
		wantOrigin int // number of allowed origins
	}{
		{
			name: "[Success] all env vars set",
			envVars: map[string]string{
				"DATABASE_URL":  "postgres://user:pass@localhost:5432/db",
				"API_PORT":      "9090",
				"CLIENT_ORIGIN": "http://example.com,http://test.com",
			},
			wantPort:   9090,
			wantOrigin: 2,
		},
		{
			name: "[Success] defaults when optional vars empty",
			envVars: map[string]string{
				"DATABASE_URL": "postgres://user:pass@localhost:5432/db",
			},
			wantPort:   8080,
			wantOrigin: 2, // default localhost origins
		},
		{
			name: "[Success] CLIENT_ORIGIN with spaces",
			envVars: map[string]string{
				"DATABASE_URL":  "postgres://user:pass@localhost:5432/db",
				"CLIENT_ORIGIN": " http://example.com , http://test.com ",
			},
			wantPort:   8080,
			wantOrigin: 2,
		},
		{
			name: "[Success] CLIENT_ORIGIN with empty values filtered",
			envVars: map[string]string{
				"DATABASE_URL":  "postgres://user:pass@localhost:5432/db",
				"CLIENT_ORIGIN": "http://example.com,,http://test.com",
			},
			wantPort:   8080,
			wantOrigin: 2,
		},
		{
			name: "[Fail] missing DATABASE_URL",
			envVars: map[string]string{
				"API_PORT": "8080",
			},
			wantErr: true,
		},
		{
			name: "[Fail] invalid API_PORT",
			envVars: map[string]string{
				"DATABASE_URL": "postgres://user:pass@localhost:5432/db",
				"API_PORT":     "not-a-number",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Clear environment
			os.Clearenv()

			// Set test environment variables
			for k, v := range tt.envVars {
				os.Setenv(k, v)
			}

			cfg, err := config.Load()

			if tt.wantErr {
				if err == nil {
					t.Fatal("expected error but got nil")
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if cfg.DatabaseURL != tt.envVars["DATABASE_URL"] {
				t.Errorf("DatabaseURL = %q, want %q", cfg.DatabaseURL, tt.envVars["DATABASE_URL"])
			}

			if cfg.ServerPort != tt.wantPort {
				t.Errorf("ServerPort = %d, want %d", cfg.ServerPort, tt.wantPort)
			}

			if len(cfg.AllowedOrigins) != tt.wantOrigin {
				t.Errorf("len(AllowedOrigins) = %d, want %d", len(cfg.AllowedOrigins), tt.wantOrigin)
			}
		})
	}
}

func TestLoad_DefaultOrigins(t *testing.T) {
	os.Clearenv()
	os.Setenv("DATABASE_URL", "postgres://user:pass@localhost:5432/db")

	cfg, err := config.Load()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expectedOrigins := []string{"http://localhost:3000", "http://127.0.0.1:3000"}
	if len(cfg.AllowedOrigins) != len(expectedOrigins) {
		t.Fatalf("expected %d origins, got %d", len(expectedOrigins), len(cfg.AllowedOrigins))
	}

	for i, origin := range expectedOrigins {
		if cfg.AllowedOrigins[i] != origin {
			t.Errorf("AllowedOrigins[%d] = %q, want %q", i, cfg.AllowedOrigins[i], origin)
		}
	}
}

func TestLoad_WhitespaceOnlyOrigin(t *testing.T) {
	os.Clearenv()
	os.Setenv("DATABASE_URL", "postgres://user:pass@localhost:5432/db")
	os.Setenv("CLIENT_ORIGIN", "   ")

	cfg, err := config.Load()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Should fall back to defaults when only whitespace
	expectedOrigins := []string{"http://localhost:3000", "http://127.0.0.1:3000"}
	if len(cfg.AllowedOrigins) != len(expectedOrigins) {
		t.Fatalf("expected %d origins, got %d", len(expectedOrigins), len(cfg.AllowedOrigins))
	}
}
