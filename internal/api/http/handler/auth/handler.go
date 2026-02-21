package auth

import (
	"tms-core-service/internal/api/http/dto"
	"tms-core-service/internal/api/http/middleware"
	"tms-core-service/internal/usecase/auth"
	"tms-core-service/internal/util/httpresponse"
	"tms-core-service/internal/util/validator"

	"github.com/gofiber/fiber/v2"
)

// Handler handles authentication requests
type Handler struct {
	useCase       *auth.AuthUseCase
	googleUseCase *auth.GoogleAuthUseCase
	lineUseCase   *auth.LineAuthUseCase
	frontendURL   string
}

// NewHandler creates a new auth handler
func NewHandler(useCase *auth.AuthUseCase, googleUseCase *auth.GoogleAuthUseCase, lineUseCase *auth.LineAuthUseCase, frontendURL string) *Handler {
	return &Handler{
		useCase:       useCase,
		googleUseCase: googleUseCase,
		lineUseCase:   lineUseCase,
		frontendURL:   frontendURL,
	}
}

// Register godoc
// @Summary Register a new user
// @Description Create a new user account
// @Tags auth
// @Accept json
// @Produce json
// @Param request body dto.RegisterRequest true "Registration details"
// @Success 201 {object} httpresponse.Response{data=dto.AuthResponse}
// @Failure 400 {object} httpresponse.Response
// @Failure 409 {object} httpresponse.Response
// @Failure 500 {object} httpresponse.Response
// @Router /api/v1/auth/register [post]
func (h *Handler) Register(c *fiber.Ctx) error {
	var req dto.RegisterRequest
	if err := c.BodyParser(&req); err != nil {
		return httpresponse.Error(c, err)
	}

	// Validate request body
	if err := validator.Validate(req); err != nil {
		return httpresponse.Error(c, err)
	}

	result, err := h.useCase.Register(c.Context(), auth.RegisterInput{
		Email:       req.Email,
		Password:    req.Password,
		FirstName:   req.FirstName,
		LastName:    req.LastName,
		PhoneNumber: req.PhoneNumber,
	})
	if err != nil {
		return httpresponse.Error(c, err)
	}

	return httpresponse.Created(c, dto.AuthResponse{
		AccessToken:  result.AccessToken,
		RefreshToken: result.RefreshToken,
		User: &dto.UserResponse{
			ID:          result.User.ID.String(),
			Email:       result.User.Email,
			FirstName:   result.User.FirstName,
			LastName:    result.User.LastName,
			PhoneNumber: result.User.PhoneNumber,
			AvatarURL:   result.User.AvatarURL,
		},
	}, "User registered successfully")
}

// Login godoc
// @Summary Login
// @Description Authenticate user and get tokens
// @Tags auth
// @Accept json
// @Produce json
// @Param request body dto.LoginRequest true "Login credentials"
// @Success 200 {object} httpresponse.Response{data=dto.AuthResponse}
// @Failure 400 {object} httpresponse.Response
// @Failure 401 {object} httpresponse.Response
// @Failure 500 {object} httpresponse.Response
// @Router /api/v1/auth/login [post]
func (h *Handler) Login(c *fiber.Ctx) error {
	var req dto.LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return httpresponse.Error(c, err)
	}

	// Validate request body
	if err := validator.Validate(req); err != nil {
		return httpresponse.Error(c, err)
	}

	result, err := h.useCase.Login(c.Context(), auth.LoginInput{
		Email:    req.Email,
		Password: req.Password,
	})
	if err != nil {
		return httpresponse.Error(c, err)
	}

	return httpresponse.Success(c, dto.AuthResponse{
		AccessToken:  result.AccessToken,
		RefreshToken: result.RefreshToken,
		User: &dto.UserResponse{
			ID:          result.User.ID.String(),
			Email:       result.User.Email,
			FirstName:   result.User.FirstName,
			LastName:    result.User.LastName,
			PhoneNumber: result.User.PhoneNumber,
			AvatarURL:   result.User.AvatarURL,
		},
	}, "Login successful")
}

