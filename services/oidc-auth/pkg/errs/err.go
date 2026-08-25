package errs

import "errors"

var (
	ErrInfoInvalidToken   = errors.New("invalid token, no related account exists")
	ErrInfoQueryUserInfo  = errors.New("query user info fail")
	ErrInfoUpdateUserInfo = errors.New("update user info fail")
	ErrInfoGenerateToken  = errors.New("generate token fail")
)

const (
	ErrBadRequestParam = "oidc-auth.badRequestParameter"
	ErrDataEncryption  = "oidc-auth.loginEncryptFailed"
	ErrDataDecryption  = "oidc-auth.LoginDecryptFailed"
	ErrUserNotFound    = "oidc-auth.userNotFound"
	ErrTokenInvalid    = "oidc-auth.tokenInvalid"
	ErrUpdateInfo      = "oidc-auth.updateInfoFailed"
	ErrBindAccount     = "oidc-auth.bindAccountFailed"
	ErrTokenGenerate   = "oidc-auth.tokenGenerateFailed"
	ErrAuthentication  = "oidc-auth.authenticationFailed"
	// ErrInvalidIdentity: cs-user rejected the identity (401) — login fails.
	ErrInvalidIdentity = "oidc-auth.invalidIdentity"
	// ErrServiceUnavailable: cs-user unreachable / failed (network, timeout,
	// 5xx) — login cannot establish the identity trust boundary.
	ErrServiceUnavailable = "oidc-auth.serviceUnavailable"
)

func ParamNeedErr(name string) error {
	return errors.New(name + " needs to be provided")
}
