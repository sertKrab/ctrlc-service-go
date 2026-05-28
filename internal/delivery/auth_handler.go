package delivery

import (
	"github.com/gin-gonic/gin"
	authdto "git.trovefin.com/poc/ctrlc-service-go/internal/application/dto/auth"
	"git.trovefin.com/poc/ctrlc-service-go/internal/application/usecase"
	authusecase "git.trovefin.com/poc/ctrlc-service-go/internal/application/usecase/auth"
	"git.trovefin.com/poc/ctrlc-service-go/internal/config"
	"git.trovefin.com/poc/ctrlc-service-go/internal/infrastructure/middleware"
	"git.trovefin.com/poc/ctrlc-service-go/internal/response"
)

type AuthHandler struct {
	login   *authusecase.LoginUseCase
	refresh *authusecase.RefreshUseCase
	logout  *authusecase.LogoutUseCase
	getMe   *authusecase.GetMeUseCase
	cfg     *config.Config
}

func NewAuthHandler(
	login *authusecase.LoginUseCase,
	refresh *authusecase.RefreshUseCase,
	logout *authusecase.LogoutUseCase,
	getMe *authusecase.GetMeUseCase,
	cfg *config.Config,
) *AuthHandler {
	return &AuthHandler{login, refresh, logout, getMe, cfg}
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req authdto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "VALIDATION_REQUIRED", err.Error())
		return
	}
	result, err := h.login.Execute(c.Request.Context(), req, usecase.MetaFromGin(c))
	if err != nil {
		response.Unauthorized(c, "AUTH_INVALID_CREDENTIALS", err.Error())
		return
	}
	h.setRefreshCookie(c, result.RefreshToken)
	response.OK(c, gin.H{
		"user":         result.User,
		"access_token": result.AccessToken,
		"expires_in":   result.ExpiresIn,
	})
}

func (h *AuthHandler) Refresh(c *gin.Context) {
	rawToken, err := c.Cookie("refresh_token")
	if err != nil || rawToken == "" {
		response.Unauthorized(c, "AUTH_SESSION_REVOKED", "refresh token missing")
		return
	}
	result, err := h.refresh.Execute(c.Request.Context(), rawToken, usecase.MetaFromGin(c))
	if err != nil {
		h.clearRefreshCookie(c)
		response.Unauthorized(c, "AUTH_SESSION_REVOKED", err.Error())
		return
	}
	h.setRefreshCookie(c, result.RefreshToken)
	response.OK(c, gin.H{
		"access_token": result.AccessToken,
		"expires_in":   result.ExpiresIn,
	})
}

func (h *AuthHandler) Logout(c *gin.Context) {
	rawToken, _ := c.Cookie("refresh_token")
	userID := middleware.GetUserID(c)
	_ = h.logout.Execute(c.Request.Context(), rawToken, userID, usecase.MetaFromGin(c))
	h.clearRefreshCookie(c)
	response.OK(c, gin.H{"message": "logged out"})
}

func (h *AuthHandler) GetMe(c *gin.Context) {
	userID := middleware.GetUserID(c)
	result, err := h.getMe.Execute(c.Request.Context(), userID)
	if err != nil {
		response.NotFound(c, "USER_NOT_FOUND", err.Error())
		return
	}
	response.OK(c, result)
}

func (h *AuthHandler) setRefreshCookie(c *gin.Context, token string) {
	c.SetCookie(
		"refresh_token", token,
		int(h.cfg.RefreshTokenExpiry.Seconds()),
		"/api/v1/auth/refresh",
		h.cfg.CookieDomain,
		h.cfg.CookieSecureFlag(),
		true,
	)
}

func (h *AuthHandler) clearRefreshCookie(c *gin.Context) {
	c.SetCookie("refresh_token", "", -1, "/", h.cfg.CookieDomain, h.cfg.CookieSecureFlag(), true)
}
