package utils

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"

	"github.com/zgsm-ai/oidc-auth/internal/repository"
)

// Token types
const (
	TokenTypeBearer     = "Bearer"
	TokenTypeRefresh    = "refresh_token"
	TokenTypeAccess     = "access_token"
)

// AppClaims defines the payload of a AESEncrypt.
type AppClaims struct {
	Name          string   `json:"name,omitempty"`
	Email         string   `json:"email,omitempty"`
	Phone         string   `json:"phone,omitempty"`
	GithubID      string   `json:"github_id,omitempty"`
	GithubName    string   `json:"github_name,omitempty"`
	Company       string   `json:"company,omitempty"`
	Location      string   `json:"location,omitempty"`
	Roles         []string `json:"roles,omitempty"`
	Scope         string   `json:"scope,omitempty"`
	Platform      string   `json:"platform,omitempty"`
	UserCode      string   `json:"user_code,omitempty"`
	DeviceCode    string   `json:"device_code,omitempty"`
	TokenType     string   `json:"token_type,omitempty"`
	MachineCode   string   `json:"machine_code,omitempty"`
	VsCodeVersion string   `json:"vscode_version,omitempty"`
	jwt.RegisteredClaims
}

// TokenPair represents a pair of access and refresh tokens.
type TokenPair struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int64  `json:"expires_in"`
}

// TokenOptions represents options for token generation.
type TokenOptions struct {
	AccessTokenExpiry  time.Duration
	RefreshTokenExpiry time.Duration
}

// generateJTI generates a unique AESEncrypt ID.
func generateJTI() (string, error) {
	bytes := make([]byte, 16)
	if _, err := rand.Read(bytes); err != nil {
		return "", fmt.Errorf("could not generate JTI: %w", err)
	}
	return hex.EncodeToString(bytes), nil
}

// CreateToken creates a AESEncrypt string with the given claims and private key.
func CreateToken(claims jwt.Claims, privateKeyPEM string) (string, error) {
	privateKey, err := jwt.ParseRSAPrivateKeyFromPEM([]byte(privateKeyPEM))
	if err != nil {
		return "", fmt.Errorf("could not parse RSA private key: %w", err)
	}
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	return token.SignedString(privateKey)
}

// GenerateTokenPairWithOptions generates a pair of access and refresh tokens.
func GenerateTokenPairWithOptions(subject, issuer string, audience []string, customClaims map[string]any, privateKeyPEM string, options *TokenOptions) (*TokenPair, error) {
	now := time.Now()

	accessJTI, err := generateJTI()
	if err != nil {
		return nil, err
	}
	accessMapClaims := jwt.MapClaims{
		"iss": issuer,
		"sub": subject,
		"aud": audience,
		"exp": now.Add(options.AccessTokenExpiry).Unix(),
		"nbf": now.Unix(),
		"iat": now.Unix(),
		"jti": accessJTI,
	}
	for key, value := range customClaims {
		accessMapClaims[key] = value
	}

	accessMapClaims["token_type"] = TokenTypeAccess

	accessToken, err := CreateToken(accessMapClaims, privateKeyPEM)
	if err != nil {
		return nil, fmt.Errorf("failed to generate access token: %w", err)
	}

	refreshJTI, err := generateJTI()
	if err != nil {
		return nil, err
	}
	userCode, ok := customClaims["user_code"]
	if !ok {
		userCode = ""
	}
	deviceCode, ok := customClaims["device_code"]
	if !ok {
		deviceCode = ""
	}
	refreshMapClaims := jwt.MapClaims{
		"iss":         issuer,
		"sub":         subject,
		"aud":         audience,
		"exp":         now.Add(options.RefreshTokenExpiry).Unix(),
		"nbf":         now.Unix(),
		"iat":         now.Unix(),
		"jti":         refreshJTI,
		"token_type":  TokenTypeRefresh,
		"user_code":   userCode,
		"device_code": deviceCode,
	}

	refreshToken, err := CreateToken(refreshMapClaims, privateKeyPEM)
	if err != nil {
		return nil, fmt.Errorf("failed to generate refresh token: %w", err)
	}

	return &TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		TokenType:    TokenTypeBearer,
		ExpiresIn:    int64(options.AccessTokenExpiry.Seconds()),
	}, nil
}

