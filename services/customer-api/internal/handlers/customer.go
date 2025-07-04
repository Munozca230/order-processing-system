package handlers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/customer-api-v2/internal/models"
	"github.com/customer-api-v2/internal/repository"
	"github.com/customer-api-v2/internal/services"
	"github.com/sirupsen/logrus"
)

// CustomerHandler handles HTTP requests for customers
type CustomerHandler struct {
	service *services.CustomerService
	logger  *logrus.Logger
}

// NewCustomerHandler creates a new customer handler
func NewCustomerHandler(service *services.CustomerService, logger *logrus.Logger) *CustomerHandler {
	return &CustomerHandler{
		service: service,
		logger:  logger,
	}
}

// GetCustomer handles GET /customers/:id
func (h *CustomerHandler) GetCustomer(c echo.Context) error {
	customerID := c.Param("id")
	
	if customerID == "" {
		return h.errorResponse(c, http.StatusBadRequest, "missing_parameter", "Customer ID is required")
	}
	
	ctx := c.Request().Context()
	customer, err := h.service.GetCustomer(ctx, customerID)
	if err != nil {
		// Check if it's a not found error
		if err.Error() == "customer with ID "+customerID+" not found" {
			return h.errorResponse(c, http.StatusNotFound, "customer_not_found", err.Error())
		}
		
		// Check if it's an inactive customer error
		if err.Error() == "customer "+customerID+" is not active" {
			return h.errorResponse(c, http.StatusGone, "customer_inactive", err.Error())
		}
		
		// Other errors are internal server errors
		return h.errorResponse(c, http.StatusInternalServerError, "internal_error", "Failed to retrieve customer")
	}
	
	return c.JSON(http.StatusOK, customer)
}

// GetCustomers handles GET /customers
func (h *CustomerHandler) GetCustomers(c echo.Context) error {
	// Parse query parameters
	filters := repository.CustomerFilters{}
	
	if activeStr := c.QueryParam("active"); activeStr != "" {
		if active, err := strconv.ParseBool(activeStr); err == nil {
			filters.Active = &active
		}
	}
	
	if email := c.QueryParam("email"); email != "" {
		filters.Email = email
	}
	
	// Parse pagination parameters
	if pageStr := c.QueryParam("page"); pageStr != "" {
		if page, err := strconv.Atoi(pageStr); err == nil && page >= 0 {
			filters.Page = page
		}
	}
	
	if pageSizeStr := c.QueryParam("page_size"); pageSizeStr != "" {
		if pageSize, err := strconv.Atoi(pageSizeStr); err == nil && pageSize > 0 && pageSize <= 100 {
			filters.PageSize = pageSize
		}
	}
	
	ctx := c.Request().Context()
	response, err := h.service.GetCustomers(ctx, filters)
	if err != nil {
		return h.errorResponse(c, http.StatusInternalServerError, "internal_error", "Failed to retrieve customers")
	}
	
	return c.JSON(http.StatusOK, response)
}

// GetActiveCustomers handles GET /customers/active
func (h *CustomerHandler) GetActiveCustomers(c echo.Context) error {
	// Parse query parameters
	filters := repository.CustomerFilters{}
	
	if email := c.QueryParam("email"); email != "" {
		filters.Email = email
	}
	
	// Parse pagination parameters
	if pageStr := c.QueryParam("page"); pageStr != "" {
		if page, err := strconv.Atoi(pageStr); err == nil && page >= 0 {
			filters.Page = page
		}
	}
	
	if pageSizeStr := c.QueryParam("page_size"); pageSizeStr != "" {
		if pageSize, err := strconv.Atoi(pageSizeStr); err == nil && pageSize > 0 && pageSize <= 100 {
			filters.PageSize = pageSize
		}
	}
	
	ctx := c.Request().Context()
	response, err := h.service.GetActiveCustomers(ctx, filters)
	if err != nil {
		return h.errorResponse(c, http.StatusInternalServerError, "internal_error", "Failed to retrieve active customers")
	}
	
	return c.JSON(http.StatusOK, response)
}

// CreateCustomer handles POST /customers
func (h *CustomerHandler) CreateCustomer(c echo.Context) error {
	var customer models.Customer
	
	if err := c.Bind(&customer); err != nil {
		return h.errorResponse(c, http.StatusBadRequest, "invalid_json", "Invalid JSON format")
	}
	
	ctx := c.Request().Context()
	if err := h.service.CreateCustomer(ctx, &customer); err != nil {
		if err.Error() == "customer "+customer.CustomerID+" already exists" {
			return h.errorResponse(c, http.StatusConflict, "customer_exists", err.Error())
		}
		
		// Check if it's a validation error
		if len(err.Error()) > 17 && err.Error()[:17] == "validation failed" {
			return h.errorResponse(c, http.StatusBadRequest, "validation_error", err.Error())
		}
		
		return h.errorResponse(c, http.StatusInternalServerError, "internal_error", "Failed to create customer")
	}
	
	return c.JSON(http.StatusCreated, map[string]interface{}{
		"message":    "Customer created successfully",
		"customerId": customer.CustomerID,
	})
}

// GetHealth handles GET /health
func (h *CustomerHandler) GetHealth(c echo.Context) error {
	ctx := c.Request().Context()
	health, err := h.service.GetHealthStatus(ctx)
	if err != nil {
		return h.errorResponse(c, http.StatusInternalServerError, "health_check_failed", "Health check failed")
	}
	
	// Return appropriate HTTP status based on health
	status := http.StatusOK
	if health.Status == "degraded" {
		status = http.StatusPartialContent
	} else if health.Status != "healthy" {
		status = http.StatusServiceUnavailable
	}
	
	return c.JSON(status, health)
}

// GetMetrics handles GET /metrics (simple metrics endpoint)
func (h *CustomerHandler) GetMetrics(c echo.Context) error {
	metrics := h.service.GetMetrics()
	return c.JSON(http.StatusOK, metrics)
}

// errorResponse creates a standardized error response
func (h *CustomerHandler) errorResponse(c echo.Context, status int, errorCode, message string) error {
	requestID := ""
	if id := c.Get("requestId"); id != nil {
		requestID = id.(string)
	}
	
	errorResp := models.ErrorResponse{
		Error:     errorCode,
		Message:   message,
		RequestID: requestID,
		Timestamp: time.Now(),
	}
	
	// Add additional context for certain errors
	switch status {
	case http.StatusBadRequest:
		errorResp.Details = map[string]interface{}{
			"hint": "Check your request parameters and try again",
		}
	case http.StatusNotFound:
		errorResp.Details = map[string]interface{}{
			"hint": "The requested customer was not found",
		}
	case http.StatusGone:
		errorResp.Details = map[string]interface{}{
			"hint": "The requested customer is inactive",
		}
	case http.StatusInternalServerError:
		errorResp.Details = map[string]interface{}{
			"hint": "An internal error occurred. Please try again later",
		}
	}
	
	return c.JSON(status, errorResp)
}