package odoo

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
)

// Client represents an Odoo client
type Client struct {
	baseURL    string
	masterPass string
	adminUser  string
	adminPass  string
	httpClient *http.Client
}

// NewClient creates a new Odoo client
func NewClient(baseURL, masterPass, adminUser, adminPass string, timeoutSeconds int) *Client {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	return &Client{
		baseURL:    strings.TrimSuffix(baseURL, "/"),
		masterPass: masterPass,
		adminUser:  adminUser,
		adminPass:  adminPass,
		httpClient: &http.Client{
			Timeout:   time.Duration(timeoutSeconds) * time.Second,
			Transport: tr,
		},
	}
}

func (c *Client) Login(dbName, username, password string, rpcID int) (int, error) {
	logger := logrus.WithFields(logrus.Fields{
		"database": dbName,
		"username": username,
	})
	logger.Info("Logging in to get UID")

	payload := map[string]interface{}{
		"jsonrpc": "2.0",
		"method":  "call",
		"params": map[string]interface{}{
			"service": "common",
			"method":  "login",
			"args":    []interface{}{dbName, username, password},
		},
		"id": rpcID,
	}

	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return 0, fmt.Errorf("failed to marshal login payload: %w", err)
	}

	req, err := http.NewRequest("POST", c.baseURL+"/jsonrpc", bytes.NewBuffer(jsonPayload))
	if err != nil {
		return 0, fmt.Errorf("failed to create login request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return 0, fmt.Errorf("login request failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, fmt.Errorf("failed to read login response: %w", err)
	}

	var response map[string]interface{}
	if err := json.Unmarshal(body, &response); err != nil {
		return 0, fmt.Errorf("failed to parse login response: %w", err)
	}

	if result, ok := response["result"]; ok {
		if uid, ok := result.(float64); ok && uid > 0 {
			logger.WithField("uid", int(uid)).Info("Login successful")
			return int(uid), nil
		}
	}

	logger.Error("Login failed - no valid UID")
	return 0, fmt.Errorf("login failed")
}

// DatabaseExists checks if a database exists by attempting login
func (c *Client) DatabaseExists(dbName string) (bool, error) {
	logger := logrus.WithField("database", dbName)
	logger.Info("Checking if database exists by attempting login")

	// Use a dummy RPC ID
	rpcID := int(time.Now().UnixNano() % 1000000)

	_, err := c.Login(dbName, c.adminUser, c.adminPass, rpcID)
	if err != nil {
		logger.Debug("Database does not exist or login failed")
		return false, nil
	}

	logger.Info("Database exists - login successful")
	return true, nil
}

// CloneDatabase clones an existing database to create a new one using JSON-RPC
func (c *Client) CloneDatabase(templateDbName, newDbName string, rpcID int) error {
	logger := logrus.WithFields(logrus.Fields{
		"template_database": templateDbName,
		"new_database":      newDbName,
	})
	logger.Info("Cloning Odoo database using JSON-RPC")

	payload := map[string]interface{}{
		"jsonrpc": "2.0",
		"method":  "call",
		"params": map[string]interface{}{
			"service": "db",
			"method":  "duplicate_database",
			"args":    []interface{}{c.masterPass, templateDbName, newDbName},
		},
		"id": rpcID,
	}

	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal clone payload: %w", err)
	}

	req, err := http.NewRequest("POST", c.baseURL+"/jsonrpc", bytes.NewBuffer(jsonPayload))
	if err != nil {
		return fmt.Errorf("failed to create clone request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to clone database: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read clone response: %w", err)
	}

	var response map[string]interface{}
	if err := json.Unmarshal(body, &response); err != nil {
		return fmt.Errorf("failed to parse clone response: %w", err)
	}

	if result, ok := response["result"]; ok && result == true {
		logger.WithField("new_database", newDbName).Info("Database cloned successfully")
		return nil
	}

	logger.WithField("response", response).Error("Clone response indicates failure")
	return fmt.Errorf("database clone failed: %v", response)
}

