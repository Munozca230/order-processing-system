package models

import (
	"time"
)

// Address represents a customer address
type Address struct {
	Street     string `json:"street,omitempty" bson:"street,omitempty"`
	City       string `json:"city,omitempty" bson:"city,omitempty"`
	PostalCode string `json:"postalCode,omitempty" bson:"postalCode,omitempty"`
	Country    string `json:"country,omitempty" bson:"country,omitempty"`
}

// Preferences represents customer preferences
type Preferences struct {
	Newsletter    bool `json:"newsletter" bson:"newsletter"`
	Notifications bool `json:"notifications" bson:"notifications"`
}

// Customer represents a customer in the system
type Customer struct {
	CustomerID       string       `json:"customerId" bson:"customerId" validate:"required"`
	Name             string       `json:"name" bson:"name" validate:"required,min=1,max=255"`
	Email            string       `json:"email,omitempty" bson:"email,omitempty"`
	Phone            string       `json:"phone,omitempty" bson:"phone,omitempty"`
	Address          Address      `json:"address,omitempty" bson:"address,omitempty"`
	Active           bool         `json:"active" bson:"active"`
	CustomerTier     string       `json:"customerTier,omitempty" bson:"customerTier,omitempty"`
	Preferences      Preferences  `json:"preferences,omitempty" bson:"preferences,omitempty"`
	RegistrationDate *time.Time   `json:"registrationDate,omitempty" bson:"registrationDate,omitempty"`
	LastLogin        *time.Time   `json:"lastLogin,omitempty" bson:"lastLogin,omitempty"`
	LoyaltyPoints    int          `json:"loyaltyPoints,omitempty" bson:"loyaltyPoints,omitempty"`
	CreatedAt        time.Time    `json:"createdAt" bson:"createdAt"`
	UpdatedAt        time.Time    `json:"updatedAt" bson:"updatedAt"`
}

// CustomerSummary represents a simplified customer view
type CustomerSummary struct {
	CustomerID string `json:"customerId"`
	Name       string `json:"name"`
	Active     bool   `json:"active"`
}

// CustomerResponse represents the response for customer listings
type CustomerResponse struct {
	Customers []CustomerSummary `json:"customers"`
	Total     int               `json:"total"`
	Active    int               `json:"active_count,omitempty"`
	Inactive  int               `json:"inactive_count,omitempty"`
	Page      int               `json:"page,omitempty"`
	PageSize  int               `json:"pageSize,omitempty"`
}

// ErrorResponse represents an API error response
type ErrorResponse struct {
	Error     string                 `json:"error"`
	Message   string                 `json:"message"`
	Details   map[string]interface{} `json:"details,omitempty"`
	RequestID string                 `json:"requestId,omitempty"`
	Timestamp time.Time              `json:"timestamp"`
}

// HealthResponse represents the health check response
type HealthResponse struct {
	Status         string            `json:"status"`
	Service        string            `json:"service"`
	Version        string            `json:"version"`
	Timestamp      time.Time         `json:"timestamp"`
	Uptime         string            `json:"uptime"`
	Environment    string            `json:"environment"`
	Metrics        map[string]int    `json:"metrics"`
	Dependencies   map[string]string `json:"dependencies,omitempty"`
	TotalCustomers int               `json:"total_customers"`
	ActiveCustomers int              `json:"active_customers"`
}