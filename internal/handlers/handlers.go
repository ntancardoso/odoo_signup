package handlers

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"odoo-signup/internal/integration/odoo"
	"odoo-signup/internal/models"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/sirupsen/logrus"
)

// Handler holds dependencies for HTTP handlers
type Handler struct {
	config     *models.Config
	odooClient *odoo.Client
	validate   *validator.Validate
}

// NewHandler creates a new handler instance
func NewHandler(config *models.Config, odooClient *odoo.Client) *Handler {
	return &Handler{
		config:     config,
		odooClient: odooClient,
		validate:   validator.New(),
	}
}

// HandleSignup handles user signup requests
func (h *Handler) HandleSignup(c *gin.Context) {
	var req models.SignupRequest

	// Determine database mode from query param or config
	dbMode := c.DefaultQuery("db_mode", h.config.DefaultDBMode)

	// Bind JSON
	if err := c.ShouldBindJSON(&req); err != nil {
		logrus.WithError(err).Warn("Invalid JSON in signup request")
		c.JSON(http.StatusBadRequest, models.SignupResponse{
			Success: false,
			Message: "Invalid request format",
		})
		return
	}

	// Validate request
	if err := h.validate.Struct(req); err != nil {
		logrus.WithError(err).Warn("Validation failed for signup request")
		c.JSON(http.StatusBadRequest, models.SignupResponse{
			Success: false,
			Message: "Validation failed: " + err.Error(),
		})
		return
	}

	// Sanitize and validate username
	req.Username = strings.ToLower(strings.TrimSpace(req.Username))
	if req.Username == "admin" || req.Username == "www" {
		c.JSON(http.StatusBadRequest, models.SignupResponse{
			Success: false,
			Message: "Username not allowed",
		})
		return
	}

	// Generate database name (use just the username)
	dbName := req.Username

	logger := logrus.WithFields(logrus.Fields{
		"username": req.Username,
		"email":    req.Email,
		"database": dbName,
	})
	logger.Info("Processing signup request")

	// Generate unique RPC ID for this signup request
	rpcID := int(time.Now().UnixNano() % 1000000)

	// Check if database already exists by attempting authentication
	exists, err := h.odooClient.DatabaseExists(dbName)
	if err != nil {
		logger.WithError(err).Error("Failed to check database existence")
		c.JSON(http.StatusInternalServerError, models.SignupResponse{
			Success: false,
			Message: "Failed to validate database name",
		})
		return
	}

	if exists {
		c.JSON(http.StatusConflict, models.SignupResponse{
			Success: false,
			Message: "Username already taken",
		})
		return
	}

	var uid int

	if dbMode == "create" {
		logger.Info("Creating new database")

		// Create new database
		if err := h.odooClient.CreateNewDatabase(dbName, req.Password, req.Email, req.Country.Code, rpcID); err != nil {
			logger.WithError(err).WithField("database", dbName).Error("Failed to create database")
			c.JSON(http.StatusInternalServerError, models.SignupResponse{
				Success: false,
				Message: "Failed to create database",
			})
			return
		}

		// Poll for database readiness with new user credentials
		maxPollingTime := time.Duration(h.config.TimeoutSeconds) * time.Second
		pollInterval := 3 * time.Second
		startTime := time.Now()

		for {
			elapsed := time.Since(startTime)
			if elapsed > maxPollingTime {
				logger.WithField("elapsed_seconds", elapsed.Seconds()).Error("Database polling timeout exceeded")
				c.JSON(http.StatusInternalServerError, models.SignupResponse{
					Success: false,
					Message: "Database creation timeout",
				})
				return
			}

			logger.WithField("elapsed_seconds", elapsed.Seconds()).Debug("Checking if database is ready...")
			tempUID, authErr := h.odooClient.Login(dbName, req.Email, req.Password, rpcID)
			if authErr != nil {
				logger.WithError(authErr).WithField("elapsed_seconds", elapsed.Seconds()).Debug("Database not ready yet, retrying...")
			} else {
				logger.WithField("elapsed_seconds", elapsed.Seconds()).Info("Database is now ready and accessible")
				uid = int(tempUID)
				break
			}

			time.Sleep(pollInterval)
		}

		// Update company details with new user credentials
		companyData := map[string]interface{}{
			"name":  req.CompanyName,
			"email": req.Email,
		}

		if req.Phone != "" {
			companyData["phone"] = req.Phone
		}

		_, err = h.odooClient.ExecuteKw(dbName, uid, req.Password, "res.company", "write", []interface{}{[]interface{}{1}, companyData}, rpcID)
		if err != nil {
			logger.WithError(err).Error("Failed to update company details")
			c.JSON(http.StatusInternalServerError, models.SignupResponse{
				Success: false,
				Message: "Database created but company update failed",
			})
			return
		}

		logger.Info("Company details updated successfully")

	} else { // clone mode
		logger.Info("Cloning database from template")

		// Clone database
		if err := h.odooClient.CloneDatabase(h.config.TemplateDatabase, dbName, rpcID); err != nil {
			logger.WithError(err).WithField("database", dbName).Error("Failed to clone database")
			c.JSON(http.StatusInternalServerError, models.SignupResponse{
				Success: false,
				Message: "Failed to clone database",
			})
			return
		}

		// Poll for database readiness with admin credentials
		maxPollingTime := time.Duration(h.config.TimeoutSeconds) * time.Second
		pollInterval := 3 * time.Second
		startTime := time.Now()

		adminReady := false
		for {
			elapsed := time.Since(startTime)
			if elapsed > maxPollingTime {
				logger.WithField("elapsed_seconds", elapsed.Seconds()).Error("Database polling timeout exceeded")
				c.JSON(http.StatusInternalServerError, models.SignupResponse{
					Success: false,
					Message: "Database cloning timeout",
				})
				return
			}

			logger.WithField("elapsed_seconds", elapsed.Seconds()).Debug("Checking if cloned database is ready...")
			tempUID, authErr := h.odooClient.Login(dbName, h.config.AdminUser, h.config.AdminPassword, rpcID)
			if authErr != nil {
				logger.WithError(authErr).WithField("elapsed_seconds", elapsed.Seconds()).Debug("Database not ready yet, retrying...")
			} else {
				logger.WithField("elapsed_seconds", elapsed.Seconds()).Info("Cloned database is now ready")
				uid = int(tempUID)
				adminReady = true
				break
			}

			time.Sleep(pollInterval)
		}

		if !adminReady {
			c.JSON(http.StatusInternalServerError, models.SignupResponse{
				Success: false,
				Message: "Database cloned but admin authentication failed",
			})
			return
		}

		// Create new user
		userData := map[string]interface{}{
			"name":       fmt.Sprintf("%s %s", req.FirstName, req.LastName),
			"login":      req.Email,
			"password":   req.Password,
			"email":      req.Email,
			"active":     true,
			"company_id": 1,
			"groups_id": []interface{}{
				[]interface{}{6, 0, []interface{}{1, 2, 4}}, // Admin group
			},
		}

		_, err = h.odooClient.ExecuteKw(dbName, uid, h.config.AdminPassword, "res.users", "create", []interface{}{userData}, rpcID)
		if err != nil {
			logger.WithError(err).Error("Failed to create new user")
			c.JSON(http.StatusInternalServerError, models.SignupResponse{
				Success: false,
				Message: "Database cloned but user creation failed",
			})
			return
		}

		logger.Info("New user created successfully")

		// Update company details
		companyData := map[string]interface{}{
			"name":  req.CompanyName,
			"email": req.Email,
		}

		if req.Phone != "" {
			companyData["phone"] = req.Phone
		}

		if req.Country.ID > 0 {
			companyData["country_id"] = req.Country.ID
		}

		// First search for country ID
		// countryIDResult, err := h.odooClient.ExecuteKw(dbName, uid, h.config.AdminPassword, "res.country", "search", []interface{}{[]interface{}{"name", "=", req.Country}}, rpcID)
		//if err != nil {
		//		logger.WithError(err).Warn("Failed to find country ID, skipping country update")
		//} else {
		// if countryIDs, ok := countryIDResult.([]interface{}); ok && len(countryIDs) > 0 {
		//	if id, ok := countryIDs[0].(float64); ok {
		//		companyData["country_id"] = int(id)
		//	}
		//}
		// }

		_, err = h.odooClient.ExecuteKw(dbName, uid, h.config.AdminPassword, "res.company", "write", []interface{}{[]interface{}{1}, companyData}, rpcID)
		if err != nil {
			logger.WithError(err).Error("Failed to update company details")
			c.JSON(http.StatusInternalServerError, models.SignupResponse{
				Success: false,
				Message: "Database cloned and user created but company update failed",
			})
			return
		}

		logger.Info("Company details updated successfully")

	}

	// Success response
	instanceURL := fmt.Sprintf("%s.%s", req.Username, h.config.Domain)

	response := models.SignupResponse{
		Success: true,
		Message: fmt.Sprintf("Signup successful using %s mode! Your Odoo instance is ready.", dbMode),
		Data: &models.SignupData{
			InstanceURL: instanceURL,
			Email:       req.Email,
			Database:    dbName,
		},
	}

	logger.WithFields(logrus.Fields{
		"username":    req.Username,
		"email":       req.Email,
		"db_mode":     dbMode,
		"instanceURL": instanceURL,
	}).Info("Signup completed successfully")

	c.JSON(http.StatusOK, response)
}

// HandleHealthCheck handles health check requests
func (h *Handler) HandleHealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":    "healthy",
		"timestamp": time.Now().UTC(),
	})
}
