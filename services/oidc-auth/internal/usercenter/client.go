// Package usercenter is the cs-user internal API client. It backs the login
// chain (parse-identity → get-or-create → auth-identities) and the userinfo
// soft-TTL refresh (get-profile). All endpoints are guarded by the shared
// X-Internal-Token secret; oidc-auth never sends X-Tenant-Id, so every call
// lands on the default tenant.
package usercenter

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"sync"
	"time"
)

// Error sentinels used by handlers to map RPC failures:
//   - ErrInvalidIdentity: cs-user rejected the identity (401 verify failure) → handler 401
//   - ErrExplicitlyUnbound: get-or-create refused to silently rebuild an
//     explicitly unbound identity (409) → login failure
//   - ErrServiceUnavailable: network / timeout / 5xx / 503 → handler 503
var (
	ErrInvalidIdentity    = errors.New("usercenter: invalid identity")
	ErrExplicitlyUnbound  = errors.New("usercenter: identity explicitly unbound")
	ErrServiceUnavailable = errors.New("usercenter: service unavailable")
)

// HTTPError is an unexpected non-2xx status from cs-user (a caller bug or an
// unmapped status). Handlers treat it with ErrServiceUnavailable semantics.
type HTTPError struct {
	StatusCode int
	Body       string
}

func (e *HTTPError) Error() string {
	return fmt.Sprintf("usercenter: unexpected status %d: %s", e.StatusCode, e.Body)
}

// IdentityProfile is the verified IdP profile returned by parse-identity
// (mirrors cs-user's reissueProfile wire shape).
type IdentityProfile struct {
	ID                string `json:"id,omitempty"`
	Sub               string `json:"sub,omitempty"`
	UniversalID       string `json:"universal_id,omitempty"`
	Name              string `json:"name,omitempty"`
	PreferredUsername string `json:"preferred_username,omitempty"`
	Email             string `json:"email,omitempty"`
	Phone             string `json:"phone,omitempty"`
	Picture           string `json:"picture,omitempty"`
	Owner             string `json:"owner,omitempty"`
	Provider          string `json:"provider,omitempty"`
	ProviderUserID    string `json:"provider_user_id,omitempty"`
}

// Claims is the get-or-create request body, mirroring cs-user's
// models.JWTClaims JSON tags.
type Claims struct {
	ID                string `json:"id,omitempty"`
	Sub               string `json:"sub,omitempty"`
	UniversalID       string `json:"universal_id,omitempty"`
	Name              string `json:"name,omitempty"`
	PreferredUsername string `json:"preferred_username,omitempty"`
	Email             string `json:"email,omitempty"`
	Phone             string `json:"phone,omitempty"`
	Picture           string `json:"picture,omitempty"`
	Owner             string `json:"owner,omitempty"`
	Provider          string `json:"provider,omitempty"`
	ProviderUserID    string `json:"provider_user_id,omitempty"`
}

// User is the get-or-create response user row, carrying only the fields
// oidc-auth consumes (SubjectID is the TTL-refresh anchor).
type User struct {
	SubjectID   string  `json:"subject_id,omitempty"`
	Username    string  `json:"username,omitempty"`
	DisplayName *string `json:"display_name,omitempty"`
	Email       *string `json:"email,omitempty"`
	Phone       *string `json:"phone,omitempty"`
	Status      string  `json:"status,omitempty"`
}

// Identity is one row of the auth-identities list. The github row's
// provider_user_id / display_name map to AuthUser.GithubID / GithubName.
type Identity struct {
	Provider       string  `json:"provider,omitempty"`
	ProviderUserID string  `json:"provider_user_id,omitempty"`
	DisplayName    *string `json:"display_name,omitempty"`
	Email          *string `json:"email,omitempty"`
	Phone          *string `json:"phone,omitempty"`
}

// Profile is the user profile returned by get-profile, the soft-TTL refresh
// source (username / display_name / email / phone are user-owned on cs-user
// side and more authoritative than the login-time Casdoor claims).
type Profile struct {
	SubjectID   string  `json:"subject_id,omitempty"`
	Username    string  `json:"username,omitempty"`
	DisplayName *string `json:"display_name,omitempty"`
	Email       *string `json:"email,omitempty"`
	Phone       *string `json:"phone,omitempty"`
}

// Client talks to cs-user's internal API.
type Client struct {
	baseURL       string
	internalToken string
	httpClient    *http.Client
	timeout       time.Duration
}

// NewClient builds a client. timeout bounds a single RPC call; the login
// callback context caps the whole three-call chain at 15s, so keep it 3-5s.
func NewClient(baseURL, internalToken string, httpClient *http.Client, timeout time.Duration) *Client {
	if timeout <= 0 {
		timeout = 5 * time.Second
	}
	return &Client{
		baseURL:       baseURL,
		internalToken: internalToken,
		httpClient:    httpClient,
		timeout:       timeout,
	}
}

