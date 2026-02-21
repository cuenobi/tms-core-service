package dto

// UserResponse represents user information in responses
type UserResponse struct {
	ID          string  `json:"id"`
	Email       *string `json:"email"`
	FirstName   string  `json:"first_name"`
	LastName    string  `json:"last_name"`
	PhoneNumber string  `json:"phone_number"`
	AvatarURL   string  `json:"avatar_url"`
}

// RegisterRequest represents a user registration request
type RegisterRequest struct {
	Email       string `json:"email" validate:"required,email"`
	Password    string `json:"password" validate:"required,min=6"`
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	PhoneNumber string `json:"phone_number" validate:"required,e164"`
}

// LoginRequest represents a user login request
type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

// AuthResponse represents an authentication response
type AuthResponse struct {
	AccessToken  string        `json:"access_token"`
	RefreshToken string        `json:"refresh_token"`
	User         *UserResponse `json:"user"`
}

// RefreshTokenRequest represents a token refresh request
type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}

// UpdateProfileRequest represents a user profile update request
type UpdateProfileRequest struct {
	FirstName   string `json:"first_name" validate:"required"`
	LastName    string `json:"last_name" validate:"required"`
	PhoneNumber string `json:"phone_number" validate:"omitempty,e164"`
	AvatarURL   string `json:"avatar_url" validate:"omitempty"`
}

// AvatarUploadRequest represents a request for a presigned avatar upload URL
type AvatarUploadRequest struct {
	ContentType string `json:"content_type" validate:"required,oneof=image/jpeg image/png image/webp image/gif"`
}

// PresignUploadResponse represents the presigned upload URL response
type PresignUploadResponse struct {
	UploadURL string `json:"upload_url"`
	ObjectKey string `json:"object_key"`
}
