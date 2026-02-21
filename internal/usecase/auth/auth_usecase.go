package auth

import (
	"context"
	"errors"
	"fmt"
	"time"

	"tms-core-service/internal/domain/entity"
	"tms-core-service/internal/domain/errs"
	"tms-core-service/internal/domain/repository"
	"tms-core-service/internal/domain/service"

	"github.com/google/uuid"
)

// AuthUseCase handles authentication operations
type AuthUseCase struct {
	userRepo       repository.UserRepository
	hashService    service.HashService
	tokenService   service.TokenService
	storageService service.StorageService
	accessExpiry   int64 // in minutes
	refreshExpiry  int64 // in hours
}

// NewAuthUseCase creates a new auth use case
func NewAuthUseCase(
	userRepo repository.UserRepository,
	hashService service.HashService,
	tokenService service.TokenService,
	storageService service.StorageService,
	accessExpiry, refreshExpiry int64,
) *AuthUseCase {
	return &AuthUseCase{
		userRepo:       userRepo,
		hashService:    hashService,
		tokenService:   tokenService,
		storageService: storageService,
		accessExpiry:   accessExpiry,
		refreshExpiry:  refreshExpiry,
	}
}

// Register registers a new user
func (uc *AuthUseCase) Register(ctx context.Context, input RegisterInput) (*AuthOutput, error) {
	// Check if user already exists (Email)
	existingUser, err := uc.userRepo.FindByEmail(ctx, input.Email)
	if err != nil && !errors.Is(err, errs.ErrNotFound) {
		return nil, fmt.Errorf("user repository: find by email: %w", err)
	}
	if existingUser != nil {
		return nil, errs.ErrConflict
	}

	// Check if user already exists (Phone Number)
	if input.PhoneNumber != "" {
		existingUser, err = uc.userRepo.FindByPhoneNumber(ctx, input.PhoneNumber)
		if err != nil && !errors.Is(err, errs.ErrNotFound) {
			return nil, fmt.Errorf("user repository: find by phone number: %w", err)
		}
		if existingUser != nil {
			return nil, errs.ErrConflict
		}
	}

	// Hash password via service
	passwordHash, err := uc.hashService.HashPassword(input.Password)
	if err != nil {
		return nil, fmt.Errorf("hash service: failed to hash password: %w", err)
	}

	// Create user
	user := &entity.User{
		Email:        stringPtr(input.Email),
		PhoneNumber:  stringPtr(input.PhoneNumber),
		PasswordHash: passwordHash,
		FirstName:    input.FirstName,
		LastName:     input.LastName,
	}

	if err := uc.userRepo.Create(ctx, user); err != nil {
		return nil, fmt.Errorf("user repository: create user: %w", err)
	}

	// Generate tokens
	return uc.generateTokens(user)
}

// Login authenticates a user
func (uc *AuthUseCase) Login(ctx context.Context, input LoginInput) (*AuthOutput, error) {
	// Find user by email
	user, err := uc.userRepo.FindByEmail(ctx, input.Email)
	if err != nil {
		if errors.Is(err, errs.ErrNotFound) {
			return nil, errs.ErrInvalidCredentials
		}
		return nil, fmt.Errorf("user repository: find by email: %w", err)
	}

	// Verify password via service
	if !uc.hashService.CheckPassword(input.Password, user.PasswordHash) {
		return nil, errs.ErrInvalidCredentials
	}

	// Generate tokens
	return uc.generateTokens(user)
}

// RefreshToken refreshes access and refresh tokens
func (uc *AuthUseCase) RefreshToken(ctx context.Context, refreshToken string) (*AuthOutput, error) {
	// Validate refresh token
	claims, err := uc.tokenService.ValidateToken(refreshToken)
	if err != nil {
		return nil, errs.ErrUnauthorized
	}

	// Find user by ID to ensure they still exist and are active
	user, err := uc.userRepo.FindByID(ctx, claims.UserID)
	if err != nil {
		if errors.Is(err, errs.ErrNotFound) {
			return nil, errs.ErrUnauthorized
		}
		return nil, fmt.Errorf("user repository: find by id: %w", err)
	}

	// Generate new tokens
	return uc.generateTokens(user)
}

