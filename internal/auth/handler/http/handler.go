package http

import (
	"encoding/json"
	"io"
	"mime/multipart"
	"net/http"
	"path/filepath"
	"time"

	"github.com/ADRPUR/event-driven-marketplace/internal/auth/model"
	"github.com/ADRPUR/event-driven-marketplace/internal/auth/service"
	"github.com/ADRPUR/event-driven-marketplace/internal/auth/utils"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	svc service.AuthService
}

func New(svc service.AuthService) *Handler { return &Handler{svc: svc} }

// RegisterPublicRoutes mounts /register, /login, /refresh
func RegisterPublicRoutes(r *gin.Engine, h *Handler) {
	r.POST("/register", h.register)
	r.POST("/login", h.login)
	r.POST("/refresh", h.refresh)
}

// RegisterProtectedRoutes mounts endpoints that require authentication (middleware)
func RegisterProtectedRoutes(r *gin.RouterGroup, h *Handler) {
	r.POST("/logout", h.logout)
	r.GET("/me", h.me)
	r.PUT("/me", h.updateMe)
	r.POST("/me/photo", h.uploadPhoto)
	r.POST("/me/password", h.changePassword)
}

// -------------------- Handlers --------------------

func (h *Handler) register(c *gin.Context) {
	var req struct {
		Email       string         `json:"email" binding:"required,email"`
		Password    string         `json:"password" binding:"required,min=6"`
		Role        string         `json:"role"`
		FirstName   string         `json:"firstName"`
		LastName    string         `json:"lastName"`
		DateOfBirth string         `json:"dateOfBirth"`
		Phone       string         `json:"phone"`
		Address     map[string]any `json:"address"` // you can unmarshal into a struct if you want
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	dob, _ := time.Parse("2006-01-02", req.DateOfBirth)
	u := &model.User{
		Email: req.Email, Role: req.Role,
	}
	d := &model.UserDetails{
		FirstName:   req.FirstName,
		LastName:    req.LastName,
		DateOfBirth: dob,
		Phone:       req.Phone,
	}
	if req.Address != nil {
		bytes, err := json.Marshal(req.Address)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid address"})
			return
		}
		d.Address = bytes
	}
	err := h.svc.Register(c, u, d, req.Password)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"id": u.ID})
}

func (h *Handler) login(c *gin.Context) {
	var req struct {
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	at, rt, st, pl, err := h.svc.Login(c, req.Email, req.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}
	user, details, _ := h.svc.GetUserWithDetails(c, pl.UserID)
	c.JSON(http.StatusOK, gin.H{
		"accessToken":  at,
		"refreshToken": rt,
		"sessionToken": st,
		"expiresAt":    pl.ExpiredAt.Unix(),
		"user": gin.H{
			"id":          user.ID,
			"email":       user.Email,
			"role":        user.Role,
			"firstName":   details.FirstName,
			"lastName":    details.LastName,
			"dateOfBirth": details.DateOfBirth,
			"phone":       details.Phone,
			"address":     details.Address,
			"photo":       details.PhotoPath,
			"thumbnail":   details.ThumbnailPath,
		},
	})
}

func (h *Handler) refresh(c *gin.Context) {
	var req struct {
		SessionToken string `json:"sessionToken" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	at, pl, err := h.svc.Refresh(c, req.SessionToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"accessToken": at,
		"expiresAt":   pl.ExpiredAt.Unix(),
	})
}

func (h *Handler) logout(c *gin.Context) {
	var req struct {
		SessionToken string `json:"sessionToken" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.svc.Logout(c, req.SessionToken); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}

func (h *Handler) me(c *gin.Context) {
	payload := utils.ExtractPayload(c)
	if payload == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "not authenticated"})
		return
	}
	user, details, err := h.svc.GetUserWithDetails(c, payload.UserID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}

	// Address conversion to map[string]any
	var address map[string]any
	if details != nil && len(details.Address) > 0 {
		if err := json.Unmarshal(details.Address, &address); err != nil {
			address = nil // or you can log/ignore
		}
	}
	c.JSON(http.StatusOK, gin.H{
		"id":      user.ID,
		"email":   user.Email,
		"role":    user.Role,
		"details": details,
	})

	//c.JSON(http.StatusOK, gin.H{
	//	"id":          user.ID,
	//	"email":       user.Email,
	//	"role":        user.Role,
	//	"firstName":   safeStr(details, func(d *model.UserDetails) string { return d.FirstName }),
	//	"lastName":    safeStr(details, func(d *model.UserDetails) string { return d.LastName }),
	//	"dateOfBirth": safeDate(details, func(d *model.UserDetails) time.Time { return d.DateOfBirth }),
	//	"phone":       safeStr(details, func(d *model.UserDetails) string { return d.Phone }),
	//	"address":     address,
	//	"photo":       safeStr(details, func(d *model.UserDetails) string { return d.PhotoPath }),
	//	"thumbnail":   safeStr(details, func(d *model.UserDetails) string { return d.ThumbnailPath }),
	//})
}

// [PUT] /me — update user details
func (h *Handler) updateMe(c *gin.Context) {
	payload := utils.ExtractPayload(c)
	if payload == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "not authenticated"})
		return
	}
	var req struct {
		FirstName   string         `json:"firstName"`
		LastName    string         `json:"lastName"`
		DateOfBirth string         `json:"dateOfBirth"`
		Phone       string         `json:"phone"`
		Address     map[string]any `json:"address"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	dob, _ := time.Parse("2006-01-02", req.DateOfBirth)
	_, details, err := h.svc.GetUserWithDetails(c, payload.UserID)
	if err != nil || details == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "user details not found"})
		return
	}
	details.FirstName = req.FirstName
	details.LastName = req.LastName
	details.DateOfBirth = dob
	details.Phone = req.Phone
	if req.Address != nil {
		bytes, err := json.Marshal(req.Address)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid address"})
			return
		}
		details.Address = bytes
	}
	if err := h.svc.UpdateUserDetails(c, payload.UserID, details); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusOK)
}

// POST /me/photo — upload profile photo
func (h *Handler) uploadPhoto(c *gin.Context) {
	payload := utils.ExtractPayload(c)
	if payload == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "not authenticated"})
		return
	}
	file, err := c.FormFile("photo")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "file required"})
		return
	}
	src, err := file.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer func(src multipart.File) {
		err := src.Close()
		if err != nil {

		}
	}(src)
	data, err := io.ReadAll(src)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ext := filepath.Ext(file.Filename)
	photoPath, thumbPath, err := h.svc.UploadPhoto(c, payload.UserID, data, ext)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"photo": photoPath, "thumbnail": thumbPath})
}

// POST /me/password — change password
func (h *Handler) changePassword(c *gin.Context) {
	payload := utils.ExtractPayload(c)
	if payload == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "not authenticated"})
		return
	}
	var req struct {
		OldPassword string `json:"oldPassword" binding:"required"`
		NewPassword string `json:"newPassword" binding:"required,min=6"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.svc.ChangePassword(c, payload.UserID, req.OldPassword, req.NewPassword); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusOK)
}

func safeStr(d *model.UserDetails, f func(*model.UserDetails) string) string {
	if d != nil {
		return f(d)
	}
	return ""
}

func safeDate(d *model.UserDetails, f func(*model.UserDetails) time.Time) time.Time {
	if d != nil {
		return f(d)
	}
	return time.Time{}
}
