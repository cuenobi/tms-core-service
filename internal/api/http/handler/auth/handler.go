package auth

import (
	"tms-core-service/internal/api/http/dto"
	"tms-core-service/internal/usecase/auth"
	"tms-core-service/internal/util/httpresponse"
	"tms-core-service/internal/util/validator"

	"github.com/gofiber/fiber/v2"
)

// Handler handles authentication requests
type Handler struct {
	useCase *auth.AuthUseCase
}

// NewHandler creates a new auth handler
func NewHandler(useCase *auth.AuthUseCase) *Handler {
	return &Handler{useCase: useCase}
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
		},
	}, "Login successful")
}
