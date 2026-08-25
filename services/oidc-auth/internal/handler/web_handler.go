package handler

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/zgsm-ai/oidc-auth/internal/constants"
	"github.com/zgsm-ai/oidc-auth/internal/providers"
	"github.com/zgsm-ai/oidc-auth/internal/repository"
	"github.com/zgsm-ai/oidc-auth/pkg/errs"
	"github.com/zgsm-ai/oidc-auth/pkg/response"
	"github.com/zgsm-ai/oidc-auth/pkg/utils"
)

// WebParameterCarrier carries web login parameters through the OAuth flow
type WebParameterCarrier struct {
	Provider string `json:"provider"`
}

// webLoginHandler handles web login requests
func (s *Server) webLoginHandler(c *gin.Context) {
	provider := c.DefaultQuery("provider", "casdoor")
	inviterCode := c.DefaultQuery("inviter_code", "")
	redirectService := c.DefaultQuery("redirect_service", "")

	oauthManager := providers.GetManager()
	providerInstance, err := oauthManager.GetProvider(provider)
	if err != nil {
		response.HandleError(c, http.StatusInternalServerError, errs.ErrBadRequestParam, err)
		return
	}

	// Use inviterCode as state parameter
	state := inviterCode
	var redirectURL = ""
	if redirectService != "" {
		redirectURL = fmt.Sprintf("%s/%s", s.BaseURL+constants.WebLoginCallbackURI, redirectService)
	} else {
		redirectURL = s.BaseURL + constants.WebLoginCallbackURI
	}
	authURL := providerInstance.GetAuthURL(state, redirectURL)
	response.JSONSuccess(c, "", map[string]interface{}{
		"state":        state,
		"inviter_code": inviterCode,
		"url":          authURL,
	})
}

