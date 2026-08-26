package handler

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/zgsm-ai/oidc-auth/internal/constants"
	"github.com/zgsm-ai/oidc-auth/internal/providers"
	"github.com/zgsm-ai/oidc-auth/internal/repository"
	"github.com/zgsm-ai/oidc-auth/pkg/errs"
	"github.com/zgsm-ai/oidc-auth/pkg/log"
	"github.com/zgsm-ai/oidc-auth/pkg/response"
	"github.com/zgsm-ai/oidc-auth/pkg/utils"
)

// tokenHandler handles token requests (return new refresh_token/access_token by refresh token)
func tokenHandler(c *gin.Context) {
	var query requestQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		response.JSONError(c, http.StatusBadRequest, errs.ErrBadRequestParam, err.Error())
		return
	}
	if query.MachineCode != "" && query.VscodeVersion == "" {
		response.JSONError(c, http.StatusBadRequest, errs.ErrBadRequestParam,
			errs.ParamNeedErr("machine_code or vscode_version").Error())
		return
	}
	// if MachineCode is provided, get the token for the first time
	// the account should have been pre-registered.
	if query.MachineCode != "" {
		if query.State == "" {
			response.JSONError(c, http.StatusBadRequest, errs.ErrBadRequestParam,
				errs.ParamNeedErr("state").Error())
			return
		}
		tokenPair, code, err := firstGetToken(query.MachineCode, query.VscodeVersion, query.State)
		if err != nil {
			response.JSONError(c, code, errs.ErrTokenGenerate, err.Error())
			return
		}
		if tokenPair == nil {
			response.JSONError(c, http.StatusInternalServerError, errs.ErrTokenGenerate,
				errs.ErrInfoGenerateToken.Error())
			return
		}
		response.JSONSuccess(c, "", gin.H{
			"access_token":  tokenPair.AccessToken,
			"refresh_token": tokenPair.RefreshToken,
			"state":         c.DefaultQuery("state", ""),
		})
		return
	}
	refreshToken, err := getTokenFromHeader(c)
	if err != nil {
		response.JSONError(c, http.StatusUnauthorized, errs.ErrAuthentication, err.Error())
		return
	}
	tokenPair, code, err := tokenRefresh(refreshToken)
	if err != nil {
		response.JSONError(c, code, errs.ErrTokenInvalid, err.Error())
		return
	}
	if tokenPair == nil {
		response.JSONError(c, http.StatusInternalServerError, errs.ErrTokenGenerate,
			errs.ErrInfoGenerateToken.Error())
		return
	}
	response.JSONSuccess(c, "", gin.H{
		"access_token":  tokenPair.AccessToken,
		"refresh_token": tokenPair.RefreshToken,
		"state":         c.DefaultQuery("state", ""),
	})
}

func firstGetToken(machineCode, vscodeVersion, state string) (*utils.TokenPair, int, error) {
	if vscodeVersion == "" {
		return nil, http.StatusUnauthorized, errs.ParamNeedErr("vscode_version")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	db := repository.GetDB()

	// Try to find a device that already has the final token (LoggedIn status)
	user, err := db.GetUserByDeviceConditions(ctx, map[string]any{
		"machine_code":      machineCode,
		"vscode_version":    vscodeVersion,
		"access_token_hash": state,
		"status":            constants.LoginStatusLoggedIn,
	})
	if err != nil {
		return nil, http.StatusUnauthorized, errs.ErrInfoQueryUserInfo
	}
	if user != nil {
		index := findDeviceIndex(user, machineCode, vscodeVersion)
		if index != -1 {
			log.Info(ctx, "firstGetToken: found existing LoggedIn device, returning token from DB")
			clientRefreshToken, err := mintClientRefreshToken(user, index, user.Devices[index].RefreshToken)
			if err != nil {
				return nil, http.StatusInternalServerError, err
			}
			return &utils.TokenPair{
				AccessToken:  user.Devices[index].AccessToken,
				RefreshToken: clientRefreshToken,
			}, http.StatusOK, nil
		}
	}

	// Try to find a device that has already been logged in via callback (state matches, status=LoggedIn)
	user, err = db.GetUserByDeviceConditions(ctx, map[string]any{
		"machine_code":   machineCode,
		"vscode_version": vscodeVersion,
		"state":          state,
		"status":         constants.LoginStatusLoggedIn,
	})
	if err != nil {
		return nil, http.StatusUnauthorized, errs.ErrInfoQueryUserInfo
	}
	if user != nil {
		index := findDeviceIndex(user, machineCode, vscodeVersion)
		if index != -1 {
			log.Info(ctx, "firstGetToken: found LoggedIn device by state, returning token from DB")
			clientRefreshToken, err := mintClientRefreshToken(user, index, user.Devices[index].RefreshToken)
			if err != nil {
				return nil, http.StatusInternalServerError, err
			}
			return &utils.TokenPair{
				AccessToken:  user.Devices[index].AccessToken,
				RefreshToken: clientRefreshToken,
			}, http.StatusOK, nil
		}
	}

	// Fallback: no final token yet, generate a new one
	user, err = db.GetUserByDeviceConditions(ctx, map[string]any{
		"machine_code":   machineCode,
		"vscode_version": vscodeVersion,
		"state":          state,
		"status":         constants.LoginStatusLoggedOut,
	})
	if err != nil {
		return nil, http.StatusUnauthorized, errs.ErrInfoQueryUserInfo
	}
	if user == nil {
		return &utils.TokenPair{
			AccessToken:  "",
			RefreshToken: "",
		}, http.StatusOK, nil
	}

	index := findDeviceIndex(user, machineCode, vscodeVersion)
	if index == -1 {
		return nil, http.StatusUnauthorized, errs.ErrInfoInvalidToken
	}

	clientPair, internalPair, err := generateTokenPair(ctx, user, index)
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}
	if clientPair == nil {
		return nil, http.StatusInternalServerError, errs.ErrInfoGenerateToken
	}

	user.Devices[index].Status = constants.LoginStatusLoggedIn
	if err := updateUserAndSave(ctx, user, index, internalPair); err != nil {
		return nil, http.StatusInternalServerError, err
	}

	return &utils.TokenPair{
		AccessToken:  clientPair.AccessToken,
		RefreshToken: clientPair.RefreshToken,
	}, http.StatusOK, nil
}

