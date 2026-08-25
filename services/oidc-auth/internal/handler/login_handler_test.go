package handler

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/zgsm-ai/oidc-auth/internal/usercenter"
	"github.com/zgsm-ai/oidc-auth/pkg/errs"
)

// chainStub is a single cs-user login-chain stub kept alive for the whole test
// binary (usercenter.InitClient is sync.Once-guarded, so every test shares one
// client URL). Behavior is driven by the token the client sends:
//
//	"bad-jwt"  → verify 401
//	"no-uuid"  → verify active=true without universal_id
//	"nogithub" → get-or-create returns subject-nogithub, no github identity row
//	otherwise  → full success path
func chainStub() *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.Method == http.MethodPost && r.URL.Path == "/api/internal/auth/verify":
			var req map[string]string
			_ = json.NewDecoder(r.Body).Decode(&req)
			switch req["token"] {
			case "bad-jwt":
				jsonRespond(w, http.StatusUnauthorized, map[string]any{"active": false, "error": "invalid token"})
				return
			case "no-uuid":
				jsonRespond(w, http.StatusOK, map[string]any{
					"active": true,
					"sub":    "sub-no-uuid",
				})
				return
			case "nogithub":
				jsonRespond(w, http.StatusOK, map[string]any{
					"active": true, "sub": "sub-nogithub",
					"universal_id": "22222222-2222-2222-2222-222222222222",
					"name":         "No GitHub", "email": "ng@example.com",
				})
				return
			case "profilefail":
				jsonRespond(w, http.StatusOK, map[string]any{
					"active": true, "sub": "sub-fail",
					"universal_id": "33333333-3333-3333-3333-333333333333",
					"name":         "Profile Fail", "email": "pf@example.com",
				})
				return
			}
			jsonRespond(w, http.StatusOK, map[string]any{
				"active": true, "sub": "sub-1",
				"universal_id": "11111111-1111-1111-1111-111111111111",
				"name":         "Zhang San", "email": "zs@example.com", "phone": "+8613800000000",
			})
		case r.Method == http.MethodPost && r.URL.Path == "/api/internal/users/get-or-create":
			var req usercenter.Claims
			_ = json.NewDecoder(r.Body).Decode(&req)
			subjectID := "subject-1"
			if req.UniversalID == "22222222-2222-2222-2222-222222222222" {
				subjectID = "subject-nogithub"
			} else if req.UniversalID == "33333333-3333-3333-3333-333333333333" {
				subjectID = "subject-fail"
			}
			jsonRespond(w, http.StatusOK, map[string]any{
				"user":        map[string]any{"subject_id": subjectID, "username": "zs"},
				"is_new_user": true,
			})
		case r.Method == http.MethodGet && r.URL.Path == "/api/internal/users/subject-1/profile":
			jsonRespond(w, http.StatusOK, map[string]any{
				"subject_id": "subject-1", "username": "zs-updated",
				"display_name": "Zhang San Updated", "email": "zs-updated@example.com", "phone": "+8613800001111",
			})
		case r.Method == http.MethodGet && r.URL.Path == "/api/internal/users/subject-fail/profile":
			jsonRespond(w, http.StatusInternalServerError, map[string]any{"error": "profile unavailable"})
		case r.Method == http.MethodGet && r.URL.Path == "/api/internal/users/subject-nogithub/auth-identities":
			jsonRespond(w, http.StatusOK, map[string]any{"identities": []any{}})
		case r.Method == http.MethodGet && r.URL.Path == "/api/internal/users/subject-fail/auth-identities":
			jsonRespond(w, http.StatusOK, map[string]any{
				"identities": []any{
					map[string]any{"provider": "github", "provider_user_id": "gh-2", "display_name": "dromedary"},
				},
			})
		case r.Method == http.MethodGet && r.URL.Path == "/api/internal/users/subject-1/auth-identities":
			jsonRespond(w, http.StatusOK, map[string]any{
				"identities": []any{
					map[string]any{"provider": "github", "provider_user_id": "gh-1", "display_name": "octocat"},
				},
			})
		default:
			http.NotFound(w, r)
		}
	}))
}

func TestMain(m *testing.M) {
	s := chainStub()
	if err := usercenter.InitClient(s.URL, "test-token", http.DefaultClient, time.Second); err != nil {
		fmt.Println("InitClient:", err)
		os.Exit(1)
	}
	code := m.Run()
	s.Close()
	os.Exit(code)
}

func TestFetchUserByIdentityChain_Success(t *testing.T) {
	user, isNew, err := fetchUserByIdentityChain(context.Background(), "ok-token")
	if err != nil {
		t.Fatalf("fetchUserByIdentityChain: %v", err)
	}
	if !isNew {
		t.Error("isNewUser = false, want true")
	}
	if user.ID.String() != "11111111-1111-1111-1111-111111111111" {
		t.Errorf("ID = %s, want universal_id", user.ID)
	}
	if user.Name != "Zhang San" || user.Email != "zs@example.com" {
		t.Errorf("identity fields = name:%q email:%q, want Zhang San / zs@example.com", user.Name, user.Email)
	}
	if user.Phone != "13800000000" {
		t.Errorf("phone = %q, want +86 prefix stripped", user.Phone)
	}
	if user.GithubID != "gh-1" || user.GithubName != "octocat" {
		t.Errorf("github = id:%q name:%q, want gh-1 / octocat", user.GithubID, user.GithubName)
	}
	if user.SubjectID != "subject-1" {
		t.Errorf("SubjectID = %q, want subject-1", user.SubjectID)
	}
	if user.IdentitySyncedAt == nil {
		t.Error("IdentitySyncedAt = nil, want set")
	}
}

