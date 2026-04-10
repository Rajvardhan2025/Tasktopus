package service

import (
	"context"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/yourusername/project-management/config"
	"github.com/yourusername/project-management/models"
	"github.com/yourusername/project-management/store"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrInvalidCredentials = errors.New("invalid email or password")
	ErrUserNotFound       = errors.New("user not found")
	ErrUserExists         = errors.New("user with this email already exists")
	ErrInvalidToken       = errors.New("invalid or expired token")
	ErrUserInactive       = errors.New("user account is inactive")
)

type AuthService struct {
	userStore *store.UserStore
	config    *config.Config
}

func NewAuthService(userStore *store.UserStore, config *config.Config) *AuthService {
	return &AuthService{
		userStore: userStore,
		config:    config,
	}
}

// Register creates a new user account
func (s *AuthService) Register(ctx context.Context, req *models.RegisterRequest) (*models.AuthResponse, error) {
	// Check if user already exists
	existingUser, err := s.userStore.FindByEmail(ctx, req.Email)
	if err == nil && existingUser != nil {
		return nil, ErrUserExists
	}

	// Hash password with bcrypt (cost 12 for security)
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), 12)
	if err != nil {
		return nil, err
	}

	// Create user
	user := &models.User{
		BaseModel:    models.BaseModel{ID: uuid.New().String()},
		Email:        req.Email,
		PasswordHash: string(hashedPassword),
		DisplayName:  req.DisplayName,
		IsActive:     true,
	}

	if err := s.userStore.Create(ctx, user); err != nil {
		return nil, err
	}

	// Generate JWT token
	token, err := s.generateToken(user)
	if err != nil {
		return nil, err
	}

	return &models.AuthResponse{
		Token: token,
		User:  user,
	}, nil
}

// Login authenticates a user and returns a JWT token
func (s *AuthService) Login(ctx context.Context, req *models.LoginRequest) (*models.AuthResponse, error) {
	// Find user by email
	user, err := s.userStore.FindByEmail(ctx, req.Email)
	if err != nil {
		return nil, ErrInvalidCredentials
	}

	// Check if user is active
	if !user.IsActive {
		return nil, ErrUserInactive
	}

	// Verify password
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		return nil, ErrInvalidCredentials
	}

	// Generate JWT token
	token, err := s.generateToken(user)
	if err != nil {
		return nil, err
	}

	return &models.AuthResponse{
		Token: token,
		User:  user,
	}, nil
}

// GetCurrentUser retrieves the authenticated user
func (s *AuthService) GetCurrentUser(ctx context.Context, userID string) (*models.User, error) {
	user, err := s.userStore.FindByID(ctx, userID)
	if err != nil {
		return nil, ErrUserNotFound
	}

	if !user.IsActive {
		return nil, ErrUserInactive
	}

	return user, nil
}

// ChangePassword updates user password
func (s *AuthService) ChangePassword(ctx context.Context, userID string, req *models.ChangePasswordRequest) error {
	user, err := s.userStore.FindByID(ctx, userID)
	if err != nil {
		return ErrUserNotFound
	}

	// Verify current password
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.CurrentPassword)); err != nil {
		return ErrInvalidCredentials
	}

	// Hash new password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), 12)
	if err != nil {
		return err
	}

	// Update password
	return s.userStore.UpdatePassword(ctx, userID, string(hashedPassword))
}

// ValidateToken validates a JWT token and returns the user ID
func (s *AuthService) ValidateToken(tokenString string) (string, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Validate signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrInvalidToken
		}
		return []byte(s.config.JWT.Secret), nil
	})

	if err != nil {
		return "", ErrInvalidToken
	}

	if !token.Valid {
		return "", ErrInvalidToken
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", ErrInvalidToken
	}

	// Check expiration
	exp, ok := claims["exp"].(float64)
	if !ok || time.Unix(int64(exp), 0).Before(time.Now()) {
		return "", ErrInvalidToken
	}

	userID, ok := claims["user_id"].(string)
	if !ok {
		return "", ErrInvalidToken
	}

	return userID, nil
}

// generateToken creates a JWT token for a user
func (s *AuthService) generateToken(user *models.User) (string, error) {
	now := time.Now()
	expiresAt := now.Add(s.config.JWT.Expiry)

	claims := jwt.MapClaims{
		"user_id": user.ID,
		"email":   user.Email,
		"iat":     now.Unix(),
		"exp":     expiresAt.Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.config.JWT.Secret))
}
