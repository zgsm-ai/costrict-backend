package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/zgsm-ai/oidc-auth/pkg/errs"
	"github.com/zgsm-ai/oidc-auth/pkg/log"

	"github.com/gin-gonic/gin"
	"github.com/zgsm-ai/oidc-auth/internal/constants"
	"github.com/zgsm-ai/oidc-auth/internal/providers"
	"github.com/zgsm-ai/oidc-auth/internal/repository"
	"github.com/zgsm-ai/oidc-auth/internal/service"
	github "github.com/zgsm-ai/oidc-auth/internal/sync"
	"github.com/zgsm-ai/oidc-auth/internal/usercenter"
	"github.com/zgsm-ai/oidc-auth/pkg/response"
	"github.com/zgsm-ai/oidc-auth/pkg/utils"
)

const (
	defaultTimeout = 45 * time.Second
	shortTimeout   = 10 * time.Second
	// identityTTL bounds how long the local identity columns are trusted before
	// a soft refresh from cs-user (decision 6). Must be >= cs-user's
	// GetOrCreateUser syncInterval (15m) so the two throttle layers never
	// invert.
	identityTTL = 24 * time.Hour
)

func getContextWithTimeout(timeout time.Duration) (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), timeout)
}

func getEncryptedData(data any) (string, error) {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return "", fmt.Errorf("failed to marshal data: %w", err)
	}

	keyManager, err := utils.GetEncryptKeyManager()
	if err != nil {
		return "", fmt.Errorf("failed to get key manager: %w", err)
	}

	return keyManager.AESEncrypt(jsonData)
}

func getDecryptedData(encryptedData string, result any) error {
	keyManager, err := utils.GetEncryptKeyManager()
	if err != nil {
		return fmt.Errorf("failed to get key manager: %w", err)
	}

	decrypted, err := keyManager.AESDecrypt(encryptedData)
	if err != nil {
		return fmt.Errorf("failed to decrypt data: %w", err)
	}

	return json.Unmarshal(decrypted, result)
}

func (s *Server) bindAccount(c *gin.Context) {
	token, err := getTokenFromHeader(c)
	if err != nil {
		response.HandleError(c, http.StatusBadRequest, errs.ErrBadRequestParam, err)
		return
	}

	tokenHash := utils.HashToken(token)
	oauthManager := providers.GetManager()
	provider := c.DefaultQuery("provider", "casdoor")
	providerInstance, err := oauthManager.GetProvider(provider)
	if err != nil {
		response.HandleError(c, http.StatusInternalServerError, errs.ErrBadRequestParam, err)
		return
	}

	if s.BaseURL == "" {
		response.HandleError(c, http.StatusInternalServerError, errs.ErrBadRequestParam,
			fmt.Errorf("base URL is not configured"))
		return
	}

	encryptedData, err := getEncryptedData(ParameterCarrier{
		TokenHash: tokenHash,
	})
	if err != nil {
		response.HandleError(c, http.StatusInternalServerError, errs.ErrDataEncryption, err)
		return
	}

	redirectURL := fmt.Sprintf("%s%s", s.BaseURL, constants.BindAccountCallbackURI)
	bindType := c.DefaultQuery("bindType", "")
	var bindParm string
	if bindType == "github" {
		bindParm = "&bindType=github"
	} else {
		bindParm = "&bindType=sms"
	}
	url := providerInstance.GetAuthURL(encryptedData, redirectURL) + bindParm

	response.JSONSuccess(c, "", map[string]interface{}{
		"state": c.DefaultQuery("state", ""),
		"url":   url,
	})
}