// ParseIdentity forwards the raw Casdoor JWT to cs-user for verification and
// returns the verified profile. ErrInvalidIdentity on verify failure.
func (c *Client) ParseIdentity(ctx context.Context, rawJWT string) (*IdentityProfile, error) {
	body, err := json.Marshal(map[string]string{"token": rawJWT})
	if err != nil {
		return nil, err
	}
	var out struct {
		Profile *IdentityProfile `json:"profile"`
	}
	if err := c.doJSON(ctx, http.MethodPost, "/api/internal/auth/parse-identity", body, &out); err != nil {
		return nil, err
	}
	if out.Profile == nil {
		return nil, ErrInvalidIdentity
	}
	return out.Profile, nil
}

// GetOrCreate upserts the user from parsed Casdoor claims and reports whether
// the row was actually created (is_new_user). ErrExplicitlyUnbound when cs-user
// refuses to rebuild an explicitly unbound identity.
func (c *Client) GetOrCreate(ctx context.Context, claims *Claims) (*User, bool, error) {
	body, err := json.Marshal(claims)
	if err != nil {
		return nil, false, err
	}
	var out struct {
		User      *User `json:"user"`
		IsNewUser bool  `json:"is_new_user"`
	}
	if err := c.doJSON(ctx, http.MethodPost, "/api/internal/users/get-or-create", body, &out); err != nil {
		return nil, false, err
	}
	if out.User == nil {
		return nil, false, ErrServiceUnavailable
	}
	return out.User, out.IsNewUser, nil
}

// ListIdentities returns the auth identities bound to subjectID. The github
// row supplies GithubID (provider_user_id) and GithubName (display_name).
func (c *Client) ListIdentities(ctx context.Context, subjectID string) ([]Identity, error) {
	var out struct {
		Identities []Identity `json:"identities"`
	}
	path := "/api/internal/users/" + url.PathEscape(subjectID) + "/auth-identities"
	if err := c.doJSON(ctx, http.MethodGet, path, nil, &out); err != nil {
		return nil, err
	}
	return out.Identities, nil
}

// GetProfile returns the user's identity profile, used by the userinfo
// soft-TTL refresh.
func (c *Client) GetProfile(ctx context.Context, subjectID string) (*Profile, error) {
	var out Profile
	path := "/api/internal/users/" + url.PathEscape(subjectID) + "/profile"
	if err := c.doJSON(ctx, http.MethodGet, path, nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// doJSON performs one RPC and maps failures to the package error sentinels:
// 401 → ErrInvalidIdentity, 409 → ErrExplicitlyUnbound, other non-2xx →
// *HTTPError, transport errors / timeout → ErrServiceUnavailable.
func (c *Client) doJSON(ctx context.Context, method, path string, body []byte, out any) error {
	ctx, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel()

	var reader io.Reader
	if body != nil {
		reader = bytes.NewReader(body)
	}
	req, err := http.NewRequestWithContext(ctx, method, c.baseURL+path, reader)
	if err != nil {
		return fmt.Errorf("%w: build request: %v", ErrServiceUnavailable, err)
	}
	req.Header.Set("Content-Type", "application/json")
	if c.internalToken != "" {
		req.Header.Set("X-Internal-Token", c.internalToken)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrServiceUnavailable, err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("%w: read response: %v", ErrServiceUnavailable, err)
	}

	switch {
	case resp.StatusCode == http.StatusOK:
		if out == nil {
			return nil
		}
		if err := json.Unmarshal(respBody, out); err != nil {
			return fmt.Errorf("%w: decode response: %v", ErrServiceUnavailable, err)
		}
		return nil
	case resp.StatusCode == http.StatusUnauthorized:
		return fmt.Errorf("%w: %s", ErrInvalidIdentity, truncateBody(respBody))
	case resp.StatusCode == http.StatusConflict:
		return fmt.Errorf("%w: %s", ErrExplicitlyUnbound, truncateBody(respBody))
	case resp.StatusCode >= http.StatusInternalServerError:
		// 5xx (incl. 503 verifier-unconfigured) → cs-user unavailable.
		return fmt.Errorf("%w: status %d: %s", ErrServiceUnavailable, resp.StatusCode, truncateBody(respBody))
	default:
		return &HTTPError{StatusCode: resp.StatusCode, Body: truncateBody(respBody)}
	}
}

func truncateBody(b []byte) string {
	const max = 512
	if len(b) > max {
		return string(b[:max])
	}
	return string(b)
}

var (
	clientInstance *Client
	clientOnce     sync.Once
)

// InitClient configures the package-global client used by handlers. Fails when
// BaseURL is missing (fail-closed: without cs-user the login chain cannot
// establish the identity trust boundary, so the process must not start).
func InitClient(baseURL, internalToken string, httpClient *http.Client, timeout time.Duration) error {
	if baseURL == "" {
		return errors.New("usercenter: baseURL is required")
	}
	clientOnce.Do(func() {
		clientInstance = NewClient(baseURL, internalToken, httpClient, timeout)
	})
	return nil
}

// GetClient returns the package-global client, or nil when InitClient has not
// run (callers must treat nil as a startup misconfiguration).
func GetClient() *Client {
	return clientInstance
}