// CreateNewDatabase creates a new database and admin user using JSON-RPC
func (c *Client) CreateNewDatabase(dbName, password, login, country string, rpcID int) error {
	logger := logrus.WithFields(logrus.Fields{
		"database": dbName,
		"login":    login,
	})
	logger.Info("Creating new Odoo database using JSON-RPC")

	payload := map[string]interface{}{
		"jsonrpc": "2.0",
		"method":  "call",
		"params": map[string]interface{}{
			"service": "db",
			"method":  "create_database",
			"args":    []interface{}{c.masterPass, dbName, false, "en_US", password, login, country},
		},
		"id": rpcID,
	}

	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal create database payload: %w", err)
	}

	req, err := http.NewRequest("POST", c.baseURL+"/jsonrpc", bytes.NewBuffer(jsonPayload))
	if err != nil {
		return fmt.Errorf("failed to create database request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to create database: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read create database response: %w", err)
	}

	var response map[string]interface{}
	if err := json.Unmarshal(body, &response); err != nil {
		return fmt.Errorf("failed to parse create database response: %w", err)
	}

	if result, ok := response["result"]; ok && result == true {
		logger.WithField("database", dbName).Info("Database created successfully")
		return nil
	}

	logger.WithField("response", response).Error("Create database response indicates failure")
	return fmt.Errorf("database creation failed: %v", response)
}

// ExecuteKw executes a kw method on a model
func (c *Client) ExecuteKw(dbName string, uid int, password string, model string, method string, args []interface{}, rpcID int) (interface{}, error) {
	logger := logrus.WithFields(logrus.Fields{
		"database": dbName,
		"model":    model,
		"method":   method,
	})
	logger.Info("Executing kw method")

	payload := map[string]interface{}{
		"jsonrpc": "2.0",
		"method":  "call",
		"params": map[string]interface{}{
			"service": "object",
			"method":  "execute_kw",
			"args": []interface{}{
				dbName,
				uid,
				password,
				model,
				method,
				args,
			},
		},
		"id": rpcID,
	}

	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal execute_kw payload: %w", err)
	}

	req, err := http.NewRequest("POST", c.baseURL+"/jsonrpc", bytes.NewBuffer(jsonPayload))
	if err != nil {
		return nil, fmt.Errorf("failed to create execute_kw request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("execute_kw request failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read execute_kw response: %w", err)
	}

	var response map[string]interface{}
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("failed to parse execute_kw response: %w", err)
	}

	if result, ok := response["result"]; ok {
		logger.Info("Execute_kw successful")
		return result, nil
	}

	logger.WithField("response", response).Error("Execute_kw failed")
	return nil, fmt.Errorf("execute_kw failed: %v", response)
}

// CreateDatabase creates a new database by duplicating from template
//func (c *Client) CreateDatabase(templateDbName, newDbName string, rpcID int) error {
//	return c.CloneDatabase(templateDbName, newDbName, rpcID)
//}

// CreateAdminUser creates a new admin user in the specified database
// func (c *Client) CreateAdminUser(dbName, email, password, firstName, lastName, domain string, rpcID int) error {
// 	return c.CreateAdminUserWithSession(dbName, email, password, firstName, lastName, domain, "", rpcID)
// }