func (s *Server) bindAccountCallback(c *gin.Context) {
	code := c.DefaultQuery("code", "")
	if code == "" {
		response.HandleError(c, http.StatusBadRequest, errs.ErrBadRequestParam,
			errs.ParamNeedErr("code"))
		return
	}
	encryptedData := c.DefaultQuery("state", "")
	if encryptedData == "" {
		response.HandleError(c, http.StatusBadRequest, errs.ErrBadRequestParam,
			errs.ParamNeedErr("state"))
		return
	}
	var parameterCarrier ParameterCarrier
	if err := getDecryptedData(encryptedData, &parameterCarrier); err != nil {
		response.HandleError(c, http.StatusInternalServerError, errs.ErrDataDecryption, err)
		return
	}
	oauthManager := providers.GetManager()
	providerInstance, err := oauthManager.GetProvider("casdoor")
	if err != nil {
		response.HandleError(c, http.StatusInternalServerError, errs.ErrBadRequestParam, err)
		return
	}
	ctx, cancel := getContextWithTimeout(defaultTimeout)
	defer cancel()

	parameterCarrier.Provider = "casdoor"
	userNew, _, err := GetUserByOauth(ctx, "plugin", code, &parameterCarrier)
	if err != nil {
		response.HandleError(c, http.StatusInternalServerError, errs.ErrUserNotFound, err)
		return
	}
	userOld, err := repository.GetDB().GetUserByDeviceConditions(ctx, map[string]any{
		"access_token_hash": parameterCarrier.TokenHash,
	})
	if err != nil || userOld == nil || userNew == nil {
		response.HandleError(c, http.StatusUnauthorized, errs.ErrUserNotFound, errs.ErrInfoQueryUserInfo)
		return
	}
	// The already bound one cannot be bound again
	if userOld.GithubID != "" && userOld.Phone != "" {
		response.HandleError(c, http.StatusConflict, errs.ErrUpdateInfo, fmt.Errorf("this account has already been bound"))
		return
	}
	var useroldToken string
	for _, device := range userOld.Devices {
		if device.AccessTokenHash == parameterCarrier.TokenHash {
			useroldToken = device.AccessToken
			break
		}
	}
	// Get a new user and first determine whether it exists in the database
	var userNewExist *repository.AuthUser

	ctx, cancel = getContextWithTimeout(defaultTimeout)
	defer cancel()
	if userOld.GithubID != "" {
		userNewExist, err = repository.GetDB().GetUserByField(ctx, "phone", userNew.Phone)
	} else if userOld.Phone != "" {
		userNewExist, err = repository.GetDB().GetUserByField(ctx, "github_id", userNew.GithubID)
	} else {
		// custom types are not considered
		response.HandleError(c, http.StatusInternalServerError, errs.ErrTokenInvalid,
			fmt.Errorf("does not support custom account binding"))
		return
	}
	// Check for conflict: if existing account is already fully bound (has both GitHub and Phone)
	if userNewExist != nil {
		if userNewExist.GithubID != "" && userNewExist.Phone != "" {
			response.HandleError(c, http.StatusConflict, errs.ErrUpdateInfo, fmt.Errorf("this account has already been bound"))
			return
		}
	}

	// Determine main account (userMarge) and other account based on GitHub priority
	// Strategy: GitHub account always becomes the main account when both accounts exist
	userMarge, otherUser, mainToken, otherToken := determineMainAccount(userOld, userNew, userNewExist, useroldToken)

	// Merge fields from otherUser into userMarge
	userMarge.Email = coalesceString(userMarge.Email, otherUser.Email)
	userMarge.Phone = coalesceString(userMarge.Phone, otherUser.Phone)
	userMarge.GithubID = coalesceString(userMarge.GithubID, otherUser.GithubID)
	userMarge.GithubName = coalesceString(userMarge.GithubName, otherUser.GithubName)
	userMarge.UpdatedAt = time.Now()
	if userMarge.GithubName != "" {
		userMarge.Name = userMarge.GithubName
	} else {
		userMarge.Name = coalesceString(userMarge.Name, otherUser.Name)
	}
	userMarge.Company = coalesceString(userMarge.Company, otherUser.Company)
	userMarge.Location = coalesceString(userMarge.Location, otherUser.Location)
	userMarge.EmployeeNumber = coalesceString(userMarge.EmployeeNumber, otherUser.EmployeeNumber)
	userMarge.GithubStar = coalesceString(userMarge.GithubStar, otherUser.GithubStar)
	if otherUser.Vip > userMarge.Vip {
		userMarge.Vip = otherUser.Vip
	}
	userMarge.InviteCode = coalesceString(userMarge.InviteCode, otherUser.InviteCode)
	if userMarge.InviterID == nil || *userMarge.InviterID == uuid.Nil {
		if otherUser.InviterID != nil && *otherUser.InviterID != uuid.Nil {
			userMarge.InviterID = otherUser.InviterID
		}
	}

	// Call Casdoor merge API
	resp, err := service.MergeByCasdoor(providerInstance, mainToken, otherToken, s.HTTPClient)
	if err != nil {
		response.HandleError(c, http.StatusInternalServerError, errs.ErrBindAccount,
			fmt.Errorf("account linking failed, %w", err))
		return
	}
	if resp.Status != "ok" {
		log.Error(c, "failed to merge account. status: %s, msg: %s, UniversalID: %s, DeletedUserID: %s",
			resp.Status, resp.Msg, resp.UniversalID, resp.DeletedUserID)
		response.HandleError(c, http.StatusInternalServerError, errs.ErrBindAccount,
			fmt.Errorf("failed to merge account, status: %s, msg: %s", resp.Status, resp.Msg))
		return
	}

	// Handle existing account deletion and quota merge
	if userNewExist != nil {
		// Merge quota before deleting account
		err = service.MergeUserQuota(userMarge.ID.String(), otherUser.ID.String(), mainToken)
		if err != nil {
			response.HandleError(c, http.StatusInternalServerError, errs.ErrBindAccount,
				fmt.Errorf("failed to merge user quota: %w", err))
			return
		}

		// Delete the existing duplicate account
		if delNum, err := repository.GetDB().DeleteUserByField(ctx, constants.DBIndexField, otherUser.ID); err != nil || delNum == 0 {
			response.HandleError(c, http.StatusInternalServerError, errs.ErrBindAccount,
				fmt.Errorf("failed to delete old user, %w", err))
			return
		}
	}

	if err := repository.GetDB().Upsert(ctx, userMarge, constants.DBIndexField, userMarge.ID); err != nil {
		response.HandleError(c, http.StatusInternalServerError, errs.ErrUpdateInfo,
			fmt.Errorf("%s: %w", errs.ErrInfoUpdateUserInfo, err))
		return
	}

	// Use main account's token hash for redirect to ensure token validity
	tokenHash := getTokenHashForRedirect(userMarge, mainToken)

	url := s.webRedirectBase(providerInstance) + constants.BindAccountBindURI + "?state=" + tokenHash
	url = url + "&bind=true"
	c.Redirect(http.StatusFound, url)
}