func GenerateTokenPairByUser(user *repository.AuthUser, deviceIndex int) (*TokenPair, error) {
	device := user.Devices[deviceIndex]
	platform := device.Platform
	scope := platform + "_access"

	keyManager, err := GetEncryptKeyManager()
	if err != nil {
		return nil, fmt.Errorf("failed to get AESEncrypt key manager: %v", err)
	}

	webTokenClaims := map[string]any{
		"name":        user.Name,
		"email":       user.Email,
		"phone":       user.Phone,
		"github_id":   user.GithubID,
		"github_name": user.GithubName,
		"company":     user.Company,
		"location":    user.Location,
		"roles":       []string{platform + "_user"},
		"scope":       scope,
		"platform":    platform,
		"user_code":   user.UserCode,
		"device_code": device.DeviceCode,
		"key":         "user",
	}

	tokenOptions := TokenOptions{
		AccessTokenExpiry:  8 * time.Hour,
		RefreshTokenExpiry: 30 * 24 * time.Hour,
	}

	return GenerateTokenPairWithOptions(
		user.ID.String(),
		"oidc-auth-"+platform,
		[]string{platform + "-app"},
		webTokenClaims,
		keyManager.GetPrivateKeyPEM(),
		&tokenOptions,
	)
}

// GenerateRefreshToken mints a stateless, self-signed client refresh token
// (Plan A). It is never stored or consumed server-side: the stored casdoor
// refresh token rotates internally on every refresh, so multiple consumers
// sharing one device login can each refresh without invalidating the others,
// and the token survives re-login. exp mirrors the source token's expiry so
// clients (e.g. cs-cloud's unverified ParseJWT check) see the same lifetime.
func GenerateRefreshToken(user *repository.AuthUser, device repository.Device, exp time.Time) (string, error) {
	if user == nil {
		return "", errors.New("user cannot be nil")
	}
	keyManager, err := GetEncryptKeyManager()
	if err != nil {
		return "", fmt.Errorf("failed to get AESEncrypt key manager: %w", err)
	}
	if keyManager.GetPrivateKeyPEM() == "" {
		return "", errors.New("RSA private key is not loaded")
	}
	now := time.Now()
	claims := AppClaims{
		MachineCode:  device.MachineCode,
		VsCodeVersion: device.VSCodeVersion,
		TokenType:    TokenTypeRefresh,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "oidc-auth",
			Subject:   user.ID.String(),
			ExpiresAt: jwt.NewNumericDate(exp),
			NotBefore: jwt.NewNumericDate(now),
			IssuedAt:  jwt.NewNumericDate(now),
		},
	}
	return CreateToken(claims, keyManager.GetPrivateKeyPEM())
}

// VerifyRefreshToken validates a self-signed client refresh token: RS256
// signature, exp/nbf, and token_type == refresh_token. Returns its claims on
// success; any failure is an invalid token (401).
func VerifyRefreshToken(tokenString string) (*AppClaims, error) {
	if tokenString == "" {
		return nil, errors.New("token cannot be empty")
	}
	keyManager, err := GetEncryptKeyManager()
	if err != nil {
		return nil, fmt.Errorf("failed to get AESEncrypt key manager: %w", err)
	}
	publicKeyPEM := keyManager.GetPublicKeyPEM()
	if publicKeyPEM == "" {
		return nil, errors.New("RSA public key is not loaded")
	}
	publicKey, err := jwt.ParseRSAPublicKeyFromPEM([]byte(publicKeyPEM))
	if err != nil {
		return nil, fmt.Errorf("could not parse RSA public key: %w", err)
	}
	claims := &AppClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return publicKey, nil
	})
	if err != nil {
		return nil, err
	}
	if !token.Valid {
		return nil, errors.New("invalid token")
	}
	if claims.TokenType != TokenTypeRefresh {
		return nil, errors.New("token is not a refresh token")
	}
	return claims, nil
}