// // CreateAdminUserWithSession creates a new admin user using an existing session
// func (c *Client) CreateAdminUserWithSession(dbName, email, password, firstName, lastName, domain, sessionID string, rpcID int) error {
// 	logger := logrus.WithFields(logrus.Fields{
// 		"database": dbName,
// 		"email":    email,
// 	})
// 	logger.Info("Creating admin user in database")
//
// 	if sessionID == "" {
// 		var err error
// 		sessionID, err = c.Authenticate(dbName, c.adminUser, c.adminPass, domain, rpcID)
// 		if err != nil {
// 			return fmt.Errorf("authentication failed: %w", err)
// 		}
// 	}
//
// 	instanceURL := fmt.Sprintf("https://%s.%s", dbName, domain)
// 	userURL := instanceURL + "/web/dataset/call_kw"
//
// 	userPayload := map[string]interface{}{
// 		"jsonrpc": "2.0",
// 		"method":  "call",
// 		"params": map[string]interface{}{
// 			"model":  "res.users",
// 			"method": "create",
// 			"args": []interface{}{
// 				map[string]interface{}{
// 					"name":       fmt.Sprintf("%s %s", firstName, lastName),
// 					"login":      email,
// 					"password":   password,
// 					"email":      email,
// 					"active":     true,
// 					"company_id": 1,
// 					"groups_id": []interface{}{
// 						[]interface{}{6, 0, []interface{}{1, 2, 4}}, // Admin group
// 					},
// 				},
// 			},
// 			"kwargs": map[string]interface{}{
// 				"context": map[string]interface{}{
// 					"tz":   "UTC",
// 					"lang": "en_US",
// 				},
// 			},
// 		},
// 	}
//
// 	jsonPayload, err := json.Marshal(userPayload)
// 	if err != nil {
// 		return fmt.Errorf("failed to marshal user payload: %w", err)
// 	}
//
// 	req, err := http.NewRequest("POST", userURL, bytes.NewBuffer(jsonPayload))
// 	if err != nil {
// 		return fmt.Errorf("failed to create user request: %w", err)
// 	}
//
// 	req.Header.Set("Content-Type", "application/json")
// 	req.Header.Set("Cookie", fmt.Sprintf("session_id=%s", sessionID))
//
// 	resp, err := c.httpClient.Do(req)
// 	if err != nil {
// 		return fmt.Errorf("user creation request failed: %w", err)
// 	}
// 	defer resp.Body.Close()
//
// 	body, err := io.ReadAll(resp.Body)
// 	if err != nil {
// 		return fmt.Errorf("failed to read user response: %w", err)
// 	}
//
// 	var response map[string]interface{}
// 	if err := json.Unmarshal(body, &response); err != nil {
// 		return fmt.Errorf("failed to parse user response: %w", err)
// 	}
//
// 	if result, ok := response["result"]; ok {
// 		if userID, ok := result.(float64); ok && userID > 0 {
// 			logger.WithField("user_id", userID).Info("User created successfully")
// 			return nil
// 		}
// 	}
//
// 	logger.WithField("response", response).Error("User creation response indicates failure")
// 	return fmt.Errorf("user creation failed: %v", response)
// }

// UpdateCompanyDetails updates company information in the specified database
// func (c *Client) UpdateCompanyDetails(dbName, companyName, email, phone, industry, companySize, country, domain string) error {
// 	return c.UpdateCompanyDetailsWithSession(dbName, companyName, email, phone, industry, companySize, country, domain, "")
// }