// webLoginCallbackHandler handles web login callback with invite code processing
func (s *Server) webLoginCallbackHandler(c *gin.Context) {
	code := c.DefaultQuery("code", "")
	state := c.DefaultQuery("state", "")
	inviterCode := state          // inviter code is in the state parameter
	service := c.Param("service") // service parameter for custom redirect

	if code == "" {
		response.JSONError(c, http.StatusBadRequest, errs.ErrBadRequestParam,
			errs.ParamNeedErr("code").Error())
		return
	}

	provider := "casdoor" // Fixed to use casdoor

	oauthManager := providers.GetManager()
	providerInstance, err := oauthManager.GetProvider(provider)
	if err != nil {
		response.HandleError(c, http.StatusBadRequest, errs.ErrBadRequestParam, err)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	// Get user info from OAuth provider
	user, isNewUser, err := GetWebUserByOauth(ctx, code, provider)
	if err != nil {
		status, codeMsg := mapUserCenterError(err)
		response.HandleError(c, status, codeMsg,
			fmt.Errorf("%s: %v", errs.ErrInfoQueryUserInfo, err))
		return
	}

	if user == nil {
		response.HandleError(c, http.StatusUnauthorized, errs.ErrTokenInvalid, errs.ErrInfoInvalidToken)
		return
	}

	// Inviter code only binds on cs-user-confirmed new users (decision 9);
	// existing users ignore it and complete login normally.
	if inviterCode != "" && isNewUser {
		inviter, err := utils.ValidateInviteCode(ctx, inviterCode)
		if err != nil {
			response.HandleError(c, http.StatusInternalServerError, errs.ErrBadRequestParam, err)
			return
		}
		user.InviterID = &inviter.ID
	}

	// Update or create user
	err = providerInstance.Update(ctx, user)
	if err != nil {
		response.HandleError(c, http.StatusInternalServerError, errs.ErrUpdateInfo,
			fmt.Errorf("%s: %v", errs.ErrInfoUpdateUserInfo, err))
		return
	}

	// Get user's access token hash as state parameter
	var tokenHash string
	if len(user.Devices) > 0 {
		tokenHash = user.Devices[0].AccessTokenHash
	}

	// Determine redirect URL based on service parameter
	var redirectURL string
	if service != "" {
		// Check RedirectConfig for custom redirect URI
		if s.RedirectURL != nil {
			if uri, ok := s.RedirectURL[service]; ok && uri != "" {
				redirectURL = fmt.Sprintf("%s%s?state=%s", s.webRedirectBase(providerInstance), uri, tokenHash)
			}
		}
	}

	// If no custom redirect configured, use default
	if redirectURL == "" {
		redirectURL = s.webRedirectBase(providerInstance) + constants.BindAccountBindURI + "?state=" + tokenHash
	}

	c.Redirect(http.StatusFound, redirectURL)
}

// GetWebUserByOauth gets user info from the cs-user login chain and builds the
// web user with a virtual Device record for account binding.
func GetWebUserByOauth(ctx context.Context, code, provider string) (*repository.AuthUser, bool, error) {
	oauthManager := providers.GetManager()
	providerInstance, err := oauthManager.GetProvider(provider)
	if err != nil {
		return nil, false, err
	}

	// Exchange code for token
	token, err := providerInstance.ExchangeToken(ctx, code)
	if err != nil {
		return nil, false, fmt.Errorf("failed to exchange token: %v", err)
	}

	// Get user info from the cs-user login chain
	user, isNewUser, err := fetchUserByIdentityChain(ctx, token.AccessToken)
	if err != nil {
		return nil, false, fmt.Errorf("%s: %w", errs.ErrInfoQueryUserInfo, err)
	}

	// Create virtual Device record for web users to enable account binding functionality
	var tokenProvider, refreshToken, accessToken string
	if provider == "casdoor" {
		refreshToken = token.RefreshToken
		accessToken = token.AccessToken
		tokenProvider = "custom" // Use token generated by provider
	}

	// Calculate token hashes for proper token management
	var refreshTokenHash, accessTokenHash string
	if refreshToken != "" {
		refreshTokenHash = utils.HashToken(refreshToken)
	}
	if accessToken != "" {
		accessTokenHash = utils.HashToken(accessToken)
	}

	// Create web device record with web-specific identifiers (following plugin field order)
	user.Devices = append(user.Devices, repository.Device{
		ID:               uuid.New(),
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now(),
		MachineCode:      "web-" + user.ID.String()[:8], // Generate unique web identifier
		UriScheme:        "https",                       // Web protocol
		VSCodeVersion:    "web-browser",                 // Fixed identifier for web platform
		PluginVersion:    "1.0.0",                       // Simplified version
		RefreshToken:     refreshToken,
		AccessToken:      accessToken,
		Provider:         provider,
		Platform:         "web",                          // Platform identifier
		Status:           constants.LoginStatusLoggedOut, // Initial status
		TokenProvider:    tokenProvider,                  // Token provider type
		RefreshTokenHash: refreshTokenHash,
		AccessTokenHash:  accessTokenHash,
		State:            "", // Will be set during callback if needed
		DeviceCode:       "",
	})

	// Note: User's own invite code will be generated when they first access the invite-code endpoint
	return user, isNewUser, nil
}

// getUserInviteCodeHandler gets current user's invite code
func (s *Server) getUserInviteCodeHandler(c *gin.Context) {
	// Get token from request header
	token, err := getTokenFromHeader(c)
	if err != nil {
		response.JSONError(c, http.StatusUnauthorized, errs.ErrAuthentication, "authentication failed: "+err.Error())
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	// Authenticate via the cs-user chain (same identity source as /manager/userinfo).
	// The previous local device-token lookup could not match after a token refresh
	// or a re-login replaced the stored device hashes (the web device MachineCode is
	// stable per user), leaving valid cs-user sessions unable to fetch their code.
	authUser, err := fetchUserInfoFromCsUser(ctx, token)
	if err != nil {
		status, codeMsg := mapUserCenterError(err)
		response.HandleError(c, status, codeMsg,
			fmt.Errorf("%s: %v", errs.ErrInfoQueryUserInfo, err))
		return
	}

	// Invite codes live in the local DB keyed by the universal_id.
	user, err := repository.GetDB().GetUserByField(ctx, "id", authUser.ID)
	if err != nil {
		response.JSONError(c, http.StatusInternalServerError, errs.ErrUpdateInfo, "failed to load user: "+err.Error())
		return
	}

	if user == nil {
		// No local record (the identity exists only in cs-user): seed a minimal
		// row so the generated invite code has a place to persist.
		user = authUser
		if user.UserCode == "" {
			user.UserCode, err = utils.GenerateRandomString(16)
			if err != nil {
				response.JSONError(c, http.StatusInternalServerError, errs.ErrUpdateInfo, "failed to create user: "+err.Error())
				return
			}
		}

		// Concurrent first-time fetches race on row creation (Upsert's create
		// path has no ON CONFLICT handling): if the create lost, another request
		// just seeded the row, so fall through to the atomic claim below.
		if err := repository.GetDB().Upsert(ctx, user, "id", user.ID); err != nil {
			if existing, getErr := repository.GetDB().GetUserByField(ctx, "id", user.ID); getErr != nil || existing == nil {
				response.JSONError(c, http.StatusInternalServerError, errs.ErrUpdateInfo, "failed to create user: "+err.Error())
				return
			}
		}
	}

	// Claim the invite code atomically when the row has none. Upsert alone is
	// not enough: its update branch can silently overwrite a code that another
	// request just committed, so the invite-code write is a single conditional
	// UPDATE (ClaimInviteCode) where the first writer wins.
	if user.InviteCode == "" {
		inviteCode, err := utils.GenerateUniqueInviteCode(ctx)
		if err != nil {
			response.JSONError(c, http.StatusInternalServerError, errs.ErrUpdateInfo, "failed to get invite code: "+err.Error())
			return
		}

		if err := repository.GetDB().ClaimInviteCode(ctx, user.ID, inviteCode); err != nil {
			response.JSONError(c, http.StatusInternalServerError, errs.ErrUpdateInfo, "failed to update invite code: "+err.Error())
			return
		}

		// Serve the authoritative code from the row: when a concurrent request
		// claimed it first, this adopts theirs instead of our (lost) one.
		final, err := repository.GetDB().GetUserByField(ctx, "id", user.ID)
		if err != nil || final == nil || final.InviteCode == "" {
			response.JSONError(c, http.StatusInternalServerError, errs.ErrUpdateInfo, "failed to load invite code: "+err.Error())
			return
		}
		user = final
	}

	// Return user invite code information
	response.JSONSuccess(c, "", gin.H{
		"invite_code": user.InviteCode,
	})
}
