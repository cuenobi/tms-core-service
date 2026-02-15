package config

import (
	"fmt"
	"strings"
	"time"

	"github.com/spf13/viper"
)

// AppConfig represents the entire application configuration
type AppConfig struct {
	Server    ServerConfig    `mapstructure:"server"`
	Database  DatabaseConfig  `mapstructure:"database"`
	Redis     RedisConfig     `mapstructure:"redis"`
	JWT       JWTConfig       `mapstructure:"jwt"`
	Migration MigrationConfig `mapstructure:"migration"`
}

// ServerConfig contains HTTP server settings
type ServerConfig struct {
	Port    int           `mapstructure:"port"`
	Mode    string        `mapstructure:"mode"` // debug, release, test
	Timeout TimeoutConfig `mapstructure:"timeout"`
}

// TimeoutConfig contains server timeout settings
type TimeoutConfig struct {
	Read  time.Duration `mapstructure:"read"`
	Write time.Duration `mapstructure:"write"`
	Idle  time.Duration `mapstructure:"idle"`
}

// DatabaseConfig contains database connection settings
type DatabaseConfig struct {
	Host     string     `mapstructure:"host"`
	Port     int        `mapstructure:"port"`
	User     string     `mapstructure:"user"`
	Password string     `mapstructure:"password"`
	DBName   string     `mapstructure:"dbname"`
	SSLMode  string     `mapstructure:"sslmode"`
	TimeZone string     `mapstructure:"timezone"`
	Pool     PoolConfig `mapstructure:"pool"`
}

// PoolConfig contains database connection pool settings
type PoolConfig struct {
	MaxIdleConns    int           `mapstructure:"max_idle_conns"`
	MaxOpenConns    int           `mapstructure:"max_open_conns"`
	ConnMaxLifetime time.Duration `mapstructure:"conn_max_lifetime"`
}

// RedisConfig contains Redis connection settings
type RedisConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Password string `mapstructure:"password"`
	DB       int    `mapstructure:"db"`
	PoolSize int    `mapstructure:"pool_size"`
}

// JWTConfig contains JWT authentication settings
type JWTConfig struct {
	Secret             string        `mapstructure:"secret"`
	Issuer             string        `mapstructure:"issuer"`
	AccessTokenExpiry  time.Duration `mapstructure:"access_token_expiry"`
	RefreshTokenExpiry time.Duration `mapstructure:"refresh_token_expiry"`
}

// MigrationConfig contains database migration settings
type MigrationConfig struct {
	Dir   string `mapstructure:"dir"`
	Table string `mapstructure:"table"`
}

// LoadConfig loads configuration from the specified file
func LoadConfig(configPath string) (*AppConfig, error) {
	viper.SetConfigFile(configPath)
	viper.SetConfigType("yaml")

	// Read from environment variables with prefix TMS_
	viper.SetEnvPrefix("TMS")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

	// Explicitly bind JWT environment variables for rotation/flexibility
	_ = viper.BindEnv("jwt.secret", "JWT_SECRET")
	_ = viper.BindEnv("jwt.access_token_expiry", "JWT_ACCESS_EXP")
	_ = viper.BindEnv("jwt.refresh_token_expiry", "JWT_REFRESH_EXP")

	// Standard K8s/Docker bindings (No prefix)
	_ = viper.BindEnv("server.port", "PORT")
	_ = viper.BindEnv("database.host", "DATABASE_HOST")
	_ = viper.BindEnv("database.port", "DATABASE_PORT")
	_ = viper.BindEnv("database.user", "DATABASE_USER")
	_ = viper.BindEnv("database.password", "DATABASE_PASSWORD")
	_ = viper.BindEnv("database.dbname", "DATABASE_NAME")
	_ = viper.BindEnv("redis.host", "REDIS_HOST")
	_ = viper.BindEnv("redis.port", "REDIS_PORT")
	_ = viper.BindEnv("redis.password", "REDIS_PASSWORD")

	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var config AppConfig
	if err := viper.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	return &config, nil
}

// GetDSN returns the PostgreSQL DSN connection string
func (c *DatabaseConfig) GetDSN() string {
	return fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s TimeZone=%s",
		c.Host, c.Port, c.User, c.Password, c.DBName, c.SSLMode, c.TimeZone,
	)
}

// GetRedisAddr returns the Redis address string
func (c *RedisConfig) GetRedisAddr() string {
	return fmt.Sprintf("%s:%d", c.Host, c.Port)
}
