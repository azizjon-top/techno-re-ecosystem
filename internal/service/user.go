package service

import (
	"context"
	"errors"

	apierrors "github.com/azizjon-top/techno-re-ecosystem/internal/errors"
	"github.com/azizjon-top/techno-re-ecosystem/internal/models"
	"github.com/azizjon-top/techno-re-ecosystem/internal/repository"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

// RegisterRequest holds the input for user registration.
type RegisterRequest struct {
	Email    string
	Password string
	Role     models.UserRole
}

// LoginRequest holds the input for user login.
type LoginRequest struct {
	Email    string
	Password string
}

// AuthTokens holds access and refresh tokens.
type AuthTokens struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

// UserService handles user business logic.
type UserService struct {
	users  repository.UserRepository
	wallet repository.WalletRepository
	jwt    *JWTService
}

// NewUserService creates a UserService.
func NewUserService(users repository.UserRepository, wallet repository.WalletRepository, jwt *JWTService) *UserService {
	return &UserService{users: users, wallet: wallet, jwt: jwt}
}

// Register creates a new user account and a corresponding wallet.
func (s *UserService) Register(ctx context.Context, req RegisterRequest) (*models.User, *AuthTokens, error) {
	// Validate role
	role := req.Role
	if role == "" {
		role = models.RoleUser
	}

	user := &models.User{Role: role, Email: req.Email}
	if !user.IsValidRole() {
		return nil, nil, apierrors.NewAPIError(apierrors.ErrCodeBadRequest, "invalid role", 400)
	}

	hash, err := user.HashPassword(req.Password)
	if err != nil {
		return nil, nil, apierrors.NewAPIError(apierrors.ErrCodeInternalError, "failed to hash password", 500)
	}
	user.PasswordHash = hash

	if err := s.users.Create(ctx, user); err != nil {
		if errors.Is(err, repository.ErrConflict) {
			return nil, nil, apierrors.NewAPIError(apierrors.ErrCodeConflict, "email already in use", 409)
		}
		return nil, nil, apierrors.NewAPIError(apierrors.ErrCodeInternalError, "failed to create user", 500)
	}

	// Create wallet for the new user
	wallet := &models.Wallet{
		UserID:          user.UserID,
		BalanceTK:       decimal.Zero,
		FrozenBalanceTK: decimal.Zero,
	}
	if err := s.wallet.Create(ctx, wallet); err != nil {
		// non-fatal: log would go here
	}

	tokens, err := s.generateTokens(user)
	if err != nil {
		return nil, nil, err
	}
	return user, tokens, nil
}

// Login authenticates a user and returns tokens.
func (s *UserService) Login(ctx context.Context, req LoginRequest) (*models.User, *AuthTokens, error) {
	user, err := s.users.GetByEmail(ctx, req.Email)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, nil, apierrors.NewAPIError(apierrors.ErrCodeInvalidCredentials, "invalid email or password", 401)
		}
		return nil, nil, apierrors.NewAPIError(apierrors.ErrCodeInternalError, "failed to look up user", 500)
	}

	if !user.VerifyPassword(req.Password) {
		return nil, nil, apierrors.NewAPIError(apierrors.ErrCodeInvalidCredentials, "invalid email or password", 401)
	}

	tokens, err := s.generateTokens(user)
	if err != nil {
		return nil, nil, err
	}
	return user, tokens, nil
}

// GetByID returns a user by ID.
func (s *UserService) GetByID(ctx context.Context, id uuid.UUID) (*models.User, error) {
	user, err := s.users.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, apierrors.NewAPIError(apierrors.ErrCodeUserNotFound, "user not found", 404)
		}
		return nil, apierrors.NewAPIError(apierrors.ErrCodeInternalError, "failed to get user", 500)
	}
	return user, nil
}

func (s *UserService) generateTokens(user *models.User) (*AuthTokens, error) {
	access, err := s.jwt.GenerateAccessToken(user.UserID, string(user.Role))
	if err != nil {
		return nil, apierrors.NewAPIError(apierrors.ErrCodeInternalError, "failed to generate access token", 500)
	}
	refresh, err := s.jwt.GenerateRefreshToken(user.UserID, string(user.Role))
	if err != nil {
		return nil, apierrors.NewAPIError(apierrors.ErrCodeInternalError, "failed to generate refresh token", 500)
	}
	return &AuthTokens{AccessToken: access, RefreshToken: refresh}, nil
}