func (s *Server) userInfoHandler(c *gin.Context) {
	token, err := getTokenFromHeader(c)
	if err != nil {
		response.HandleError(c, http.StatusBadRequest, errs.ErrBadRequestParam, err)
		return
	}

	tokenHash := utils.HashToken(token)
	ctx, cancel := getContextWithTimeout(shortTimeout)
	defer cancel()

	user, err := repository.GetDB().GetUserByDeviceConditions(ctx, map[string]any{
		"access_token_hash": tokenHash,
	})
	if err != nil || user == nil {
		response.HandleError(c, http.StatusBadRequest, errs.ErrTokenInvalid, errs.ErrInfoInvalidToken)
		return
	}

	// decision 6: identity fields are local columns with a soft TTL. Refresh
	// lazily when stale; refresh failure keeps stale values + warning, never
	// errors to the client.
	if user.IdentitySyncedAt == nil || time.Since(*user.IdentitySyncedAt) > identityTTL {
		user = refreshUserIdentity(ctx, user)
	}

	isStar := true
	starProject := user.GithubStar

	project := fmt.Sprintf("%s.%s", github.Owner, github.Repo)
	if starProject == "" || starProject != project {
		isStar = false
	}

	data := gin.H{
		"state":      c.DefaultQuery("state", ""),
		"username":   user.Name,
		"uuid":       user.ID.String(),
		"email":      user.Email,
		"phone":      user.Phone,
		"githubID":   user.GithubID,
		"githubName": user.GithubName,
		"isPrivate":  s.IsPrivate,
		"isStar":     isStar,
	}

	response.JSONSuccess(c, "", data)
}

// identityRefresh carries one in-flight cs-user identity refresh so concurrent
// requests for the same user coalesce into a single refresh (F-05).
type identityRefresh struct {
	done chan struct{}
	user *repository.AuthUser
	err  error
}

var (
	identityRefreshMu     sync.Mutex
	identityRefreshFlight = map[string]*identityRefresh{}
)

