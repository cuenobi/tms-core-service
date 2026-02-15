package server

import (
	"fmt"

	"tms-core-service/internal/api/http/handler/auth"
	"tms-core-service/internal/api/http/handler/healthcheck"
	"tms-core-service/internal/api/http/route"
	"tms-core-service/internal/config"
	"tms-core-service/internal/infra/db"
	healthcheckRepo "tms-core-service/internal/infra/db/repository/healthcheck"
	userRepo "tms-core-service/internal/infra/db/repository/user"
	"tms-core-service/internal/infra/redis"
	hashSvc "tms-core-service/internal/infra/service/hash"
	tokenSvc "tms-core-service/internal/infra/service/token"
	authUseCase "tms-core-service/internal/usecase/auth"
	healthcheckUseCase "tms-core-service/internal/usecase/healthcheck"
	"tms-core-service/pkg/jwt"

	"github.com/gofiber/fiber/v2"
)

// WireDependencies manually wires all dependencies
func WireDependencies(app *fiber.App, cfg *config.AppConfig) error {
	// Initialize database connection
	dbConn, err := db.NewConnection(&cfg.Database)
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}

	// Initialize Redis connection
	redisClient, err := redis.NewConnection(&cfg.Redis)
	if err != nil {
		return fmt.Errorf("failed to connect to Redis: %w", err)
	}

	// Initialize core packages (concrete implementations)
	jwtProvider := jwt.NewJWTService(cfg.JWT.Secret, cfg.JWT.Issuer)

	// Initialize SOLID service wrappers (Domain Abstractions)
	hashService := hashSvc.NewBcryptHashService()
	tokenService := tokenSvc.NewJWTTokenService(jwtProvider)

	// Initialize repositories
	healthCheckRepo := healthcheckRepo.NewHealthCheckRepository(dbConn)
	userRepository := userRepo.NewUserRepository(dbConn)

	// Initialize cache repository (if needed)
	_ = redis.NewCacheRepository(redisClient)

	// Initialize use cases
	healthCheckUC := healthcheckUseCase.NewHealthCheckUseCase(healthCheckRepo)
	authUC := authUseCase.NewAuthUseCase(
		userRepository,
		hashService,
		tokenService,
		int64(cfg.JWT.AccessTokenExpiry.Minutes()),
		int64(cfg.JWT.RefreshTokenExpiry.Hours()),
	)

	// Initialize handlers
	healthCheckHandler := healthcheck.NewHandler(healthCheckUC)
	authHandler := auth.NewHandler(authUC)

	// Setup routes
	deps := &route.Dependencies{
		HealthCheckHandler: healthCheckHandler,
		AuthHandler:        authHandler,
		JWTService:         jwtProvider,
	}
	route.SetupRoutes(app, deps)

	return nil
}
