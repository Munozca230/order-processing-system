package models

import (
	"time"
)

// Product represents a product in the catalog
type Product struct {
	ProductID   string    `json:"productId" validate:"required"`
	Name        string    `json:"name" validate:"required,min=1,max=255"`
	Description string    `json:"description,omitempty"`
	Price       float64   `json:"price" validate:"required,gt=0"`
	Category    string    `json:"category,omitempty"`
	Stock       int       `json:"stock,omitempty"`
	Active      bool      `json:"active"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

// ProductSummary represents a simplified product view
type ProductSummary struct {
	ProductID string  `json:"productId"`
	Name      string  `json:"name"`
	Price     float64 `json:"price"`
}

// ProductCatalogResponse represents the response for product listings
type ProductCatalogResponse struct {
	Products []ProductSummary `json:"products"`
	Total    int              `json:"total"`
	Page     int              `json:"page,omitempty"`
	PageSize int              `json:"pageSize,omitempty"`
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
	Status      string            `json:"status"`
	Service     string            `json:"service"`
	Version     string            `json:"version"`
	Timestamp   time.Time         `json:"timestamp"`
	Uptime      string            `json:"uptime"`
	Environment string            `json:"environment"`
	Metrics     map[string]int    `json:"metrics"`
	Dependencies map[string]string `json:"dependencies,omitempty"`
}