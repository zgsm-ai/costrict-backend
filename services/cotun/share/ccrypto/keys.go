package ccrypto

import (
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"os"

	"golang.org/x/crypto/ssh"
)

// GenerateKey generates a PEM key
func GenerateKey(seed string) ([]byte, error) {
	return Seed2PEM(seed)
}

// GenerateKeyFile generates an CotunKey
func GenerateKeyFile(keyFilePath, seed string) error {
	cotunKey, err := seed2CotunKey(seed)
	if err != nil {
		return err
	}

	if keyFilePath == "-" {
		fmt.Print(string(cotunKey))
		return nil
	}
	return os.WriteFile(keyFilePath, cotunKey, 0600)
}

// FingerprintKey calculates the SHA256 hash of an SSH public key
func FingerprintKey(k ssh.PublicKey) string {
	bytes := sha256.Sum256(k.Marshal())
	return base64.StdEncoding.EncodeToString(bytes[:])
}