// GoogleLogin godoc
// @Summary Google Login
// @Description Redirect to Google OAuth login page
// @Tags auth
// @Success 302
// @Router /api/v1/auth/google/login [get]
func (h *Handler) GoogleLogin(c *fiber.Ctx) error {
	url := h.googleUseCase.GetGoogleLoginURL("random-state") // In production, use a secure random state
	return c.Redirect(url)
}

// GoogleCallback godoc
// @Summary Google OAuth Callback
// @Description Handle Google OAuth redirect and authenticate user
// @Tags auth
// @Accept json
// @Produce json
// @Param code query string true "Authorization code"
// @Success 200 {object} httpresponse.Response{data=dto.AuthResponse}
// @Failure 401 {object} httpresponse.Response
// @Failure 500 {object} httpresponse.Response
// @Router /api/v1/auth/google/callback [get]
func (h *Handler) GoogleCallback(c *fiber.Ctx) error {
	code := c.Query("code")
	if code == "" {
		return c.Redirect(h.frontendURL + "/signin?error=no_code")
	}

	result, err := h.googleUseCase.HandleGoogleCallback(c.Context(), code)
	if err != nil {
		return c.Redirect(h.frontendURL + "/signin?error=auth_failed")
	}

	// Redirect to frontend callback with tokens
	return c.Redirect(h.frontendURL + "/auth/callback?token=" + result.AccessToken + "&refresh_token=" + result.RefreshToken)
}

// LineLogin godoc
// @Summary LINE Login
// @Description Redirect to LINE OAuth login page
// @Tags auth
// @Success 302
// @Router /api/v1/auth/line/login [get]
func (h *Handler) LineLogin(c *fiber.Ctx) error {
	url := h.lineUseCase.GetLineLoginURL("random-state") // In production, use a secure random state
	return c.Redirect(url)
}

// LineCallback godoc
// @Summary LINE OAuth Callback
// @Description Handle LINE OAuth redirect and authenticate user
// @Tags auth
// @Accept json
// @Produce json
// @Param code query string true "Authorization code"
// @Success 200 {object} httpresponse.Response{data=dto.AuthResponse}
// @Failure 401 {object} httpresponse.Response
// @Failure 500 {object} httpresponse.Response
// @Router /api/v1/auth/line/callback [get]
func (h *Handler) LineCallback(c *fiber.Ctx) error {
	code := c.Query("code")
	if code == "" {
		return c.Redirect(h.frontendURL + "/signin?error=no_code")
	}

	result, err := h.lineUseCase.HandleLineCallback(c.Context(), code)
	if err != nil {
		return c.Redirect(h.frontendURL + "/signin?error=auth_failed")
	}

	// Redirect to frontend callback with tokens
	return c.Redirect(h.frontendURL + "/auth/callback?token=" + result.AccessToken + "&refresh_token=" + result.RefreshToken)
}

// GetProfile godoc
// @Summary Get User Profile
// @Description Get currently authenticated user's profile
// @Tags auth
// @Accept json
// @Produce json
// @Security Bearer
// @Success 200 {object} httpresponse.Response{data=dto.UserResponse}
// @Failure 401 {object} httpresponse.Response
// @Failure 404 {object} httpresponse.Response
// @Failure 500 {object} httpresponse.Response
// @Router /api/v1/auth/me [get]
func (h *Handler) GetProfile(c *fiber.Ctx) error {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		return httpresponse.Error(c, fiber.ErrUnauthorized)
	}

	user, err := h.useCase.GetProfile(c.Context(), userID)
	if err != nil {
		return httpresponse.Error(c, err)
	}

	return httpresponse.Success(c, dto.UserResponse{
		ID:          user.ID.String(),
		Email:       user.Email,
		FirstName:   user.FirstName,
		LastName:    user.LastName,
		PhoneNumber: user.PhoneNumber,
		AvatarURL:   user.AvatarURL,
	}, "Profile retrieved successfully")
}

