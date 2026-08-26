package utils

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"

	"github.com/zgsm-ai/oidc-auth/internal/config"
	"github.com/zgsm-ai/oidc-auth/internal/repository"
)

// initKeys generates an ephemeral RSA keypair, writes the private key to a temp
// PEM file, and initializes the global key manager so the self-signed token
// functions have a signing key. The public key is derived from the private key
// at load time. pkg/utils has no other tests, so the sync.Once in
// GetEncryptKeyManager is still un-consumed when the first test calls it.
func initKeys(t *testing.T) {
	t.Helper()
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatalf("GenerateKey: %v", err)
	}
	dir := t.TempDir()
	privPath := filepath.Join(dir, "private.pem")

	privPEM := pem.EncodeToMemory(&pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(privateKey)})
	if err := os.WriteFile(privPath, privPEM, 0o600); err != nil {
		t.Fatalf("write private key: %v", err)
	}

	SetGlobalConfig(&config.AppConfig{
		Encrypt: config.EncryptConfig{
			PrivateKey: privPath,
			AesKey:     "0123456789abcdef0123456789abcdef",
		},
	})
	if _, err := GetEncryptKeyManager(); err != nil {
		t.Fatalf("GetEncryptKeyManager: %v", err)
	}
}

func testUserAndDevice() (*repository.AuthUser, repository.Device) {
	user := &repository.AuthUser{ID: uuid.MustParse("11111111-1111-1111-1111-111111111111")}
	device := repository.Device{MachineCode: "mc-1", VSCodeVersion: "vs-1"}
	return user, device
}

func TestGenerateRefreshToken_RoundTrip(t *testing.T) {
	initKeys(t)
	user, device := testUserAndDevice()
	exp := time.Now().Add(48 * time.Hour)

	rt, err := GenerateRefreshToken(user, device, exp)
	if err != nil {
		t.Fatalf("GenerateRefreshToken: %v", err)
	}

	claims, err := VerifyRefreshToken(rt)
	if err != nil {
		t.Fatalf("VerifyRefreshToken: %v", err)
	}
	if claims.Subject != user.ID.String() {
		t.Errorf("Subject = %q, want %q", claims.Subject, user.ID.String())
	}
	if claims.MachineCode != "mc-1" {
		t.Errorf("MachineCode = %q, want mc-1", claims.MachineCode)
	}
	if claims.VsCodeVersion != "vs-1" {
		t.Errorf("VsCodeVersion = %q, want vs-1", claims.VsCodeVersion)
	}
	if claims.TokenType != TokenTypeRefresh {
		t.Errorf("TokenType = %q, want %q", claims.TokenType, TokenTypeRefresh)
	}
	if claims.ExpiresAt == nil {
		t.Fatal("ExpiresAt = nil, want set")
	}
	if diff := claims.ExpiresAt.Time.Sub(exp); diff > time.Second || diff < -time.Second {
		t.Errorf("ExpiresAt = %v, want %v (diff %v)", claims.ExpiresAt.Time, exp, diff)
	}
}

func TestVerifyRefreshToken_Expired(t *testing.T) {
	initKeys(t)
	user, device := testUserAndDevice()

	rt, err := GenerateRefreshToken(user, device, time.Now().Add(-time.Hour))
	if err != nil {
		t.Fatalf("GenerateRefreshToken: %v", err)
	}
	if _, err := VerifyRefreshToken(rt); err == nil {
		t.Fatal("want error for expired token, got nil")
	}
}

func TestVerifyRefreshToken_WrongType(t *testing.T) {
	initKeys(t)
	keyManager, err := GetEncryptKeyManager()
	if err != nil {
		t.Fatalf("GetEncryptKeyManager: %v", err)
	}
	accessToken, err := CreateToken(AppClaims{
		TokenType: TokenTypeAccess,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
		},
	}, keyManager.GetPrivateKeyPEM())
	if err != nil {
		t.Fatalf("CreateToken: %v", err)
	}
	if _, err := VerifyRefreshToken(accessToken); err == nil {
		t.Fatal("want error for access token, got nil")
	}
}

func TestVerifyRefreshToken_Tampered(t *testing.T) {
	initKeys(t)
	user, device := testUserAndDevice()

	rt, err := GenerateRefreshToken(user, device, time.Now().Add(time.Hour))
	if err != nil {
		t.Fatalf("GenerateRefreshToken: %v", err)
	}
	parts := strings.Split(rt, ".")
	if len(parts) != 3 {
		t.Fatalf("token has %d parts, want 3", len(parts))
	}
	payload := parts[1]
	var tampered string
	if payload[len(payload)-1] == 'A' {
		tampered = payload[:len(payload)-1] + "B"
	} else {
		tampered = payload[:len(payload)-1] + "A"
	}
	bad := parts[0] + "." + tampered + "." + parts[2]
	if _, err := VerifyRefreshToken(bad); err == nil {
		t.Fatal("want error for tampered token, got nil")
	}
}

func TestVerifyRefreshToken_ForeignKey(t *testing.T) {
	initKeys(t)
	otherKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatalf("GenerateKey: %v", err)
	}
	otherPrivPEM := pem.EncodeToMemory(&pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(otherKey)})
	foreign, err := CreateToken(AppClaims{
		TokenType: TokenTypeRefresh,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
		},
	}, string(otherPrivPEM))
	if err != nil {
		t.Fatalf("CreateToken: %v", err)
	}
	if _, err := VerifyRefreshToken(foreign); err == nil {
		t.Fatal("want error for foreign-key token, got nil")
	}
}

func TestVerifyRefreshToken_Empty(t *testing.T) {
	initKeys(t)
	if _, err := VerifyRefreshToken(""); err == nil {
		t.Fatal("want error for empty token, got nil")
	}
}

func TestExtractTokenExp(t *testing.T) {
	initKeys(t)
	keyManager, err := GetEncryptKeyManager()
	if err != nil {
		t.Fatalf("GetEncryptKeyManager: %v", err)
	}
	exp := time.Now().Add(7 * 24 * time.Hour).Truncate(time.Second)
	tok, err := CreateToken(AppClaims{
		TokenType: TokenTypeRefresh,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(exp),
		},
	}, keyManager.GetPrivateKeyPEM())
	if err != nil {
		t.Fatalf("CreateToken: %v", err)
	}
	got, err := ExtractTokenExp(tok)
	if err != nil {
		t.Fatalf("ExtractTokenExp: %v", err)
	}
	if !got.Equal(exp) {
		t.Errorf("exp = %v, want %v", got, exp)
	}
}

func TestExtractTokenExp_NoExp(t *testing.T) {
	initKeys(t)
	keyManager, err := GetEncryptKeyManager()
	if err != nil {
		t.Fatalf("GetEncryptKeyManager: %v", err)
	}
	tok, err := CreateToken(jwt.MapClaims{"token_type": TokenTypeRefresh}, keyManager.GetPrivateKeyPEM())
	if err != nil {
		t.Fatalf("CreateToken: %v", err)
	}
	if _, err := ExtractTokenExp(tok); err == nil {
		t.Fatal("want error for token without exp, got nil")
	}
}

func TestVerifyRefreshToken_NotAJWT(t *testing.T) {
	initKeys(t)
	if _, err := VerifyRefreshToken("not-a-jwt"); err == nil {
		t.Fatal("want error for non-JWT, got nil")
	}
}
