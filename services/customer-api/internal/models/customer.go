package models

import (
	"time"
)

// Customer represents a customer in the system
type Customer struct {
	CustomerID   string    `json:"customerId" validate:"required"`
	Name         string    `json:"name" validate:"required,min=1,max=255"`
	Email        string    `json:"email,omitempty"`
	Phone        string    `json:"phone,omitempty"`
	Address      string    `json:"address,omitempty"`
	Active       bool      `json:"active"`
	CustomerTier string    `json:"customerTier,omitempty"`
	CreatedAt    time.Time `json:"createdAt"`
	UpdatedAt    time.Time `json:"updatedAt"`
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