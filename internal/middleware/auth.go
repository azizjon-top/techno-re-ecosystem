package middleware

import (
	"strings"

	"github.com/azizjon-top/techno-re-ecosystem/internal/service"
	"github.com/gin-gonic/gin"
)

const (
	// ContextKeyUserID is the gin context key for the authenticated user's ID.
	ContextKeyUserID = "user_id"
	// ContextKeyUserRole is the gin context key for the authenticated user's role.
	ContextKeyUserRole = "user_role"
)

// Auth returns a Gin middleware that validates JWT tokens.
func Auth(jwtSvc *service.JWTService) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(401, gin.H{
				"code":    "UNAUTHORIZED",
				"message": "Authorization header is required",
			})
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
			c.AbortWithStatusJSON(401, gin.H{
				"code":    "UNAUTHORIZED",
				"message": "Invalid authorization header format",
			})
			return
		}

		claims, err := jwtSvc.ValidateToken(parts[1])
		if err != nil {
			c.AbortWithStatusJSON(401, gin.H{
				"code":    "UNAUTHORIZED",
				"message": "Invalid or expired token",
			})
			return
		}

		c.Set(ContextKeyUserID, claims.UserID)
		c.Set(ContextKeyUserRole, claims.Role)
		c.Next()
	}
}

// RequireRole returns a middleware that enforces a minimum role.
func RequireRole(allowedRoles ...string) gin.HandlerFunc {
	allowed := make(map[string]bool, len(allowedRoles))
	for _, r := range allowedRoles {
		allowed[r] = true
	}
	return func(c *gin.Context) {
		role, exists := c.Get(ContextKeyUserRole)
		if !exists || !allowed[role.(string)] {
			c.AbortWithStatusJSON(403, gin.H{
				"code":    "FORBIDDEN",
				"message": "You do not have permission to access this resource",
			})
			return
		}
		c.Next()
	}
}
