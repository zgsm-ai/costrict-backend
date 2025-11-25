package chclient

import (
	"encoding/json"
	"fmt"
	"os"
)

// AuthFileConfig represents the authentication configuration loaded from authfile
type AuthFileConfig struct {
	Fingerprint *string            `json:"fingerprint"`
	Auth        *string            `json:"auth"`
	TLS         *AuthTLSConfig     `json:"tls"`
	Headers     *map[string]string `json:"headers"`
	Hostname    *string            `json:"hostname"`
	SNI         *string            `json:"sni"`
	ClientID    *string            `json:"clientId"`
	AppName     *string            `json:"appName"`
	UserID      *string            `json:"userId"`
}

// AuthTLSConfig represents TLS-related authentication configuration
type AuthTLSConfig struct {
	CA         *string `json:"ca"`
	SkipVerify *bool   `json:"skip_verify"`
	Cert       *string `json:"cert"`
	Key        *string `json:"key"`
}

// LoadAuthFile loads and parses the authentication configuration from a file
func LoadAuthFile(filename string) (*AuthFileConfig, error) {
	if filename == "" {
		return nil, nil
	}

	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to read auth file %s: %v", filename, err)
	}

	var config AuthFileConfig
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse auth file %s: %v", filename, err)
	}

	return &config, nil
}
