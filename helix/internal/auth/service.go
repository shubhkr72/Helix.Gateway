package auth

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/shubhkr72/helix/internal/jwt"
	"github.com/shubhkr72/helix/internal/password"
)

type Service struct {
	repo Repository
	jwt  *jwt.Manager
}

func NewService(
	repo Repository,
	jwtManager *jwt.Manager,
) *Service {

	return &Service{
		repo: repo,
		jwt:  jwtManager,
	}
}

func (s *Service) Register(
	ctx context.Context,
	req RegisterRequest,
) error {

	// Check whether email already exists.
	existingUser, err := s.repo.GetByEmail(ctx, req.Email)
	if err != nil {
		return err
	}

	if existingUser != nil {
		return ErrUserAlreadyExists
	}

	// Hash password.
	hash, err := password.HashPassword(req.Password)
	if err != nil {
		return err
	}

	now := time.Now()

	user := &User{
		ID:           uuid.New().String(),
		Email:        req.Email,
		PasswordHash: hash,
		Role:         "user",
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	return s.repo.CreateUser(ctx, user)
}

func (s *Service) Login(
	ctx context.Context,
	req LoginRequest,
) (string, error) {

	user, err := s.repo.GetByEmail(ctx, req.Email)
	if err != nil {
		return "", err
	}

	if user == nil {
		return "", ErrInvalidCredentials
	}

	err = password.VerifyPassword(
		req.Password,
		user.PasswordHash,
	)

	if err != nil {
		return "", ErrInvalidCredentials
	}

	token, err := s.jwt.GenerateToken(
		user.ID,
		user.Role,
		[]string{user.Role},
	)

	if err != nil {
		return "", err
	}

	return token, nil
}
