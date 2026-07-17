package config

import (
	"log/slog"
	"os"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/joho/godotenv"
)

type EnvConfig struct {
	// Structural constraints using declarative tags
	Environment string `validate:"required,oneof=development staging production"`
	ServerPort  int    `validate:"required,min=1024,max=65535"`
	DBHost      string `validate:"required,hostname|ip"`
	DBPort      int    `validate:"required,min=1,max=65535"`
	DBUser      string `validate:"required"`
	DBPassword  string `validate:"required"`
	DBName      string `validate:"required"`
	JWTSecret   string `validate:"required,min=32"`
}

func getEnvWithDefault(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

// Helper: type conversion wrapper with standard fallback checks
func getEnvAsInt(key string, defaultValue int) int {
	valueStr, exists := os.LookupEnv(key)
	if !exists {
		return defaultValue
	}
	value, err := strconv.Atoi(valueStr)
	if err != nil {
		slog.Warn("Failed to cast string config value to integer", "key", key, "fallback", defaultValue)
		return defaultValue
	}
	return value
}
func LoadConfig() *EnvConfig {
	// Loading the .env file
	err := godotenv.Load()
	if err != nil {
		slog.Debug("No local .env file found. Utilizing standard environment space.")
	}
	cfg := &EnvConfig{
		Environment: getEnvWithDefault("APP_ENV", "development"),
		ServerPort:  getEnvAsInt("SERVER_PORT", 8080),
		DBHost:      getEnvWithDefault("DB_HOST", "127.0.0.1"),
		DBPort:      getEnvAsInt("DB_PORT", 5432),
		DBUser:      os.Getenv("DB_USER"),
		DBPassword:  os.Getenv("DB_PASSWORD"), //these have no default;they are required secrets.
		DBName:      os.Getenv("DB_NAME"),
		JWTSecret:   os.Getenv("JWT_SECRET"),
	}
	validate := validator.New()
	if err := validate.Struct(cfg); err != nil {
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			for _, fieldErr := range validationErrors {
				slog.Error("Critical initialization check faliled",    //here FieldError is an interface with methods like Field(),Tag(),Param()
					"variable", fieldErr.Field(),                      //ValidationErrors = concrete named type = []FieldError
					"violation", fieldErr.Tag(),
					"param", fieldErr.Param(),
				)
			}
		}
		slog.Error("Configuration loading failed.")
		os.Exit(1)
	}
	return cfg

}
