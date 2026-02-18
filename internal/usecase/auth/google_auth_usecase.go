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

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	oauth2api "google.golang.org/api/oauth2/v2"
	"google.golang.org/api/option"
)

// GoogleAuthUseCase handles google authentication operations
type GoogleAuthUseCase struct {
	userRepo      repository.UserRepository
	tokenService  service.TokenService
	config        *oauth2.Config
	accessExpiry  int64 // in minutes
	refreshExpiry int64 // in hours
}

// NewGoogleAuthUseCase creates a new google auth use case
func NewGoogleAuthUseCase(
	userRepo repository.UserRepository,
	tokenService service.TokenService,
	clientID, clientSecret, redirectURL string,
	accessExpiry, refreshExpiry int64,
) *GoogleAuthUseCase {
	conf := &oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		RedirectURL:  redirectURL,
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.profile",
			"https://www.googleapis.com/auth/userinfo.email",
		},
		Endpoint: google.Endpoint,
	}

	return &GoogleAuthUseCase{
		userRepo:      userRepo,
		tokenService:  tokenService,
		config:        conf,
		accessExpiry:  accessExpiry,
		refreshExpiry: refreshExpiry,
	}
}

// GetGoogleLoginURL returns the Google OAuth login URL
func (uc *GoogleAuthUseCase) GetGoogleLoginURL(state string) string {
	return uc.config.AuthCodeURL(state)
}

// HandleGoogleCallback handles the Google OAuth callback
func (uc *GoogleAuthUseCase) HandleGoogleCallback(ctx context.Context, code string) (*AuthOutput, error) {
	// Exchange code for token
	token, err := uc.config.Exchange(ctx, code)
	if err != nil {
		return nil, fmt.Errorf("oauth2 exchange: %w", err)
	}

	// Get user info from Google
	oauth2Service, err := oauth2api.NewService(ctx, option.WithTokenSource(uc.config.TokenSource(ctx, token)))
	if err != nil {
		return nil, fmt.Errorf("oauth2 service: %w", err)
	}

	userInfo, err := oauth2Service.Userinfo.Get().Do()
	if err != nil {
		return nil, fmt.Errorf("get user info: %w", err)
	}

	// Find or create user
	user, err := uc.userRepo.FindByGoogleID(ctx, userInfo.Id)
	if err != nil && !errors.Is(err, errs.ErrNotFound) {
		return nil, fmt.Errorf("user repository: find by google id: %w", err)
	}

	if user == nil {
		// Try to find by email
		user, err = uc.userRepo.FindByEmail(ctx, userInfo.Email)
		if err != nil && !errors.Is(err, errs.ErrNotFound) {
			return nil, fmt.Errorf("user repository: find by email: %w", err)
		}

		if user != nil {
			// Link account
			user.GoogleID = stringPtr(userInfo.Id)
			if user.AvatarURL == "" {
				user.AvatarURL = userInfo.Picture
			}
			if err := uc.userRepo.Update(ctx, user); err != nil {
				return nil, fmt.Errorf("user repository: update user: %w", err)
			}
		} else {
			// Create new user
			user = &entity.User{
				Email:     stringPtr(userInfo.Email),
				FirstName: userInfo.GivenName,
				LastName:  userInfo.FamilyName,
				AvatarURL: userInfo.Picture,
				GoogleID:  stringPtr(userInfo.Id),
			}
			if err := uc.userRepo.Create(ctx, user); err != nil {
				return nil, fmt.Errorf("user repository: create user: %w", err)
			}
		}
	} else {
		// Update user info if needed
		updated := false
		if user.FirstName == "" && userInfo.GivenName != "" {
			user.FirstName = userInfo.GivenName
			updated = true
		}
		if user.LastName == "" && userInfo.FamilyName != "" {
			user.LastName = userInfo.FamilyName
			updated = true
		}
		if user.AvatarURL != userInfo.Picture {
			user.AvatarURL = userInfo.Picture
			updated = true
		}

		if updated {
			if err := uc.userRepo.Update(ctx, user); err != nil {
				return nil, fmt.Errorf("user repository: update user: %w", err)
			}
		}
	}

	// Generate tokens
	return uc.generateTokens(user)
}

// generateTokens generates access and refresh tokens (duplicated from AuthUseCase for simplicity, should ideally be shared)
func (uc *GoogleAuthUseCase) generateTokens(user *entity.User) (*AuthOutput, error) {
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
