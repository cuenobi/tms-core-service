package auth

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"tms-core-service/internal/domain/entity"
	"tms-core-service/internal/domain/errs"
	"tms-core-service/internal/domain/repository"
	"tms-core-service/internal/domain/service"
)

const (
	lineAuthURL    = "https://access.line.me/oauth2/v2.1/authorize"
	lineTokenURL   = "https://api.line.me/oauth2/v2.1/token"
	lineProfileURL = "https://api.line.me/v2/profile"
)

// lineTokenResponse represents LINE token exchange response
type lineTokenResponse struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int    `json:"expires_in"`
	Scope        string `json:"scope"`
	IDToken      string `json:"id_token"`
}

// lineProfile represents LINE user profile
type lineProfile struct {
	UserID        string `json:"userId"`
	DisplayName   string `json:"displayName"`
	PictureURL    string `json:"pictureUrl"`
	StatusMessage string `json:"statusMessage"`
}

// LineAuthUseCase handles LINE authentication operations
type LineAuthUseCase struct {
	userRepo      repository.UserRepository
	tokenService  service.TokenService
	channelID     string
	channelSecret string
	redirectURL   string
	accessExpiry  int64 // in minutes
	refreshExpiry int64 // in hours
}

// NewLineAuthUseCase creates a new LINE auth use case
func NewLineAuthUseCase(
	userRepo repository.UserRepository,
	tokenService service.TokenService,
	channelID, channelSecret, redirectURL string,
	accessExpiry, refreshExpiry int64,
) *LineAuthUseCase {
	return &LineAuthUseCase{
		userRepo:      userRepo,
		tokenService:  tokenService,
		channelID:     channelID,
		channelSecret: channelSecret,
		redirectURL:   redirectURL,
		accessExpiry:  accessExpiry,
		refreshExpiry: refreshExpiry,
	}
}

// GetLineLoginURL returns the LINE OAuth login URL
func (uc *LineAuthUseCase) GetLineLoginURL(state string) string {
	params := url.Values{}
	params.Set("response_type", "code")
	params.Set("client_id", uc.channelID)
	params.Set("redirect_uri", uc.redirectURL)
	params.Set("state", state)
	params.Set("scope", "profile openid email")

	return lineAuthURL + "?" + params.Encode()
}

// HandleLineCallback handles the LINE OAuth callback
func (uc *LineAuthUseCase) HandleLineCallback(ctx context.Context, code string) (*AuthOutput, error) {
	// Exchange code for token
	lineToken, err := uc.exchangeCode(ctx, code)
	if err != nil {
		return nil, fmt.Errorf("line token exchange: %w", err)
	}

	// Get user profile from LINE
	profile, err := uc.getProfile(ctx, lineToken.AccessToken)
	if err != nil {
		return nil, fmt.Errorf("line get profile: %w", err)
	}

	// Find or create user
	user, err := uc.userRepo.FindByLineID(ctx, profile.UserID)
	if err != nil && !errors.Is(err, errs.ErrNotFound) {
		return nil, fmt.Errorf("user repository: find by line id: %w", err)
	}

	if user == nil {
		// Create new user (LINE does not provide email via profile API by default)
		user = &entity.User{
			FirstName: profile.DisplayName,
			AvatarURL: profile.PictureURL,
			Email:     nil, // Explicitly set to nil to ensure NULL in DB
			LineID:    stringPtr(profile.UserID),
		}
		if err := uc.userRepo.Create(ctx, user); err != nil {
			return nil, fmt.Errorf("user repository: create user: %w", err)
		}
	} else {
		// Update user info if needed
		updated := false
		if user.FirstName == "" && profile.DisplayName != "" {
			user.FirstName = profile.DisplayName
			updated = true
		}
		// Only update avatar if currently empty
		if user.AvatarURL == "" && profile.PictureURL != "" {
			user.AvatarURL = profile.PictureURL
			updated = true
		}

		if updated {
			if err := uc.userRepo.Update(ctx, user); err != nil {
				return nil, fmt.Errorf("user repository: update user: %w", err)
			}
		}
	}

	// Generate JWT tokens
	return uc.generateTokens(user)
}

// exchangeCode exchanges authorization code for LINE access token
func (uc *LineAuthUseCase) exchangeCode(_ context.Context, code string) (*lineTokenResponse, error) {
	data := url.Values{}
	data.Set("grant_type", "authorization_code")
	data.Set("code", code)
	data.Set("redirect_uri", uc.redirectURL)
	data.Set("client_id", uc.channelID)
	data.Set("client_secret", uc.channelSecret)

	resp, err := http.Post(lineTokenURL, "application/x-www-form-urlencoded", strings.NewReader(data.Encode()))
	if err != nil {
		return nil, fmt.Errorf("http post: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("line token exchange failed: status=%d body=%s", resp.StatusCode, string(body))
	}

	var tokenResp lineTokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&tokenResp); err != nil {
		return nil, fmt.Errorf("decode token response: %w", err)
	}

	return &tokenResp, nil
}

// getProfile fetches the LINE user profile
func (uc *LineAuthUseCase) getProfile(_ context.Context, accessToken string) (*lineProfile, error) {
	req, err := http.NewRequest(http.MethodGet, lineProfileURL, nil)
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+accessToken)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("http get: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("line get profile failed: status=%d body=%s", resp.StatusCode, string(body))
	}

	var profile lineProfile
	if err := json.NewDecoder(resp.Body).Decode(&profile); err != nil {
		return nil, fmt.Errorf("decode profile: %w", err)
	}

	return &profile, nil
}

// generateTokens generates access and refresh tokens
func (uc *LineAuthUseCase) generateTokens(user *entity.User) (*AuthOutput, error) {
	accessToken, err := uc.tokenService.GenerateToken(
		user.ID,
		stringFromPtr(user.Email),
		time.Duration(uc.accessExpiry)*time.Minute,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to generate access token: %w", err)
	}

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