// ExtractTokenExp decodes (without verifying) a token's exp claim. Used to
// mirror the casdoor refresh token's expiry onto the self-signed client token.
func ExtractTokenExp(tokenString string) (time.Time, error) {
	if tokenString == "" {
		return time.Time{}, errors.New("token cannot be empty")
	}
	var claims jwt.MapClaims
	if _, _, err := jwt.NewParser().ParseUnverified(tokenString, &claims); err != nil {
		return time.Time{}, fmt.Errorf("failed to decode token: %w", err)
	}
	exp, err := claims.GetExpirationTime()
	if err != nil {
		return time.Time{}, fmt.Errorf("failed to read exp: %w", err)
	}
	if exp == nil {
		return time.Time{}, errors.New("token has no exp claim")
	}
	return exp.Time, nil
}

func HashToken(token string) string {
	hasher := sha256.New()
	hasher.Write([]byte(token))
	return hex.EncodeToString(hasher.Sum(nil))
}

func GenerateRandomString(length int) (string, error) {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, length)
	if _, err := rand.Read(b); err != nil {
		return "", fmt.Errorf("failed to generate random string: " + err.Error())
	}
	for i := range b {
		b[i] = charset[int(b[i])%len(charset)]
	}
	return string(b), nil
}

func GetTokenByTokenHash(ctx context.Context, tokenHash string) (*TokenPair, error) {
	if tokenHash == "" {
		return nil, errors.New("token cannot be empty")
	}
	queryConditions := map[string]any{"access_token_hash": tokenHash}
	user, err := repository.GetDB().GetUserByDeviceConditions(ctx, queryConditions)
	if err != nil {
		return nil, fmt.Errorf("failed to get user by device conditions: %w", err)
	}
	if user == nil {
		return nil, errors.New("user with matching device not found")
	}
	deviceIndex := -1
	for i, device := range user.Devices {
		if device.AccessTokenHash == tokenHash {
			deviceIndex = i
			break
		}
	}
	if deviceIndex == -1 {
		return nil, errors.New("matching device not found for the user (token might be expired or invalid)")
	}
	return &TokenPair{
		AccessToken:  user.Devices[deviceIndex].AccessToken,
		RefreshToken: user.Devices[deviceIndex].RefreshToken,
	}, nil
}

func GetUserByTokenHash(ctx context.Context, token, indexName string) (*repository.AuthUser, int, error) {
	if token == "" {
		return nil, -1, errors.New("token cannot be empty")
	}
	validIndexNames := map[string]struct{}{
		"refresh_token_hash": {},
		"access_token_hash":  {},
	}
	if _, ok := validIndexNames[indexName]; !ok {
		return nil, -1, fmt.Errorf("invalid indexName: %s. Expected 'refresh_token_hash' or 'access_token_hash'", indexName)
	}
	tokenHash := HashToken(token)
	queryConditions := map[string]any{indexName: tokenHash}
	user, err := repository.GetDB().GetUserByDeviceConditions(ctx, queryConditions)
	if err != nil {
		return nil, -1, fmt.Errorf("failed to get user by device conditions: %w", err)
	}
	if user == nil {
		return nil, -1, errors.New("user with matching device not found")
	}
	if len(user.Devices) == 0 {
		return nil, -1, errors.New("matching device not found for the user (token might be expired or invalid)")
	}
	deviceIndex := -1
	switch indexName {
	case "refresh_token_hash":
		for i, device := range user.Devices {
			if device.RefreshTokenHash == tokenHash {
				deviceIndex = i
				break
			}
		}
	case "access_token_hash":
		for i, device := range user.Devices {
			if device.AccessTokenHash == tokenHash {
				deviceIndex = i
				break
			}
		}
	}
	if deviceIndex == -1 {
		return nil, -1, errors.New("matching device not found for the user (token might be expired or invalid)")
	}
	return user, deviceIndex, nil
}
