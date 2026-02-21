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
	storageSvc "tms-core-service/internal/infra/service/storage"
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

	// Initialize S3 storage service
	storageService := storageSvc.NewS3StorageService(
		cfg.S3.Region,
		cfg.S3.Bucket,
		cfg.S3.AccessKey,
		cfg.S3.SecretKey,
		cfg.S3.PresignExpiry,
	)

	// Initialize use cases
	healthCheckUC := healthcheckUseCase.NewHealthCheckUseCase(healthCheckRepo)
	authUC := authUseCase.NewAuthUseCase(
		userRepository,
		hashService,
		tokenService,
		storageService,
		int64(cfg.JWT.AccessTokenExpiry.Minutes()),
		int64(cfg.JWT.RefreshTokenExpiry.Hours()),
	)
	googleAuthUC := authUseCase.NewGoogleAuthUseCase(
		userRepository,
		tokenService,
		cfg.Google.ClientID,
		cfg.Google.ClientSecret,
		cfg.Google.RedirectURL,
		int64(cfg.JWT.AccessTokenExpiry.Minutes()),
		int64(cfg.JWT.RefreshTokenExpiry.Hours()),
	)
	lineAuthUC := authUseCase.NewLineAuthUseCase(
		userRepository,
		tokenService,
		cfg.Line.ChannelID,
		cfg.Line.ChannelSecret,
		cfg.Line.RedirectURL,
		int64(cfg.JWT.AccessTokenExpiry.Minutes()),
		int64(cfg.JWT.RefreshTokenExpiry.Hours()),
	)

	// Initialize handlers
	healthCheckHandler := healthcheck.NewHandler(healthCheckUC)
	authHandler := auth.NewHandler(authUC, googleAuthUC, lineAuthUC, cfg.Server.FrontendURL)

	// Setup routes
	deps := &route.Dependencies{
		HealthCheckHandler: healthCheckHandler,
		AuthHandler:        authHandler,
		JWTService:         jwtProvider,
	}
	route.SetupRoutes(app, deps)

	return nil
}
