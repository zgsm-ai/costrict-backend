package utils

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
	"sync"

	"github.com/zgsm-ai/oidc-auth/internal/config"
)

var (
	once              sync.Once
	encryptKeyManager *EncryptKeyManager
	globalConfig      *config.AppConfig
)

type EncryptKeyManager struct {
	privateKeyPEM string
	publicKeyPEM  string
	aesKey        string
	mu            sync.RWMutex
	Config        *config.EncryptConfig
}

func GetEncryptKeyManager() (*EncryptKeyManager, error) {
	var initErr error
	once.Do(func() {
		if globalConfig == nil {
			initErr = fmt.Errorf("global config not initialized")
			return
		}
		encryptKeyManager = &EncryptKeyManager{
			Config: &globalConfig.Encrypt,
		}
		initErr = encryptKeyManager.loadKeys()
	})
	if initErr != nil {
		return nil, initErr
	}
	return encryptKeyManager, nil
}

func SetGlobalConfig(cfg *config.AppConfig) {
	globalConfig = cfg
}

func (m *EncryptKeyManager) loadKeys() error {
	m.mu.Lock()
	defer m.mu.Unlock()

	keyValue := strings.TrimSpace(m.Config.PrivateKey)
	if keyValue == "" {
		return errors.New("failed to load private key: empty encrypt.privateKey value")
	}

	// encrypt.privateKey accepts three forms: a path to a PEM file, the inline
	// PEM content itself, or base64(PEM). Inline PEM always carries the BEGIN
	// marker; a non-readable value falls back to base64 decoding. The fallback
	// must not gate on os.IsNotExist: a real base64(PEM) string is one huge
	// filename component, and Windows/Linux reject it with ERROR_INVALID_NAME /
	// ENAMETOOLONG rather than a not-exist error.
	var privateKeyBytes []byte
	if strings.HasPrefix(keyValue, "-----BEGIN") {
		privateKeyBytes = []byte(keyValue)
	} else {
		var err error
		privateKeyBytes, err = os.ReadFile(keyValue)
		if err != nil {
			decoded, derr := base64.StdEncoding.DecodeString(keyValue)
			if derr != nil {
				return fmt.Errorf("failed to read private key file %s: %v", keyValue, err)
			}
			privateKeyBytes = decoded
		}
	}
	m.privateKeyPEM = string(privateKeyBytes)

	publicKeyPEM, err := derivePublicKeyPEM(privateKeyBytes)
	if err != nil {
		return fmt.Errorf("failed to derive public key from private key: %v", err)
	}
	m.publicKeyPEM = publicKeyPEM
	m.aesKey = m.Config.AesKey

	return nil
}

// derivePublicKeyPEM extracts the public key embedded in an RSA private key PEM
// and re-encodes it as a PKIX "PUBLIC KEY" PEM block, so callers never need a
// separate public key file.
func derivePublicKeyPEM(privateKeyPEM []byte) (string, error) {
	block, _ := pem.Decode(privateKeyPEM)
	if block == nil {
		return "", errors.New("failed to decode private key PEM")
	}

	var privKey *rsa.PrivateKey
	if parsed, err := x509.ParsePKCS1PrivateKey(block.Bytes); err == nil {
		privKey = parsed
	} else if parsed, err := x509.ParsePKCS8PrivateKey(block.Bytes); err == nil {
		var ok bool
		privKey, ok = parsed.(*rsa.PrivateKey)
		if !ok {
			return "", errors.New("private key is not RSA")
		}
	} else {
		return "", errors.New("unsupported private key format")
	}

	pubDER, err := x509.MarshalPKIXPublicKey(&privKey.PublicKey)
	if err != nil {
		return "", fmt.Errorf("failed to marshal public key: %v", err)
	}
	return string(pem.EncodeToMemory(&pem.Block{Type: "PUBLIC KEY", Bytes: pubDER})), nil
}

func (m *EncryptKeyManager) GetPrivateKeyPEM() string {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.privateKeyPEM
}

func (m *EncryptKeyManager) GetPublicKeyPEM() string {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.publicKeyPEM
}

func (m *EncryptKeyManager) GetAesKey() (string, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.aesKey, nil
}

func (m *EncryptKeyManager) ReloadKeys() error {
	return m.loadKeys()
}

func (m *EncryptKeyManager) AESEncrypt(plaintext []byte) (string, error) {
	key := m.Config.AesKey
	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		return "", fmt.Errorf("failed to create cipher: %w", err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", fmt.Errorf("failed to create GCM: %w", err)
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", fmt.Errorf("failed to generate nonce: %w", err)
	}

	ciphertext := gcm.Seal(nonce, nonce, plaintext, nil)

	return base64.URLEncoding.EncodeToString(ciphertext), nil
}

func (m *EncryptKeyManager) AESDecrypt(encryptedText string) ([]byte, error) {
	key := m.Config.AesKey
	ciphertext, err := base64.URLEncoding.DecodeString(encryptedText)
	if err != nil {
		return nil, fmt.Errorf("failed to decode base64: %w", err)
	}

	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		return nil, fmt.Errorf("failed to create cipher: %w", err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("failed to create GCM: %w", err)
	}

	if len(ciphertext) < gcm.NonceSize() {
		return nil, errors.New("ciphertext too short")
	}

	nonce := ciphertext[:gcm.NonceSize()]
	ciphertext = ciphertext[gcm.NonceSize():]

	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt: %w", err)
	}

	return plaintext, nil
}
