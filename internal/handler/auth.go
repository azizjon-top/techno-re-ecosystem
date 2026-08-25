package handler

import (
	"net/http"

	apierrors "github.com/azizjon-top/techno-re-ecosystem/internal/errors"
	"github.com/azizjon-top/techno-re-ecosystem/internal/middleware"
	"github.com/azizjon-top/techno-re-ecosystem/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// AuthHandler handles authentication endpoints.
type AuthHandler struct {
	userSvc *service.UserService
}

// NewAuthHandler creates an AuthHandler.
func NewAuthHandler(userSvc *service.UserService) *AuthHandler {
	return &AuthHandler{userSvc: userSvc}
}

// registerRequest is the JSON body for POST /auth/register.
type registerRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
}

// Register godoc
// POST /api/v1/auth/register
func (h *AuthHandler) Register(c *gin.Context) {
	var req registerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    apierrors.ErrCodeValidationFailed,
			"message": err.Error(),
		})
		return
	}

	user, tokens, err := h.userSvc.Register(c.Request.Context(), service.RegisterRequest{
		Email:    req.Email,
		Password: req.Password,
	})
	if err != nil {
		respondError(c, err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"user":   user,
		"tokens": tokens,
	})
}

// loginRequest is the JSON body for POST /auth/login.
type loginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// Login godoc
// POST /api/v1/auth/login
func (h *AuthHandler) Login(c *gin.Context) {
	var req loginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    apierrors.ErrCodeValidationFailed,
			"message": err.Error(),
		})
		return
	}

	user, tokens, err := h.userSvc.Login(c.Request.Context(), service.LoginRequest{
		Email:    req.Email,
		Password: req.Password,
	})
	if err != nil {
		respondError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"user":   user,
		"tokens": tokens,
	})
}

// UserHandler handles user endpoints.
type UserHandler struct {
	userSvc *service.UserService
}

// NewUserHandler creates a UserHandler.
func NewUserHandler(userSvc *service.UserService) *UserHandler {
	return &UserHandler{userSvc: userSvc}
}

// GetUser godoc
// GET /api/v1/users/:id
func (h *UserHandler) GetUser(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    apierrors.ErrCodeBadRequest,
			"message": "invalid user id",
		})
		return
	}

	// Users can only fetch their own profile unless they're an admin
	callerID := c.MustGet(middleware.ContextKeyUserID).(uuid.UUID)
	callerRole := c.MustGet(middleware.ContextKeyUserRole).(string)
	if callerID != id && callerRole != "admin" {
		c.JSON(http.StatusForbidden, gin.H{
			"code":    apierrors.ErrCodeForbidden,
			"message": "access denied",
		})
		return
	}

	user, svcErr := h.userSvc.GetByID(c.Request.Context(), id)
	if svcErr != nil {
		respondError(c, svcErr)
		return
	}

	c.JSON(http.StatusOK, user)
}

// GetMe godoc
// GET /api/v1/users/me
func (h *UserHandler) GetMe(c *gin.Context) {
	callerID := c.MustGet(middleware.ContextKeyUserID).(uuid.UUID)
	user, err := h.userSvc.GetByID(c.Request.Context(), callerID)
	if err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusOK, user)
}
