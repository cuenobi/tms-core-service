package auth

import "github.com/google/uuid"

// RegisterInput represents user registration data
type RegisterInput struct {
	Email       string
	PhoneNumber string
	Password    string
	FirstName   string
	LastName    string
}

// LoginInput represents user login data
type LoginInput struct {
	Email    string
	Password string
}

// UserOutput represents user output data
type UserOutput struct {
	ID          uuid.UUID
	Email       *string
	PhoneNumber string
	FirstName   string
	LastName    string
	AvatarURL   string
}

// UpdateProfileInput represents data for updating a user profile
type UpdateProfileInput struct {
	FirstName   string
	LastName    string
	PhoneNumber string
	AvatarURL   string
}

// AuthOutput represents authentication results
type AuthOutput struct {
	AccessToken  string
	RefreshToken string
	User         *UserOutput
}

// PresignUploadOutput represents the result of a presigned upload URL request
type PresignUploadOutput struct {
	UploadURL string
	ObjectKey string
}
