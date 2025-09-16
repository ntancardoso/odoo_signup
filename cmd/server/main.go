package main

import (
	"net/http"

	"odoo-signup/config"
	"odoo-signup/internal/handlers"
	"odoo-signup/internal/integration/odoo"
	"odoo-signup/internal/middleware"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"golang.org/x/time/rate"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		logrus.Fatal("Failed to load configuration:", err)
	}

	// Initialize logger
	setupLogger(cfg.LogLevel)

	// Set Gin mode
	if cfg.Environment == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	// Initialize Odoo client
	odooClient := odoo.NewClient(cfg.OdooURL, cfg.OdooMasterPass, cfg.AdminUser, cfg.AdminPassword, cfg.TimeoutSeconds)

	// Initialize handlers
	handler := handlers.NewHandler(cfg, odooClient)

	// Create Gin router
	r := gin.New()

	// Load HTML templates
	r.LoadHTMLFiles("./static/index.html")

	// Add middleware
	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	r.Use(cors.Default())

	// Rate limiting
	limiter := rate.NewLimiter(cfg.RateLimit, cfg.BurstLimit)

	// Serve static files
	r.Static("/static", "./static")

	// Serve index.html with template rendering
	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", gin.H{
			"Domain":      cfg.Domain,
			"OdooCompany": cfg.OdooCompany,
		})
	})

	// API routes
	api := r.Group("/api")
	api.Use(middleware.RateLimitMiddleware(limiter))
	{
		api.POST("/signup", handler.HandleSignup)
		api.GET("/health", handler.HandleHealthCheck)
	}

	// Start server
	logrus.WithField("port", cfg.Port).Info("Starting Odoo Signup server")
	if err := r.Run(":" + cfg.Port); err != nil {
		logrus.Fatal("Failed to start server:", err)
	}
}

// setupLogger configures the logger based on the log level
func setupLogger(logLevel string) {
	logger := logrus.New()
	logger.SetFormatter(&logrus.JSONFormatter{})

	switch logLevel {
	case "debug":
		logger.SetLevel(logrus.DebugLevel)
	case "info":
		logger.SetLevel(logrus.InfoLevel)
	case "warn":
		logger.SetLevel(logrus.WarnLevel)
	case "error":
		logger.SetLevel(logrus.ErrorLevel)
	default:
		logger.SetLevel(logrus.InfoLevel)
	}

	logrus.SetFormatter(logger.Formatter)
	logrus.SetLevel(logger.Level)
}