func tokenRefresh(refreshToken string) (*utils.TokenPair, int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Primary path: self-signed client refresh token (Plan A). It is stateless —
	// verified by signature, then mapped back to the device row via its claims,
	// so it survives re-login and coexists with other consumers' tokens.
	if claims, err := utils.VerifyRefreshToken(refreshToken); err == nil {
		userID, err := uuid.Parse(claims.Subject)
		if err != nil {
			return nil, http.StatusUnauthorized, errs.ErrInfoInvalidToken
		}
		user, err := repository.GetDB().GetUserByField(ctx, "id", userID)
		if err != nil || user == nil {
			return nil, http.StatusUnauthorized, errs.ErrInfoInvalidToken
		}
		index := findDeviceIndex(user, claims.MachineCode, claims.VsCodeVersion)
		if index == -1 || user.Devices[index].RefreshToken == "" {
			// Device gone or logged out (logout clears the stored refresh token).
			return nil, http.StatusUnauthorized, errs.ErrInfoInvalidToken
		}
		return refreshWithDevice(ctx, user, index)
	}

	// Legacy fallback: clients still holding the raw (pre-migration) refresh
	// token. Refreshing migrates them onto the self-signed scheme.
	user, index, err := utils.GetUserByTokenHash(ctx, refreshToken, "refresh_token_hash")
	if err != nil {
		return nil, http.StatusUnauthorized, err
	}
	if user == nil {
		return nil, http.StatusUnauthorized, errs.ErrInfoInvalidToken
	}
	return refreshWithDevice(ctx, user, index)
}

// refreshWithDevice rotates the stored provider refresh token, persists the new
// internal tokens in the device row, and returns the client-facing pair.
func refreshWithDevice(ctx context.Context, user *repository.AuthUser, index int) (*utils.TokenPair, int, error) {
	clientPair, internalPair, err := generateTokenPair(ctx, user, index)
	if err != nil {
		// A concurrent refresh already rotated the provider token; re-read the
		// device row (the winner persisted the fresh token) and retry once.
		if errors.Is(err, providers.ErrTokenRotated) {
			if fresh, freshIndex, rerr := reReadDevice(ctx, user.ID,
				user.Devices[index].MachineCode, user.Devices[index].VSCodeVersion); rerr == nil {
				user, index = fresh, freshIndex
				clientPair, internalPair, err = generateTokenPair(ctx, user, index)
			}
		}
		if err != nil {
			return nil, http.StatusInternalServerError, err
		}
	}

	if err := updateUserAndSave(ctx, user, index, internalPair); err != nil {
		return nil, http.StatusInternalServerError, err
	}

	return &utils.TokenPair{
		AccessToken:  clientPair.AccessToken,
		RefreshToken: clientPair.RefreshToken,
	}, http.StatusOK, nil
}

// reReadDevice reloads the user row by ID and resolves the device index again.
func reReadDevice(ctx context.Context, userID uuid.UUID, machineCode, vscodeVersion string) (*repository.AuthUser, int, error) {
	fresh, err := repository.GetDB().GetUserByField(ctx, "id", userID)
	if err != nil || fresh == nil {
		return nil, -1, errs.ErrInfoInvalidToken
	}
	index := findDeviceIndex(fresh, machineCode, vscodeVersion)
	if index == -1 {
		return nil, -1, errs.ErrInfoInvalidToken
	}
	return fresh, index, nil
}

// mintClientRefreshToken builds a fresh self-signed client refresh token whose
// expiry mirrors sourceToken's exp claim, defaulting to 30 days when that token
// has no decodable exp.
func mintClientRefreshToken(user *repository.AuthUser, index int, sourceToken string) (string, error) {
	exp := time.Now().Add(30 * 24 * time.Hour)
	if parsedExp, err := utils.ExtractTokenExp(sourceToken); err == nil {
		exp = parsedExp
	}
	return utils.GenerateRefreshToken(user, user.Devices[index], exp)
}

