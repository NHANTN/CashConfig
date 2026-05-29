package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/cashier-config/server/internal/service"
)

type AuthHandler struct {
	Svc       *service.AuthService
	LDAPService *service.LDAPService
	SSOService  *service.SSOService
}

func NewAuthHandler(svc *service.AuthService, ldap *service.LDAPService, sso *service.SSOService) *AuthHandler {
	return &AuthHandler{Svc: svc, LDAPService: ldap, SSOService: sso}
}

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// @Summary      Login
// @Description  Authenticate with username and password
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        body  body  LoginRequest  true  "Login credentials"
// @Success      200   {object}  map[string]interface{}
// @Failure      401   {object}  map[string]interface{}
// @Router       /auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": err.Error()})
		return
	}
	token, user, err := 	h.Svc.Login(req.Username, req.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"code": 401, "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "ok",
		"data": gin.H{
			"token": token,
			"user":  user,
		},
	})
}

// @Summary      Logout
// @Description  Logout current session
// @Tags         Auth
// @Success      200  {object}  map[string]interface{}
// @Router       /auth/logout [post]
func (h *AuthHandler) Logout(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "logged out"})
}

// @Summary      Refresh token
// @Description  Refresh JWT token using current session
// @Tags         Auth
// @Security     ApiKeyAuth
// @Security     BearerAuth
// @Produce      json
// @Success      200  {object}  map[string]interface{}
// @Failure      500  {object}  map[string]interface{}
// @Router       /auth/refresh [post]
func (h *AuthHandler) Refresh(c *gin.Context) {
	claims := &service.Claims{
		UserID:   c.GetInt64("user_id"),
		Username: c.GetString("username"),
		RoleCode: c.GetString("role_code"),
	}
	token, err := 	h.Svc.GenerateTokenFromClaims(claims)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "ok",
		"data":    gin.H{"token": token},
	})
}

// @Summary      Get permissions
// @Description  Get current user's permissions
// @Tags         Auth
// @Security     ApiKeyAuth
// @Security     BearerAuth
// @Produce      json
// @Success      200  {object}  map[string]interface{}
// @Failure      500  {object}  map[string]interface{}
// @Router       /auth/permissions [get]
func (h *AuthHandler) GetPermissions(c *gin.Context) {
	userID, _ := c.Get("user_id")
	perms, err := 	h.Svc.GetPermissions(userID.(int64))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "ok", "data": perms})
}

type LDAPLoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// @Summary      LDAP Login
// @Description  Authenticate via LDAP
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        body  body  LDAPLoginRequest  true  "LDAP credentials"
// @Success      200   {object}  map[string]interface{}
// @Failure      400   {object}  map[string]interface{}
// @Failure      401   {object}  map[string]interface{}
// @Router       /auth/login/ldap [post]
func (h *AuthHandler) LDAPLogin(c *gin.Context) {
	if h.LDAPService == nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "LDAP not configured"})
		return
	}
	var req LDAPLoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": err.Error()})
		return
	}
	user, err := h.LDAPService.Authenticate(req.Username, req.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"code": 401, "message": err.Error()})
		return
	}
	token, err := h.Svc.GenerateTokenFromClaims(&service.Claims{
		UserID:   user.ID,
		Username: user.Username,
		RoleCode: user.Role.Code,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "ok",
		"data": gin.H{
			"token": token,
			"user":  user,
		},
	})
}

// @Summary      SSO Login URL
// @Description  Get SSO authentication URL
// @Tags         Auth
// @Produce      json
// @Success      200  {object}  map[string]interface{}
// @Failure      500  {object}  map[string]interface{}
// @Router       /auth/sso/login [get]
func (h *AuthHandler) SSOLogin(c *gin.Context) {
	if h.SSOService == nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "SSO not configured"})
		return
	}
	authURL, _, err := h.SSOService.AuthURL()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "ok", "data": gin.H{"auth_url": authURL}})
}

func (h *AuthHandler) Register(api gin.IRouter, authed gin.IRouter) {
	api.POST("/auth/login", h.Login)
	api.POST("/auth/login/ldap", h.LDAPLogin)
	api.GET("/auth/sso/login", h.SSOLogin)
	api.GET("/auth/sso/callback", h.SSOCallback)
	api.POST("/auth/logout", h.Logout)

	authed.POST("/auth/refresh", h.Refresh)
	authed.GET("/auth/permissions", h.GetPermissions)
}

// @Summary      SSO Callback
// @Description  Handle SSO callback with code and state
// @Tags         Auth
// @Produce      json
// @Param        code   query  string  true  "Authorization code"
// @Param        state  query  string  true  "State parameter"
// @Success      200  {object}  map[string]interface{}
// @Failure      400  {object}  map[string]interface{}
// @Failure      401  {object}  map[string]interface{}
// @Router       /auth/sso/callback [get]
func (h *AuthHandler) SSOCallback(c *gin.Context) {
	if h.SSOService == nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "SSO not configured"})
		return
	}
	code := c.Query("code")
	state := c.Query("state")
	if code == "" || state == "" {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "missing code or state"})
		return
	}
	user, err := h.SSOService.Callback(code, state)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"code": 401, "message": err.Error()})
		return
	}
	token, err := h.Svc.GenerateTokenFromClaims(&service.Claims{
		UserID:   user.ID,
		Username: user.Username,
		RoleCode: user.Role.Code,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "ok",
		"data": gin.H{
			"token": token,
			"user":  user,
		},
	})
}
