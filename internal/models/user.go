package models

import (
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"time"
)

type UserRole string

const (
	RoleUser   UserRole = "user"
	RoleSeller UserRole = "seller"
	RoleAdmin  UserRole = "admin"
)

type User struct {
	UserID          uuid.UUID             `db:"user_id" json:"user_id"`
	Email           string                `db:"email" json:"email"`
	Phone           *string               `db:"phone" json:"phone"`
	PasswordHash    string                `db:"password_hash" json:"-"`
	AvatarURL       *string               `db:"avatar_url" json:"avatar_url"`
	Role            UserRole              `db:"role" json:"role"`
	MiningEnabled   bool                  `db:"mining_enabled" json:"mining_enabled"`
	MiningCPULimit  int                   `db:"mining_cpu_limit" json:"mining_cpu_limit"`
	SecuritySettings map[string]interface{} `db:"security_settings" json:"security_settings"`
	CreatedAt       time.Time             `db:"created_at" json:"created_at"`
	UpdatedAt       time.Time             `db:"updated_at" json:"updated_at"`
}

// HashPassword hashes a plaintext password using bcrypt
func (u *User) HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

// VerifyPassword checks if the provided password matches the hash
func (u *User) VerifyPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(password))
	return err == nil
}

// IsValidRole checks if the role is valid
func (u *User) IsValidRole() bool {
	switch u.Role {
	case RoleUser, RoleSeller, RoleAdmin:
		return true
	default:
		return false
	}
}

// CanEnableMining checks if user can enable mining
func (u *User) CanEnableMining() bool {
	return u.Role == RoleUser || u.Role == RoleSeller
}

// CanSetCPULimit validates CPU limit is within bounds
func (u *User) CanSetCPULimit(limit int) bool {
	return limit >= 1 && limit <= 10
}
