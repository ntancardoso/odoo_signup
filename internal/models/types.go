package models

import "golang.org/x/time/rate"

// Config holds application configuration
type Config struct {
	Port             string
	OdooURL          string
	OdooMasterPass   string
	OdooCompany      string
	Environment      string
	Domain           string
	TemplateDatabase string // Name of the template database to clone
	AdminUser        string // Admin username for template database
	AdminPassword    string // Password for the admin user in template database
	DefaultDBMode    string // Default database mode: "create" or "clone"
	RateLimit        rate.Limit
	BurstLimit       int
	LogLevel         string
	TimeoutSeconds   int // HTTP client timeout in seconds
}

type Country struct {
	ID   int    `json:"id" validate:"required"`
	Code string `json:"code" validate:"required,len=2"`
	Name string `json:"name"`
}

// SignupRequest represents the signup form data
type SignupRequest struct {
	Username    string  `json:"username" validate:"required,min=3,max=20,alphanum"`
	Email       string  `json:"email" validate:"required,email"`
	Password    string  `json:"password" validate:"required,min=8"`
	FirstName   string  `json:"firstName" validate:"required,min=2"`
	LastName    string  `json:"lastName" validate:"required,min=2"`
	Phone       string  `json:"phone"`
	CompanyName string  `json:"companyName" validate:"required,min=2"`
	Industry    string  `json:"industry"`
	CompanySize string  `json:"companySize"`
	Country     Country `json:"country" validate:"required,dive"`
	DbMode      string  `json:"dbMode,omitempty" validate:"omitempty,oneof=create clone"`
	Terms       bool    `json:"terms" validate:"required"`
}

// SignupResponse represents the API response for signup
type SignupResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    *SignupData `json:"data,omitempty"`
}

// SignupData contains the signup result data
type SignupData struct {
	InstanceURL string `json:"instanceUrl"`
	Email       string `json:"email"`
	Database    string `json:"database"`
}

// DatabaseInfo represents database information
type DatabaseInfo struct {
	Name string `json:"name"`
}
