package hash

import (
	"tms-core-service/internal/domain/service"
	"tms-core-service/pkg/hash"
)

type bcryptHashService struct{}

// NewBcryptHashService creates a new bcrypt hash service
func NewBcryptHashService() service.HashService {
	return &bcryptHashService{}
}

func (s *bcryptHashService) HashPassword(password string) (string, error) {
	return hash.HashPassword(password)
}

func (s *bcryptHashService) CheckPassword(password, hashed string) bool {
	return hash.CheckPassword(password, hashed)
}
