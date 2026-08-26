package usercenter

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"
)

// stubServer is an httptest server with per-path canned responses and a
// request recorder for asserting the X-Internal-Token header.
type stubServer struct {
	server   *httptest.Server
	mu       sync.Mutex
	lastPath string
	lastHdr  http.Header
}

func newStubServer(handler http.HandlerFunc) *stubServer {
	s := &stubServer{}
	s.server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		s.mu.Lock()
		s.lastPath = r.URL.Path
		s.lastHdr = r.Header.Clone()
		s.mu.Unlock()
		handler(w, r)
	}))
	return s
}

func (s *stubServer) close() { s.server.Close() }

func (s *stubServer) headers() http.Header {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.lastHdr
}

func newTestClient(baseURL string, timeout time.Duration) *Client {
	return NewClient(baseURL, "test-internal-token", http.DefaultClient, timeout)
}

func jsonRespond(w http.ResponseWriter, code int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(v)
}

func TestVerify_Success(t *testing.T) {
	s := newStubServer(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost || r.URL.Path != "/api/internal/auth/verify" {
			http.NotFound(w, r)
			return
		}
		var req map[string]string
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req["token"] == "" {
			jsonRespond(w, http.StatusBadRequest, ginH("token required"))
			return
		}
		jsonRespond(w, http.StatusOK, map[string]any{
			"active": true, "token_source": "casdoor", "sub": "sub-1",
			"universal_id": "universal-123", "name": "Zhang San",
			"email": "zs@example.com", "phone": "+8613800000000", "iss": "casdoor.example.com",
			"reissued_token": "eyJ-reissued", "is_new_user": true,
		})
	})
	defer s.close()

	result, err := newTestClient(s.server.URL, time.Second).Verify(context.Background(), "jwt-token")
	if err != nil {
		t.Fatalf("Verify: %v", err)
	}
	if !result.Active {
		t.Error("Active = false, want true")
	}
	if result.UniversalID != "universal-123" || result.Name != "Zhang San" || result.Subject != "sub-1" {
		t.Errorf("result = %+v, want universal_id=universal-123 name=Zhang San sub=sub-1", result)
	}
	if result.ReissuedToken != "eyJ-reissued" {
		t.Errorf("ReissuedToken = %q, want passthrough eyJ-reissued", result.ReissuedToken)
	}
	if !result.IsNewUser {
		t.Error("IsNewUser = false, want true (verify reports inline provisioning)")
	}
	if got := s.headers().Get("X-Internal-Token"); got != "test-internal-token" {
		t.Errorf("X-Internal-Token = %q, want test-internal-token", got)
	}
}

func TestVerify_Rejected_ReturnsErrInvalidIdentity(t *testing.T) {
	s := newStubServer(func(w http.ResponseWriter, r *http.Request) {
		jsonRespond(w, http.StatusUnauthorized, map[string]any{"active": false, "error": "invalid token"})
	})
	defer s.close()

	_, err := newTestClient(s.server.URL, time.Second).Verify(context.Background(), "bad-jwt")
	if !errors.Is(err, ErrInvalidIdentity) {
		t.Fatalf("err = %v, want ErrInvalidIdentity", err)
	}
}

func TestVerify_Inactive_ReturnsErrInvalidIdentity(t *testing.T) {
	s := newStubServer(func(w http.ResponseWriter, r *http.Request) {
		jsonRespond(w, http.StatusOK, map[string]any{"active": false, "error": "invalid token"})
	})
	defer s.close()

	_, err := newTestClient(s.server.URL, time.Second).Verify(context.Background(), "bad-jwt")
	if !errors.Is(err, ErrInvalidIdentity) {
		t.Fatalf("err = %v, want ErrInvalidIdentity", err)
	}
}

func TestVerify_Unconfigured_ReturnsErrServiceUnavailable(t *testing.T) {
	s := newStubServer(func(w http.ResponseWriter, r *http.Request) {
		jsonRespond(w, http.StatusServiceUnavailable, ginH("casdoor verifier not configured"))
	})
	defer s.close()

	_, err := newTestClient(s.server.URL, time.Second).Verify(context.Background(), "jwt")
	if !errors.Is(err, ErrServiceUnavailable) {
		t.Fatalf("err = %v, want ErrServiceUnavailable", err)
	}
}

func TestVerify_NetworkError_ReturnsErrServiceUnavailable(t *testing.T) {
	s := newStubServer(func(w http.ResponseWriter, r *http.Request) {})
	s.close() // server down → connection refused

	_, err := newTestClient(s.server.URL, time.Second).Verify(context.Background(), "jwt")
	if !errors.Is(err, ErrServiceUnavailable) {
		t.Fatalf("err = %v, want ErrServiceUnavailable", err)
	}
}

func TestVerify_Timeout_ReturnsErrServiceUnavailable(t *testing.T) {
	s := newStubServer(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(200 * time.Millisecond)
		jsonRespond(w, http.StatusOK, ginH("{}"))
	})
	defer s.close()

	start := time.Now()
	_, err := newTestClient(s.server.URL, 50*time.Millisecond).Verify(context.Background(), "jwt")
	if !errors.Is(err, ErrServiceUnavailable) {
		t.Fatalf("err = %v, want ErrServiceUnavailable", err)
	}
	if elapsed := time.Since(start); elapsed > 150*time.Millisecond {
		t.Errorf("per-call timeout not enforced, elapsed = %v", elapsed)
	}
}

