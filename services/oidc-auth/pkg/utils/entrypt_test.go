package utils

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/zgsm-ai/oidc-auth/internal/config"
)

// freshKeyManager builds an EncryptKeyManager with the given encrypt.privateKey
// config value and runs loadKeys directly, bypassing the package-level
// sync.Once singleton so each resolution branch can be exercised independently.
func freshKeyManager(t *testing.T, value string) (*EncryptKeyManager, error) {
	t.Helper()
	m := &EncryptKeyManager{
		Config: &config.EncryptConfig{
			PrivateKey: value,
			AesKey:     "0123456789abcdef0123456789abcdef",
		},
	}
	return m, m.loadKeys()
}

// testRSAKeyPEM returns a throwaway 2048-bit RSA private key as PKCS#1 PEM.
func testRSAKeyPEM(t *testing.T) []byte {
	t.Helper()
	k, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatalf("GenerateKey: %v", err)
	}
	return pem.EncodeToMemory(&pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(k)})
}

func TestLoadKeys_FromPath(t *testing.T) {
	path := filepath.Join(t.TempDir(), "private.pem")
	if err := os.WriteFile(path, testRSAKeyPEM(t), 0o600); err != nil {
		t.Fatalf("write private key: %v", err)
	}
	m, err := freshKeyManager(t, path)
	if err != nil {
		t.Fatalf("loadKeys: %v", err)
	}
	if m.privateKeyPEM == "" || m.publicKeyPEM == "" {
		t.Fatal("keys not loaded from path")
	}
}

func TestLoadKeys_InlinePEM(t *testing.T) {
	pemContent := string(testRSAKeyPEM(t))
	m, err := freshKeyManager(t, pemContent)
	if err != nil {
		t.Fatalf("loadKeys: %v", err)
	}
	if m.privateKeyPEM != strings.TrimSpace(pemContent) {
		t.Error("privateKeyPEM does not match inline PEM content")
	}
	if m.publicKeyPEM == "" {
		t.Fatal("publicKeyPEM empty for inline PEM")
	}
}

func TestLoadKeys_Base64PEM(t *testing.T) {
	b64 := base64.StdEncoding.EncodeToString(testRSAKeyPEM(t))
	m, err := freshKeyManager(t, b64)
	if err != nil {
		t.Fatalf("loadKeys: %v", err)
	}
	if m.privateKeyPEM == "" || m.publicKeyPEM == "" {
		t.Fatal("keys not loaded from base64 PEM")
	}
}

func TestLoadKeys_InvalidValue(t *testing.T) {
	// Not a path, not inline PEM, not valid base64 → must error.
	if _, err := freshKeyManager(t, "definitely-not-a-key"); err == nil {
		t.Fatal("want error for invalid value, got nil")
	}
}

func TestLoadKeys_EmptyValue(t *testing.T) {
	if _, err := freshKeyManager(t, ""); err == nil {
		t.Fatal("want error for empty value, got nil")
	}
}
