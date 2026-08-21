package config

import (
	"os"
	"strconv"
	"time"
)

// Config holds all application configuration
type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	JWT      JWTConfig
	Cache    CacheConfig
	Mining   MiningConfig
	Email    EmailConfig
	S3       S3Config
}

// ServerConfig holds server-related configuration
type ServerConfig struct {
	Host            string
	Port            string
	Environment     string
	ShutdownTimeout time.Duration
	ReadTimeout     time.Duration
	WriteTimeout    time.Duration
}

// DatabaseConfig holds database configuration
type DatabaseConfig struct {
	DSN             string
	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxLifetime time.Duration
	LogLevel        string
}

// JWTConfig holds JWT token configuration
type JWTConfig struct {
	Secret           string
	AccessTokenTTL   time.Duration
	RefreshTokenTTL  time.Duration
	TokenIssuer      string
	TokenAudience    string
}

// CacheConfig holds cache configuration
type CacheConfig struct {
	Enabled        bool
	RedisAddr      string
	DefaultTTL     time.Duration
	MaxRetries     int
}

// MiningConfig holds mining-related configuration
type MiningConfig struct {
	RewardPerFact        string // decimal string
	MinSessionDuration   time.Duration
	MaxCPUUsagePercent   int
	ConsensusThreshold   float64
}

// EmailConfig holds email service configuration
type EmailConfig struct {
	Enabled  bool
	SMTPHost string
	SMTPPort int
	Username string
	Password string
	FromAddr string
}

// S3Config holds AWS S3 configuration
type S3Config struct {
	Enabled       bool
	Region        string
	AccessKeyID   string
	SecretAccessKey string
	BucketName    string
	Endpoint      string
}

// Load loads configuration from environment variables
func Load() *Config {
	return &Config{
		Server: ServerConfig{
			Host:            getEnv("SERVER_HOST", "0.0.0.0"),
			Port:            getEnv("SERVER_PORT", "8080"),
			Environment:     getEnv("ENVIRONMENT", "development"),
			ShutdownTimeout: getDurationEnv("SHUTDOWN_TIMEOUT", 30*time.Second),
			ReadTimeout:     getDurationEnv("READ_TIMEOUT", 10*time.Second),
			WriteTimeout:    getDurationEnv("WRITE_TIMEOUT", 10*time.Second),
		},
		Database: DatabaseConfig{
			DSN:             getEnv("DATABASE_URL", ""),
			MaxOpenConns:    getIntEnv("DB_MAX_OPEN_CONNS", 25),
			MaxIdleConns:    getIntEnv("DB_MAX_IDLE_CONNS", 5),
			ConnMaxLifetime: getDurationEnv("DB_CONN_MAX_LIFETIME", 5*time.Minute),
			LogLevel:        getEnv("DB_LOG_LEVEL", "info"),
		},
		JWT: JWTConfig{
			Secret:          getEnv("JWT_SECRET", "your-secret-key"),
			AccessTokenTTL:  getDurationEnv("ACCESS_TOKEN_TTL", 15*time.Minute),
			RefreshTokenTTL: getDurationEnv("REFRESH_TOKEN_TTL", 7*24*time.Hour),
			TokenIssuer:     getEnv("JWT_ISSUER", "techno-re-ecosystem"),
			TokenAudience:   getEnv("JWT_AUDIENCE", "techno-re-users"),
		},
		Cache: CacheConfig{
			Enabled:    getBoolEnv("CACHE_ENABLED", true),
			RedisAddr:  getEnv("REDIS_ADDR", "localhost:6379"),
			DefaultTTL: getDurationEnv("CACHE_DEFAULT_TTL", 1*time.Hour),
			MaxRetries: getIntEnv("CACHE_MAX_RETRIES", 3),
		},
		Mining: MiningConfig{
			RewardPerFact:      getEnv("MINING_REWARD_PER_FACT", "10.5"),
			MinSessionDuration: getDurationEnv("MINING_MIN_SESSION_DURATION", 5*time.Minute),
			MaxCPUUsagePercent: getIntEnv("MINING_MAX_CPU_PERCENT", 50),
			ConsensusThreshold: getFloat64Env("MINING_CONSENSUS_THRESHOLD", 0.8),
		},
		Email: EmailConfig{
			Enabled:  getBoolEnv("EMAIL_ENABLED", false),
			SMTPHost: getEnv("SMTP_HOST", "smtp.gmail.com"),
			SMTPPort: getIntEnv("SMTP_PORT", 587),
			Username: getEnv("SMTP_USERNAME", ""),
			Password: getEnv("SMTP_PASSWORD", ""),
			FromAddr: getEnv("EMAIL_FROM_ADDR", "noreply@techno-re.io"),
		},
		S3: S3Config{
			Enabled:         getBoolEnv("S3_ENABLED", false),
			Region:          getEnv("S3_REGION", "us-east-1"),
			AccessKeyID:     getEnv("S3_ACCESS_KEY_ID", ""),
			SecretAccessKey: getEnv("S3_SECRET_ACCESS_KEY", ""),
			BucketName:      getEnv("S3_BUCKET_NAME", ""),
			Endpoint:        getEnv("S3_ENDPOINT", ""),
		},
	}
}

// Helper functions
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getIntEnv(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intVal, err := strconv.Atoi(value); err == nil {
			return intVal
		}
	}
	return defaultValue
}

func getBoolEnv(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		if boolVal, err := strconv.ParseBool(value); err == nil {
			return boolVal
		}
	}
	return defaultValue
}

func getDurationEnv(key string, defaultValue time.Duration) time.Duration {
	if value := os.Getenv(key); value != "" {
		if duration, err := time.ParseDuration(value); err == nil {
			return duration
		}
	}
	return defaultValue
}

func getFloat64Env(key string, defaultValue float64) float64 {
	if value := os.Getenv(key); value != "" {
		if floatVal, err := strconv.ParseFloat(value, 64); err == nil {
			return floatVal
		}
	}
	return defaultValue
}