func TestFetchUserByIdentityChain_NoGithubRow(t *testing.T) {
	user, _, err := fetchUserByIdentityChain(context.Background(), "nogithub")
	if err != nil {
		t.Fatalf("fetchUserByIdentityChain: %v", err)
	}
	if user.GithubID != "" || user.GithubName != "" {
		t.Errorf("github = id:%q name:%q, want empty", user.GithubID, user.GithubName)
	}
	if user.SubjectID != "subject-nogithub" {
		t.Errorf("SubjectID = %q, want subject-nogithub", user.SubjectID)
	}
}

func TestFetchUserByIdentityChain_IdentityRejected(t *testing.T) {
	_, _, err := fetchUserByIdentityChain(context.Background(), "bad-jwt")
	if err == nil {
		t.Fatal("want error for rejected identity, got nil")
	}
	if !errors.Is(err, usercenter.ErrInvalidIdentity) {
		t.Errorf("err = %v, want ErrInvalidIdentity", err)
	}
	if code, _ := mapUserCenterError(err); code != http.StatusUnauthorized {
		t.Errorf("mapUserCenterError status = %d, want 401", code)
	}
}

func TestFetchUserByIdentityChain_InvalidUniversalID(t *testing.T) {
	if _, _, err := fetchUserByIdentityChain(context.Background(), "no-uuid"); err == nil {
		t.Fatal("want error for empty universal_id, got nil")
	}
}

func TestMapUserCenterError(t *testing.T) {
	cases := []struct {
		name string
		err  error
		want int
	}{
		{"invalid identity", fmt.Errorf("%w: x", usercenter.ErrInvalidIdentity), http.StatusUnauthorized},
		{"explicitly unbound", fmt.Errorf("%w: x", usercenter.ErrExplicitlyUnbound), http.StatusConflict},
		{"unavailable", fmt.Errorf("%w: x", usercenter.ErrServiceUnavailable), http.StatusServiceUnavailable},
		{"generic http error", &usercenter.HTTPError{StatusCode: http.StatusBadRequest, Body: "x"}, http.StatusServiceUnavailable},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got, _ := mapUserCenterError(tc.err); got != tc.want {
				t.Errorf("status = %d, want %d", got, tc.want)
			}
		})
	}
}

func TestMapUserCenterError_WrappedChainError(t *testing.T) {
	// The production path wraps chain errors with %w in GetUserByOauth /
	// GetWebUserByOauth; classification must survive the wrap.
	wrapped := fmt.Errorf("%s: %w", errs.ErrInfoQueryUserInfo, fmt.Errorf("%w: x", usercenter.ErrInvalidIdentity))
	if code, _ := mapUserCenterError(wrapped); code != http.StatusUnauthorized {
		t.Errorf("wrapped invalid-identity status = %d, want 401", code)
	}
}

func TestNormalizePhone(t *testing.T) {
	if got := normalizePhone("+8613800000000"); got != "13800000000" {
		t.Errorf("normalizePhone(+8613800000000) = %q", got)
	}
	if got := normalizePhone("13800000000"); got != "13800000000" {
		t.Errorf("normalizePhone(13800000000) = %q", got)
	}
}

func TestFetchUserInfoFromCsUser_Success(t *testing.T) {
	user, err := fetchUserInfoFromCsUser(context.Background(), "ok-token")
	if err != nil {
		t.Fatalf("fetchUserInfoFromCsUser: %v", err)
	}
	if user.ID.String() != "11111111-1111-1111-1111-111111111111" {
		t.Errorf("ID = %s, want universal_id", user.ID)
	}
	if user.Name != "zs-updated" {
		t.Errorf("Name = %q, want profile username zs-updated", user.Name)
	}
	if user.Email != "zs-updated@example.com" {
		t.Errorf("Email = %q, want profile email", user.Email)
	}
	if user.Phone != "13800001111" {
		t.Errorf("Phone = %q, want profile phone with +86 stripped", user.Phone)
	}
	if user.GithubID != "gh-1" || user.GithubName != "octocat" {
		t.Errorf("github = id:%q name:%q, want gh-1 / octocat", user.GithubID, user.GithubName)
	}
}

func TestFetchUserInfoFromCsUser_SoftProfileFailure(t *testing.T) {
	// GetProfile failing is a soft read: login-time claims are served, no error.
	user, err := fetchUserInfoFromCsUser(context.Background(), "profilefail")
	if err != nil {
		t.Fatalf("fetchUserInfoFromCsUser: %v", err)
	}
	if user.Name != "Profile Fail" || user.Email != "pf@example.com" {
		t.Errorf("login-time claims = name:%q email:%q, want Profile Fail / pf@example.com", user.Name, user.Email)
	}
	// ListIdentities still applies even after the profile read failed.
	if user.GithubID != "gh-2" || user.GithubName != "dromedary" {
		t.Errorf("github = id:%q name:%q, want gh-2 / dromedary", user.GithubID, user.GithubName)
	}
}

func TestFetchUserInfoFromCsUser_InvalidIdentity(t *testing.T) {
	if _, err := fetchUserInfoFromCsUser(context.Background(), "bad-jwt"); !errors.Is(err, usercenter.ErrInvalidIdentity) {
		t.Errorf("err = %v, want ErrInvalidIdentity", err)
	}
}

func jsonRespond(w http.ResponseWriter, code int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(v)
}
