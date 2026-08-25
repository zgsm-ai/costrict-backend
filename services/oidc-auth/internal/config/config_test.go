package config

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/spf13/viper"
)

// TestLoadConfigWithUserCenter verifies the userCenter section unmarshals from
// YAML, including the "5s"-style string → time.Duration decode hook.
func TestLoadConfigWithUserCenter(t *testing.T) {
	yaml := `
serverPort: "8080"
providers:
  casdoor:
    clientID: id
    clientSecret: secret
    baseURL: http://casdoor:30080
    internalURL: http://casdoor:8000
userCenter:
  baseURL: http://cs-user:8082
  internalToken: shared-secret
  timeout: "5s"
`
	dir := t.TempDir()
	path := filepath.Join(dir, "config.yaml")
	if err := os.WriteFile(path, []byte(yaml), 0o644); err != nil {
		t.Fatal(err)
	}

	viper.Reset()
	viper.SetConfigFile(path)
	if err := viper.ReadInConfig(); err != nil {
		t.Fatal(err)
	}
	cfg := new(AppConfig)
	if err := viper.Unmarshal(cfg); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}

	if cfg.UserCenter.BaseURL != "http://cs-user:8082" {
		t.Errorf("BaseURL = %q, want http://cs-user:8082", cfg.UserCenter.BaseURL)
	}
	if cfg.UserCenter.InternalToken != "shared-secret" {
		t.Errorf("InternalToken = %q, want shared-secret", cfg.UserCenter.InternalToken)
	}
	if cfg.UserCenter.Timeout != 5*time.Second {
		t.Errorf("Timeout = %v, want 5s", cfg.UserCenter.Timeout)
	}
}