// // UpdateCompanyDetailsWithSession updates company information using an existing session
// func (c *Client) UpdateCompanyDetailsWithSession(dbName, companyName, email, phone, industry, companySize, country, domain, sessionID string) error {
// 	logger := logrus.WithFields(logrus.Fields{
// 		"database": dbName,
// 		"company":  companyName,
// 	})
// 	logger.Info("Updating company details in database")
//
// 	if sessionID == "" {
// 		var err error
// 		sessionID, err = c.Authenticate(dbName, c.adminUser, c.adminPass, domain, 0)
// 		if err != nil {
// 			return fmt.Errorf("authentication failed: %w", err)
// 		}
// 	}
//
// 	instanceURL := fmt.Sprintf("https://%s.%s", dbName, domain)
// 	companyURL := instanceURL + "/web/dataset/call_kw"
//
// 	// Get default company ID
// 	searchPayload := map[string]interface{}{
// 		"jsonrpc": "2.0",
// 		"method":  "call",
// 		"params": map[string]interface{}{
// 			"model":  "res.company",
// 		"method": "search",
// 		"args": []interface{}{
// 				[]interface{}{},
// 			},
// 			"kwargs": map[string]interface{}{
// 				"context": map[string]interface{}{
// 					"tz":   "UTC",
// 					"lang": "en_US",
// 				},
// 				"limit": 1,
// 			},
// 		},
// 	}
//
// 	jsonSearch, err := json.Marshal(searchPayload)
// 	if err != nil {
// 		return fmt.Errorf("failed to marshal search payload: %w", err)
// 	}
//
// 	req, err := http.NewRequest("POST", companyURL, bytes.NewBuffer(jsonSearch))
// 	if err != nil {
// 		return fmt.Errorf("failed to create search request: %w", err)
// 	}
//
// 	req.Header.Set("Content-Type", "application/json")
// 	req.Header.Set("Cookie", fmt.Sprintf("session_id=%s", sessionID))
//
// 	resp, err := c.httpClient.Do(req)
// 	if err != nil {
// 		return fmt.Errorf("company search failed: %w", err)
// 	}
// 	defer resp.Body.Close()
//
// 	body, err := io.ReadAll(resp.Body)
// 	if err != nil {
// 		return fmt.Errorf("failed to read search response: %w", err)
// 	}
//
// 	var searchResponse map[string]interface{}
// 	if err := json.Unmarshal(body, &searchResponse); err != nil {
// 		return fmt.Errorf("failed to parse search response: %w", err)
// 	}
//
// 	var companyID int = 1
// 	if result, ok := searchResponse["result"].([]interface{}); ok && len(result) > 0 {
// 		if id, ok := result[0].(float64); ok {
// 			companyID = int(id)
// 		}
// 	}
//
// 	// Update company
// 	companyData := map[string]interface{}{
// 		"name":  companyName,
// 		"email": email,
// 	}
//
// 	if phone != "" {
// 		companyData["phone"] = phone
// 	}
// 	// Skip industry and companySize for now as these fields don't exist on res.company model
//
// 	updatePayload := map[string]interface{}{
// 		"jsonrpc": "2.0",
// 		"method":  "call",
// 		"params": map[string]interface{}{
// 			"model":  "res.company",
// 			"method": "write",
// 			"args": []interface{}{
// 				[]interface{}{companyID},
// 				companyData,
// 			},
// 			"kwargs": map[string]interface{}{
// 				"context": map[string]interface{}{
// 					"tz":   "UTC",
// 					"lang": "en_US",
// 				},
// 			},
// 		},
// 	}
//
// 	jsonUpdate, err := json.Marshal(updatePayload)
// 	if err != nil {
// 		return fmt.Errorf("failed to marshal update payload: %w", err)
// 	}
//
// 	req, err = http.NewRequest("POST", companyURL, bytes.NewBuffer(jsonUpdate))
// 	if err != nil {
// 		return fmt.Errorf("failed to create update request: %w", err)
// 	}
//
// 	req.Header.Set("Content-Type", "application/json")
// 	req.Header.Set("Cookie", fmt.Sprintf("session_id=%s", sessionID))
//
// 	resp, err = c.httpClient.Do(req)
// 	if err != nil {
// 		return fmt.Errorf("company update failed: %w", err)
// 	}
// 	defer resp.Body.Close()
//
// 	body, err = io.ReadAll(resp.Body)
// 	if err != nil {
// 		return fmt.Errorf("failed to read update response: %w", err)
// 	}
//
// 	var updateResponse map[string]interface{}
// 	if err := json.Unmarshal(body, &updateResponse); err != nil {
// 		return fmt.Errorf("failed to parse update response: %w", err)
// 	}
//
// 	if result, ok := updateResponse["result"]; ok && result == true {
// 		logger.WithField("company_id", companyID).Info("Company details updated successfully")
// 		return nil
// 	}
//
// 	logger.WithField("response", updateResponse).Error("Company update response indicates failure")
// 	return fmt.Errorf("company update failed: %v", updateResponse)
// }

// TestConnection performs basic connectivity tests
func (c *Client) TestConnection() map[string]interface{} {
	return map[string]interface{}{
		"status":  "ok",
		"message": "Basic connectivity test",
	}
}
