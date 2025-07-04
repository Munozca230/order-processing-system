package handlers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/product-api-v2/internal/models"
	"github.com/product-api-v2/internal/repository"
	"github.com/product-api-v2/internal/services"
	"github.com/sirupsen/logrus"
)

// ProductHandler handles HTTP requests for products
type ProductHandler struct {
	service *services.ProductService
	logger  *logrus.Logger
}

// NewProductHandler creates a new product handler
func NewProductHandler(service *services.ProductService, logger *logrus.Logger) *ProductHandler {
	return &ProductHandler{
		service: service,
		logger:  logger,
	}
}

// GetProduct handles GET /products/:id
func (h *ProductHandler) GetProduct(c echo.Context) error {
	productID := c.Param("id")
	
	if productID == "" {
		return h.errorResponse(c, http.StatusBadRequest, "missing_parameter", "Product ID is required")
	}
	
	ctx := c.Request().Context()
	product, err := h.service.GetProduct(ctx, productID)
	if err != nil {
		// Check if it's a not found error
		if err.Error() == "product with ID "+productID+" not found" {
			return h.errorResponse(c, http.StatusNotFound, "product_not_found", err.Error())
		}
		
		// Check if it's an availability error
		if err.Error() == "product "+productID+" is not available" {
			return h.errorResponse(c, http.StatusGone, "product_unavailable", err.Error())
		}
		
		// Other errors are internal server errors
		return h.errorResponse(c, http.StatusInternalServerError, "internal_error", "Failed to retrieve product")
	}
	
	return c.JSON(http.StatusOK, product)
}

// GetProducts handles GET /products
func (h *ProductHandler) GetProducts(c echo.Context) error {
	// Parse query parameters
	filters := repository.ProductFilters{}
	
	if category := c.QueryParam("category"); category != "" {
		filters.Category = category
	}
	
	if activeStr := c.QueryParam("active"); activeStr != "" {
		if active, err := strconv.ParseBool(activeStr); err == nil {
			filters.Active = &active
		}
	}
	
	if minPriceStr := c.QueryParam("min_price"); minPriceStr != "" {
		if minPrice, err := strconv.ParseFloat(minPriceStr, 64); err == nil {
			filters.MinPrice = &minPrice
		}
	}
	
	if maxPriceStr := c.QueryParam("max_price"); maxPriceStr != "" {
		if maxPrice, err := strconv.ParseFloat(maxPriceStr, 64); err == nil {
			filters.MaxPrice = &maxPrice
		}
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
	response, err := h.service.GetProducts(ctx, filters)
	if err != nil {
		return h.errorResponse(c, http.StatusInternalServerError, "internal_error", "Failed to retrieve products")
	}
	
	return c.JSON(http.StatusOK, response)
}

// CreateProduct handles POST /products
func (h *ProductHandler) CreateProduct(c echo.Context) error {
	var product models.Product
	
	if err := c.Bind(&product); err != nil {
		return h.errorResponse(c, http.StatusBadRequest, "invalid_json", "Invalid JSON format")
	}
	
	ctx := c.Request().Context()
	if err := h.service.CreateProduct(ctx, &product); err != nil {
		if err.Error() == "product "+product.ProductID+" already exists" {
			return h.errorResponse(c, http.StatusConflict, "product_exists", err.Error())
		}
		
		// Check if it's a validation error
		if err.Error()[:17] == "validation failed" {
			return h.errorResponse(c, http.StatusBadRequest, "validation_error", err.Error())
		}
		
		return h.errorResponse(c, http.StatusInternalServerError, "internal_error", "Failed to create product")
	}
	
	return c.JSON(http.StatusCreated, map[string]interface{}{
		"message":   "Product created successfully",
		"productId": product.ProductID,
	})
}

// GetHealth handles GET /health
func (h *ProductHandler) GetHealth(c echo.Context) error {
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
func (h *ProductHandler) GetMetrics(c echo.Context) error {
	metrics := h.service.GetMetrics()
	return c.JSON(http.StatusOK, metrics)
}

// errorResponse creates a standardized error response
func (h *ProductHandler) errorResponse(c echo.Context, status int, errorCode, message string) error {
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
			"hint": "The requested resource was not found",
		}
	case http.StatusInternalServerError:
		errorResp.Details = map[string]interface{}{
			"hint": "An internal error occurred. Please try again later",
		}
	}
	
	return c.JSON(status, errorResp)
}