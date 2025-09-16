package config

import (
	"os"
	"strconv"

	"odoo-signup/internal/models"

	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"golang.org/x/time/rate"
)

// Load loads configuration from environment variables
func Load() (*models.Config, error) {
	// Load environment variables from .env file if it exists
	if err := godotenv.Load(); err != nil {
		logrus.Warn("No .env file found, using environment variables")
	}

	config := &models.Config{
		Port:             getEnv("PORT", "8080"),
		OdooURL:          getEnv("ODOO_URL", "http://localhost:8069"),
		OdooMasterPass:   getEnv("ODOO_MASTER_PASSWORD", ""),
		OdooCompany:      getEnv("ODOO_COMPANY", "Sample"),
		Environment:      getEnv("ENVIRONMENT", "development"),
		Domain:           getEnv("DOMAIN", "odoo.the9o.com"),
		TemplateDatabase: getEnv("TEMPLATE_DATABASE", "odoo-template"),
		AdminUser:        getEnv("ADMIN_USER", "admin"),
		AdminPassword:    getEnv("ADMIN_PASSWORD", "admin"),
		LogLevel:         getEnv("LOG_LEVEL", "info"),
	}

	// Parse rate limiting
	rateLimitStr := getEnv("RATE_LIMIT", "10")
	burstLimitStr := getEnv("BURST_LIMIT", "20")

	if rateLimit, err := strconv.ParseFloat(rateLimitStr, 64); err == nil {
		config.RateLimit = rate.Limit(rateLimit)
	} else {
		config.RateLimit = rate.Limit(10)
	}

	if burstLimit, err := strconv.Atoi(burstLimitStr); err == nil {
		config.BurstLimit = burstLimit
	} else {
		config.BurstLimit = 20
	}

	// Parse timeout configuration
	timeoutStr := getEnv("HTTP_TIMEOUT_SECONDS", "300") // Default 5 minutes
	if timeoutSeconds, err := strconv.Atoi(timeoutStr); err == nil {
		config.TimeoutSeconds = timeoutSeconds
	} else {
		config.TimeoutSeconds = 300 // Default 5 minutes
	}

	// Validate required configuration
	if config.OdooMasterPass == "" {
		logrus.Fatal("ODOO_MASTER_PASSWORD environment variable is required")
	}

	return config, nil
}

// getEnv gets an environment variable with a fallback value
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