// refreshUserIdentity refreshes the user's identity fields from cs-user when
// identity_synced_at is stale (stale-while-revalidate, decision 6). It mutates
// user in place on success and persists the new values; on failure it leaves
// user untouched, logs a warning, and never errors (F-03).
func refreshUserIdentity(ctx context.Context, user *repository.AuthUser) *repository.AuthUser {
	if user == nil || user.SubjectID == "" {
		// Legacy user without a cs-user subject: nothing to refresh until the
		// next cs-user-backed login writes it.
		return user
	}
	key := user.ID.String()

	identityRefreshMu.Lock()
	if inFlight, ok := identityRefreshFlight[key]; ok {
		identityRefreshMu.Unlock()
		select {
		case <-inFlight.done:
			if inFlight.err == nil {
				return inFlight.user
			}
			return user
		case <-ctx.Done():
			return user
		}
	}
	flight := &identityRefresh{done: make(chan struct{})}
	identityRefreshFlight[key] = flight
	identityRefreshMu.Unlock()

	flight.user, flight.err = doRefreshIdentityFields(ctx, user)
	close(flight.done)

	identityRefreshMu.Lock()
	delete(identityRefreshFlight, key)
	identityRefreshMu.Unlock()

	if flight.err != nil {
		log.Warn(nil, "userinfo identity refresh failed for user %s (stale values served): %v", key, flight.err)
		return user
	}
	return flight.user
}

// doRefreshIdentityFields fetches the authoritative identity fields from
// cs-user (get-profile + auth-identities), applies them to the local columns,
// and persists via Upsert. Soft semantics: a field with no value in the
// refresh keeps the local value.
func doRefreshIdentityFields(ctx context.Context, user *repository.AuthUser) (*repository.AuthUser, error) {
	client := usercenter.GetClient()
	if client == nil {
		return user, fmt.Errorf("usercenter client not initialized")
	}
	profile, err := client.GetProfile(ctx, user.SubjectID)
	if err != nil {
		return user, err
	}
	identities, err := client.ListIdentities(ctx, user.SubjectID)
	if err != nil {
		return user, err
	}

	if profile.Username != "" {
		user.Name = profile.Username
	} else if profile.DisplayName != nil && *profile.DisplayName != "" {
		user.Name = *profile.DisplayName
	}
	if profile.Email != nil && *profile.Email != "" {
		user.Email = *profile.Email
	}
	if profile.Phone != nil && *profile.Phone != "" {
		user.Phone = *profile.Phone
	}
	for _, identity := range identities {
		if identity.Provider == "github" {
			user.GithubID = identity.ProviderUserID
			if identity.DisplayName != nil {
				user.GithubName = *identity.DisplayName
			}
			break
		}
	}
	now := time.Now()
	user.IdentitySyncedAt = &now
	if err := repository.GetDB().Upsert(ctx, user, "id", user.ID); err != nil {
		return user, fmt.Errorf("persist identity refresh: %w", err)
	}
	return user, nil
}

func coalesceString(values ...string) string {
	for _, v := range values {
		if v != "" {
			return v
		}
	}
	return ""
}

// determineMainAccount determines the main account and other account based on GitHub priority
// Strategy: GitHub account always becomes the main account when both accounts exist
// Returns: main account user, other account user, main token, other token
func determineMainAccount(userOld, userNew, userNewExist *repository.AuthUser, useroldToken string) (*repository.AuthUser, *repository.AuthUser, string, string) {
	var userMarge, otherUser *repository.AuthUser
	var mainToken, otherToken string

	if userNewExist == nil {
		// Scenario 1: Binding account doesn't exist - simple binding
		// Current logged-in user becomes main account, new OAuth info is supplementary
		userMarge = userOld
		otherUser = userNew
		mainToken = useroldToken
		otherToken = userNew.Devices[0].AccessToken
	} else {
		// Scenario 2: Binding account exists - GitHub account becomes main account
		if userNewExist.GithubID != "" {
			// Existing account has GitHub info, use it as main account
			userMarge = userNewExist
			otherUser = userOld
			mainToken = userNewExist.Devices[0].AccessToken
			otherToken = useroldToken
		} else {
			// Current account becomes main
			userMarge = userOld
			otherUser = userNewExist
			mainToken = useroldToken
			otherToken = userNewExist.Devices[0].AccessToken
		}
	}

	return userMarge, otherUser, mainToken, otherToken
}

// getTokenHashForRedirect gets the token hash for redirect to ensure token validity
// It searches for the device with the matching main token, falls back to first available token hash
func getTokenHashForRedirect(userMarge *repository.AuthUser, mainToken string) string {
	var tokenHash string
	for _, device := range userMarge.Devices {
		if device.AccessToken == mainToken {
			tokenHash = device.AccessTokenHash
			break
		}
	}
	// Fallback: use first available token hash if main token not found
	if tokenHash == "" && len(userMarge.Devices) > 0 {
		tokenHash = userMarge.Devices[0].AccessTokenHash
	}
	return tokenHash
}
