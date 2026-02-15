package token

import (
	"time"

	"tms-core-service/internal/domain/service"
	"tms-core-service/pkg/jwt"

	"github.com/google/uuid"
)

type jwtTokenService struct {
	jwtService *jwt.JWTService
}

// NewJWTTokenService creates a new JWT token service
func NewJWTTokenService(jwtService *jwt.JWTService) service.TokenService {
	return &jwtTokenService{jwtService: jwtService}
}

func (s *jwtTokenService) GenerateToken(userID uuid.UUID, email string, expiry time.Duration) (string, error) {
	// We map 'email' to 'status' in the existing package for now,
	// or we update the package/claims. For this refactor, let's just pass it.
	return s.jwtService.GenerateToken(userID, email, expiry)
}