// UpdateProfile godoc
// @Summary Update User Profile
// @Description Update the currently authenticated user's profile information
// @Tags auth
// @Accept json
// @Produce json
// @Security Bearer
// @Param request body dto.UpdateProfileRequest true "Profile update details"
// @Success 200 {object} httpresponse.Response{data=dto.UserResponse}
// @Failure 400 {object} httpresponse.Response
// @Failure 401 {object} httpresponse.Response
// @Failure 404 {object} httpresponse.Response
// @Failure 500 {object} httpresponse.Response
// @Router /api/v1/auth/profile [put]
func (h *Handler) UpdateProfile(c *fiber.Ctx) error {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		return httpresponse.Error(c, fiber.ErrUnauthorized)
	}

	var req dto.UpdateProfileRequest
	if err := c.BodyParser(&req); err != nil {
		return httpresponse.Error(c, err)
	}

	if err := validator.Validate(req); err != nil {
		return httpresponse.Error(c, err)
	}

	user, err := h.useCase.UpdateProfile(c.Context(), userID, auth.UpdateProfileInput{
		FirstName:   req.FirstName,
		LastName:    req.LastName,
		PhoneNumber: req.PhoneNumber,
		AvatarURL:   req.AvatarURL,
	})
	if err != nil {
		return httpresponse.Error(c, err)
	}

	return httpresponse.Success(c, dto.UserResponse{
		ID:          user.ID.String(),
		Email:       user.Email,
		FirstName:   user.FirstName,
		LastName:    user.LastName,
		PhoneNumber: user.PhoneNumber,
		AvatarURL:   user.AvatarURL,
	}, "Profile updated successfully")
}

// GetAvatarUploadURL godoc
// @Summary Get Avatar Upload URL
// @Description Get a presigned URL for uploading a profile avatar
// @Tags auth
// @Accept json
// @Produce json
// @Security Bearer
// @Param request body dto.AvatarUploadRequest true "Content type of the image"
// @Success 200 {object} httpresponse.Response{data=dto.PresignUploadResponse}
// @Failure 400 {object} httpresponse.Response
// @Failure 401 {object} httpresponse.Response
// @Failure 500 {object} httpresponse.Response
// @Router /api/v1/auth/avatar/upload-url [post]
func (h *Handler) GetAvatarUploadURL(c *fiber.Ctx) error {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		return httpresponse.Error(c, fiber.ErrUnauthorized)
	}

	var req dto.AvatarUploadRequest
	if err := c.BodyParser(&req); err != nil {
		return httpresponse.Error(c, err)
	}

	if err := validator.Validate(req); err != nil {
		return httpresponse.Error(c, err)
	}

	result, err := h.useCase.GenerateAvatarUploadURL(c.Context(), userID, req.ContentType)
	if err != nil {
		return httpresponse.Error(c, err)
	}

	return httpresponse.Success(c, dto.PresignUploadResponse{
		UploadURL: result.UploadURL,
		ObjectKey: result.ObjectKey,
	}, "Upload URL generated successfully")
}

// RefreshToken godoc
// @Summary Refresh Tokens
// @Description refresh access and refresh tokens
// @Tags auth
// @Accept json
// @Produce json
// @Param request body dto.RefreshTokenRequest true "Refresh token details"
// @Success 200 {object} httpresponse.Response{data=dto.AuthResponse}
// @Failure 401 {object} httpresponse.Response
// @Failure 500 {object} httpresponse.Response
// @Router /api/v1/auth/refresh [post]
func (h *Handler) RefreshToken(c *fiber.Ctx) error {
	var req dto.RefreshTokenRequest
	if err := c.BodyParser(&req); err != nil {
		return httpresponse.Error(c, err)
	}

	// Validate request body
	if err := validator.Validate(req); err != nil {
		return httpresponse.Error(c, err)
	}

	result, err := h.useCase.RefreshToken(c.Context(), req.RefreshToken)
	if err != nil {
		return httpresponse.Error(c, err)
	}

	return httpresponse.Success(c, dto.AuthResponse{
		AccessToken:  result.AccessToken,
		RefreshToken: result.RefreshToken,
		User: &dto.UserResponse{
			ID:          result.User.ID.String(),
			Email:       result.User.Email,
			FirstName:   result.User.FirstName,
			LastName:    result.User.LastName,
			PhoneNumber: result.User.PhoneNumber,
			AvatarURL:   result.User.AvatarURL,
		},
	}, "Token refreshed successfully")
}
