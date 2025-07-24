package http

import (
	"net/http"

	"github.com/ADRPUR/event-driven-marketplace/internal/auth/service"
	"github.com/gin-gonic/gin"
)

// Handler exposes HTTP endpoints for auth.
type Handler struct {
	svc *service.Service
}

func New(svc *service.Service) *Handler { return &Handler{svc: svc} }

// RegisterPublicRoutes mounts login & refresh.
func RegisterPublicRoutes(r *gin.Engine, h *Handler) {
	grp := r.Group("/auth")
	grp.POST("/login", h.login)
	grp.POST("/refresh", h.refresh)
}

// RegisterProtectedRoutes mounts logout under /auth.
func RegisterProtectedRoutes(r *gin.RouterGroup, h *Handler) {
	r.POST("/logout", h.logout)
}

// DTOs
type loginReq struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

type loginResp struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
	ExpiresAt    int64  `json:"expiresAt"`
}

func (h *Handler) login(c *gin.Context) {
	var req loginReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	at, rt, pl, err := h.svc.Login(c, req.Email, req.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, loginResp{AccessToken: at, RefreshToken: rt, ExpiresAt: pl.ExpiredAt.Unix()})
}

func (h *Handler) refresh(c *gin.Context) {
	var body struct {
		RefreshToken string `json:"refreshToken" binding:"required"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	at, pl, err := h.svc.Refresh(c, body.RefreshToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"accessToken": at, "expiresAt": pl.ExpiredAt.Unix()})
}

func (h *Handler) logout(c *gin.Context) {
	var body struct {
		RefreshToken string `json:"refreshToken" binding:"required"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.svc.Logout(c, body.RefreshToken); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}