func TestGetOrCreate_ReturnsUserAndIsNewUser(t *testing.T) {
	s := newStubServer(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost || r.URL.Path != "/api/internal/users/get-or-create" {
			http.NotFound(w, r)
			return
		}
		var req Claims
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.UniversalID == "" {
			jsonRespond(w, http.StatusBadRequest, ginH("no valid user identifier"))
			return
		}
		jsonRespond(w, http.StatusOK, map[string]any{
			"user":        map[string]any{"subject_id": "subject-456", "username": "zs", "email": "zs@example.com"},
			"is_new_user": true,
		})
	})
	defer s.close()

	user, isNew, err := newTestClient(s.server.URL, time.Second).GetOrCreate(context.Background(), &Claims{UniversalID: "universal-123"})
	if err != nil {
		t.Fatalf("GetOrCreate: %v", err)
	}
	if !isNew {
		t.Error("is_new_user = false, want true")
	}
	if user.SubjectID != "subject-456" || user.Username != "zs" {
		t.Errorf("user = %+v, want subject_id=subject-456 username=zs", user)
	}
}

func TestGetOrCreate_ExplicitlyUnbound_ReturnsErrExplicitlyUnbound(t *testing.T) {
	s := newStubServer(func(w http.ResponseWriter, r *http.Request) {
		jsonRespond(w, http.StatusConflict, map[string]any{"error": "explicitly_unbound", "error_code": "explicitly_unbound"})
	})
	defer s.close()

	_, _, err := newTestClient(s.server.URL, time.Second).GetOrCreate(context.Background(), &Claims{UniversalID: "u"})
	if !errors.Is(err, ErrExplicitlyUnbound) {
		t.Fatalf("err = %v, want ErrExplicitlyUnbound", err)
	}
}

func TestListIdentities_Empty(t *testing.T) {
	s := newStubServer(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/internal/users/subject-456/auth-identities" {
			http.NotFound(w, r)
			return
		}
		jsonRespond(w, http.StatusOK, map[string]any{"identities": []any{}})
	})
	defer s.close()

	ids, err := newTestClient(s.server.URL, time.Second).ListIdentities(context.Background(), "subject-456")
	if err != nil {
		t.Fatalf("ListIdentities: %v", err)
	}
	if len(ids) != 0 {
		t.Errorf("identities = %+v, want empty", ids)
	}
}

func TestListIdentities_GithubRow(t *testing.T) {
	s := newStubServer(func(w http.ResponseWriter, r *http.Request) {
		jsonRespond(w, http.StatusOK, map[string]any{
			"identities": []any{
				map[string]any{"provider": "github", "provider_user_id": "gh-999", "display_name": "octocat"},
				map[string]any{"provider": "phone", "phone": "+8613800000000"},
			},
		})
	})
	defer s.close()

	ids, err := newTestClient(s.server.URL, time.Second).ListIdentities(context.Background(), "subject-456")
	if err != nil {
		t.Fatalf("ListIdentities: %v", err)
	}
	if len(ids) != 2 {
		t.Fatalf("identities = %+v, want 2 rows", ids)
	}
	if ids[0].Provider != "github" || ids[0].ProviderUserID != "gh-999" || ids[0].DisplayName == nil || *ids[0].DisplayName != "octocat" {
		t.Errorf("github row = %+v, want provider=github provider_user_id=gh-999 display_name=octocat", ids[0])
	}
}

func TestGetProfile_Success(t *testing.T) {
	s := newStubServer(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/internal/users/subject-456/profile" {
			http.NotFound(w, r)
			return
		}
		jsonRespond(w, http.StatusOK, map[string]any{
			"subject_id": "subject-456", "username": "zs", "display_name": "Zhang San",
			"email": "zs@example.com", "phone": "13800000000",
		})
	})
	defer s.close()

	profile, err := newTestClient(s.server.URL, time.Second).GetProfile(context.Background(), "subject-456")
	if err != nil {
		t.Fatalf("GetProfile: %v", err)
	}
	if profile.Username != "zs" || profile.Email == nil || *profile.Email != "zs@example.com" {
		t.Errorf("profile = %+v, want username=zs email=zs@example.com", profile)
	}
}

func TestGetProfile_NotFound_ReturnsHTTPError(t *testing.T) {
	s := newStubServer(func(w http.ResponseWriter, r *http.Request) {
		jsonRespond(w, http.StatusNotFound, ginH("user not found"))
	})
	defer s.close()

	_, err := newTestClient(s.server.URL, time.Second).GetProfile(context.Background(), "missing")
	var httpErr *HTTPError
	if !errors.As(err, &httpErr) || httpErr.StatusCode != http.StatusNotFound {
		t.Fatalf("err = %v, want *HTTPError with 404", err)
	}
}

func TestInitClient_MissingBaseURL_ReturnsError(t *testing.T) {
	if err := InitClient("", "tok", http.DefaultClient, time.Second); err == nil {
		t.Fatal("InitClient with empty baseURL: want error (fail-closed), got nil")
	}
}

func ginH(v string) map[string]any { return map[string]any{"error": v} }