// generateTokenPair returns the pair served to the client and the pair
// persisted in the device row. For custom (casdoor) devices these differ: the
// client gets a self-signed refresh token while the internal pair holds the
// rotated casdoor access/refresh tokens.
func generateTokenPair(ctx context.Context, user *repository.AuthUser, index int) (*utils.TokenPair, *utils.TokenPair, error) {
	if user.Devices[index].TokenProvider == "custom" {
		return generateCustomTokenPair(ctx, user, index)
	}

	tokenPair, err := utils.GenerateTokenPairByUser(user, index)
	if err != nil || tokenPair == nil {
		return nil, nil, fmt.Errorf("%s, %v", errs.ErrInfoGenerateToken, err)
	}
	return tokenPair, tokenPair, nil
}

func findDeviceIndex(user *repository.AuthUser, machineCode, vscodeVersion string) int {
	if user == nil {
		return -1
	}
	for i, device := range user.Devices {
		if device.MachineCode == machineCode && device.VSCodeVersion == vscodeVersion {
			user.Devices[i].UpdatedAt = time.Now()
			return i
		}
	}
	return -1
}

func updateUserAndSave(ctx context.Context, user *repository.AuthUser, index int, tokenPair *utils.TokenPair) error {
	if user == nil || len(user.Devices) <= index {
		return errs.ErrInfoUpdateUserInfo
	}
	updateUserInfoMid(user, index, tokenPair)
	return repository.GetDB().Upsert(ctx, user, constants.DBIndexField, user.ID)
}

func updateUserInfoMid(user *repository.AuthUser, index int, tokenPair *utils.TokenPair) {
	accessTokenNew := tokenPair.AccessToken
	refreshTokenNew := tokenPair.RefreshToken
	refreshTokenHash := utils.HashToken(refreshTokenNew)
	accessTokenHash := utils.HashToken(accessTokenNew)
	user.UpdatedAt = time.Now()
	user.AccessTime = time.Now()
	user.Devices[index].UpdatedAt = time.Now()
	user.Devices[index].AccessToken = accessTokenNew
	user.Devices[index].RefreshToken = refreshTokenNew
	user.Devices[index].AccessTokenHash = accessTokenHash
	user.Devices[index].RefreshTokenHash = refreshTokenHash
}

func getTokenByHash(c *gin.Context) {
	accessTokenHash, err := getTokenFromHeader(c)
	if err != nil {
		response.JSONError(c, http.StatusUnauthorized, errs.ErrBadRequestParam,
			errs.ParamNeedErr("token").Error())
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	tokenPair, err := utils.GetTokenByTokenHash(ctx, accessTokenHash)
	if err != nil {
		response.JSONError(c, http.StatusUnauthorized, errs.ErrUserNotFound,
			fmt.Sprintf("%s, %s", errs.ErrInfoQueryUserInfo, err.Error()))
		return
	}
	response.JSONSuccess(c, "", gin.H{
		"state":        c.DefaultQuery("state", ""),
		"access_token": tokenPair.AccessToken,
	})
}

func getTokenFromHeader(c *gin.Context) (string, error) {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		return "", errs.ParamNeedErr("Authorization")
	}

	parts := strings.SplitN(authHeader, " ", 2)
	if len(parts) == 1 {
		return parts[0], nil
	}
	if !(len(parts) == 2 && parts[0] == "Bearer") {
		return "", errs.ParamNeedErr("Bearer")
	}

	tokenString := parts[1]
	return tokenString, nil
}

// generateCustomTokenPair rotates the stored provider (casdoor) refresh token
// and returns the client-facing pair (casdoor access + self-signed refresh) and
// the internal pair to persist (casdoor access + rotated casdoor refresh). The
// casdoor refresh token is never handed to clients, so multiple consumers share
// one device login without invalidating each other.
func generateCustomTokenPair(ctx context.Context, user *repository.AuthUser, index int) (*utils.TokenPair, *utils.TokenPair, error) {
	if user == nil {
		return nil, nil, fmt.Errorf("parameter user is nil")
	}
	if len(user.Devices) <= index {
		return nil, nil, fmt.Errorf("device not found")
	}
	storedRefreshToken := user.Devices[index].RefreshToken
	if storedRefreshToken == "" {
		return nil, nil, errs.ErrInfoInvalidToken
	}
	provider := user.Devices[index].Provider
	oauthManager := providers.GetManager()
	providerInstance, err := oauthManager.GetProvider(provider)
	if err != nil {
		return nil, nil, err
	}
	token, err := providerInstance.RefreshToken(ctx, storedRefreshToken)
	if err != nil {
		return nil, nil, err
	}

	// Mirror the rotated casdoor refresh token's expiry so clients observe the
	// same lifetime they would have with the raw token.
	clientRefreshToken, err := mintClientRefreshToken(user, index, token.RefreshToken)
	if err != nil {
		return nil, nil, err
	}

	return &utils.TokenPair{
			AccessToken:  token.AccessToken,
			RefreshToken: clientRefreshToken,
		}, &utils.TokenPair{
			AccessToken:  token.AccessToken,
			RefreshToken: token.RefreshToken,
		}, nil
}
