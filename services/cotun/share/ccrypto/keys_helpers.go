package ccrypto

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"strings"
)

const CotunKeyPrefix = "ck-"

//  Relations between entities:
//
//   .............> PEM <...........
//   .               ^             .
//   .               |             .
//   .               |             .
// Seed -------> PrivateKey        .
//   .               ^             .
//   .               |             .
//   .               V             .
//   ..........> CotunKey ..........

func Seed2PEM(seed string) ([]byte, error) {
	privateKey, err := seed2PrivateKey(seed)
	if err != nil {
		return nil, err
	}

	return privateKey2PEM(privateKey)
}

func seed2CotunKey(seed string) ([]byte, error) {
	privateKey, err := seed2PrivateKey(seed)
	if err != nil {
		return nil, err
	}

	return privateKey2CotunKey(privateKey)
}

func seed2PrivateKey(seed string) (*ecdsa.PrivateKey, error) {
	if seed == "" {
		return ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	} else {
		return GenerateKeyGo119(elliptic.P256(), NewDetermRand([]byte(seed)))
	}
}

func privateKey2CotunKey(privateKey *ecdsa.PrivateKey) ([]byte, error) {
	b, err := x509.MarshalECPrivateKey(privateKey)
	if err != nil {
		return nil, err
	}

	encodedPrivateKey := make([]byte, base64.RawStdEncoding.EncodedLen(len(b)))
	base64.RawStdEncoding.Encode(encodedPrivateKey, b)

	return append([]byte(CotunKeyPrefix), encodedPrivateKey...), nil
}

func privateKey2PEM(privateKey *ecdsa.PrivateKey) ([]byte, error) {
	b, err := x509.MarshalECPrivateKey(privateKey)
	if err != nil {
		return nil, err
	}

	return pem.EncodeToMemory(&pem.Block{Type: "EC PRIVATE KEY", Bytes: b}), nil
}

func cotunKey2PrivateKey(cotunKey []byte) (*ecdsa.PrivateKey, error) {
	rawCotunKey := cotunKey[len(CotunKeyPrefix):]

	decodedPrivateKey := make([]byte, base64.RawStdEncoding.DecodedLen(len(rawCotunKey)))
	_, err := base64.RawStdEncoding.Decode(decodedPrivateKey, rawCotunKey)
	if err != nil {
		return nil, err
	}

	return x509.ParseECPrivateKey(decodedPrivateKey)
}

func CotunKey2PEM(cotunKey []byte) ([]byte, error) {
	privateKey, err := cotunKey2PrivateKey(cotunKey)
	if err == nil {
		return privateKey2PEM(privateKey)
	}

	return nil, err
}

func IsCotunKey(cotunKey []byte) bool {
	return strings.HasPrefix(string(cotunKey), CotunKeyPrefix)
}