// GetProfile returns user profile
func (uc *AuthUseCase) GetProfile(ctx context.Context, userID uuid.UUID) (*UserOutput, error) {
	user, err := uc.userRepo.FindByID(ctx, userID)
	if err != nil {
		if errors.Is(err, errs.ErrNotFound) {
			return nil, errs.ErrNotFound
		}
		return nil, fmt.Errorf("user repository: find by id: %w", err)
	}

	avatarURL := uc.resolveAvatarURL(ctx, user.AvatarURL)

	return &UserOutput{
		ID:          user.ID,
		Email:       user.Email,
		PhoneNumber: stringFromPtr(user.PhoneNumber),
		FirstName:   user.FirstName,
		LastName:    user.LastName,
		AvatarURL:   avatarURL,
	}, nil
}

// UpdateProfile updates a user's profile information
func (uc *AuthUseCase) UpdateProfile(ctx context.Context, userID uuid.UUID, input UpdateProfileInput) (*UserOutput, error) {
	user, err := uc.userRepo.FindByID(ctx, userID)
	if err != nil {
		if errors.Is(err, errs.ErrNotFound) {
			return nil, errs.ErrNotFound
		}
		return nil, fmt.Errorf("user repository: find by id: %w", err)
	}

	// Update allowed fields
	user.FirstName = input.FirstName
	user.LastName = input.LastName
	if input.PhoneNumber != "" {
		user.PhoneNumber = stringPtr(input.PhoneNumber)
	}
	if input.AvatarURL != "" {
		user.AvatarURL = input.AvatarURL
	}

	if err := uc.userRepo.Update(ctx, user); err != nil {
		return nil, fmt.Errorf("user repository: update user: %w", err)
	}

	return &UserOutput{
		ID:          user.ID,
		Email:       user.Email,
		PhoneNumber: stringFromPtr(user.PhoneNumber),
		FirstName:   user.FirstName,
		LastName:    user.LastName,
		AvatarURL:   uc.resolveAvatarURL(ctx, user.AvatarURL),
	}, nil
}

// GenerateAvatarUploadURL generates a presigned upload URL for a user's avatar
func (uc *AuthUseCase) GenerateAvatarUploadURL(ctx context.Context, userID uuid.UUID, contentType string) (*PresignUploadOutput, error) {
	// Determine file extension from content type
	ext := "jpg"
	switch contentType {
	case "image/png":
		ext = "png"
	case "image/webp":
		ext = "webp"
	case "image/gif":
		ext = "gif"
	}

	key := fmt.Sprintf("avatars/%s.%s", userID.String(), ext)

	uploadURL, err := uc.storageService.GenerateUploadURL(ctx, key, contentType)
	if err != nil {
		return nil, fmt.Errorf("storage service: generate upload url: %w", err)
	}

	return &PresignUploadOutput{
		UploadURL: uploadURL,
		ObjectKey: key,
	}, nil
}

// resolveAvatarURL returns a presigned GET URL if the stored value is an S3 key,
// otherwise returns the value as-is (could be empty or a full URL).
func (uc *AuthUseCase) resolveAvatarURL(ctx context.Context, stored string) string {
	if stored == "" {
		return ""
	}
	// If it already looks like a URL, return it as-is
	if len(stored) > 4 && stored[:4] == "http" {
		return stored
	}
	// Treat it as an S3 key and generate a presigned download URL
	if uc.storageService != nil {
		url, err := uc.storageService.GenerateDownloadURL(ctx, stored)
		if err == nil {
			return url
		}
	}
	return stored
}

// generateTokens generates access and refresh tokens
func (uc *AuthUseCase) generateTokens(user *entity.User) (*AuthOutput, error) {
	// Generate access token
	accessToken, err := uc.tokenService.GenerateToken(
		user.ID,
		stringFromPtr(user.Email),
		time.Duration(uc.accessExpiry)*time.Minute,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to generate access token: %w", err)
	}

	// Generate refresh token
	refreshToken, err := uc.tokenService.GenerateToken(
		user.ID,
		stringFromPtr(user.Email),
		time.Duration(uc.refreshExpiry)*time.Hour,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to generate refresh token: %w", err)
	}

	return &AuthOutput{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		User: &UserOutput{
			ID:          user.ID,
			Email:       user.Email,
			PhoneNumber: stringFromPtr(user.PhoneNumber),
			FirstName:   user.FirstName,
			LastName:    user.LastName,
			AvatarURL:   user.AvatarURL,
		},
	}, nil
}

func stringPtr(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}

func stringFromPtr(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}
