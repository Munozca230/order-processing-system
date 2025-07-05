package services

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"github.com/product-api-v2/configs"
	"github.com/product-api-v2/internal/models"
	"github.com/product-api-v2/internal/repository"
	"github.com/sirupsen/logrus"
)

// ProductService handles business logic for products
type ProductService struct {
	repo   repository.ProductRepository
	config *configs.Config
	logger *logrus.Logger
	
	// Metrics
	startTime time.Time
	requests  int64
	errors    int64
}

// NewProductService creates a new product service
func NewProductService(repo repository.ProductRepository, config *configs.Config, logger *logrus.Logger) *ProductService {
	return &ProductService{
		repo:      repo,
		config:    config,
		logger:    logger,
		startTime: time.Now(),
	}
}

// GetProduct retrieves a product by ID with business logic and error simulation
func (s *ProductService) GetProduct(ctx context.Context, productID string) (*models.Product, error) {
	s.requests++
	
	logger := s.logger.WithFields(logrus.Fields{
		"operation": "GetProduct",
		"productId": productID,
		"requestId": ctx.Value("requestId"),
	})
	
	logger.Info("🔍 Starting product lookup")
	
	// Simulate latency for testing if enabled
	if s.config.Features.SimulateLatency {
		delay := time.Duration(rand.Intn(s.config.Features.MaxLatencyMs)+50) * time.Millisecond
		logger.WithField("latency", delay).Info("💤 Simulating latency")
		time.Sleep(delay)
	}
	
	// Simulate errors for testing if enabled
	if s.config.Features.SimulateErrors && rand.Float64() < s.config.Features.ErrorRate {
		s.errors++
		logger.Error("💥 Simulating error for testing")
		return nil, fmt.Errorf("simulated error for testing")
	}
	
	// Special case to always return error (for testing)
	if productID == "product-error" {
		s.errors++
		logger.WithField("reason", "test_product").Error("💥 Test product error")
		return nil, fmt.Errorf("this product always returns an error")
	}
	
	// Get product from repository
	product, err := s.repo.GetByID(ctx, productID)
	if err != nil {
		if err == repository.ErrProductNotFound {
			logger.WithField("reason", "not_found").Warn("⚠️ Product not found")
			return nil, fmt.Errorf("product with ID %s not found", productID)
		}
		s.errors++
		logger.WithError(err).Error("💥 Repository error")
		return nil, fmt.Errorf("failed to retrieve product: %w", err)
	}
	
	// Business logic: Check if product is active
	if !product.Active {
		logger.WithField("reason", "inactive").Warn("⚠️ Product is inactive")
		return nil, fmt.Errorf("product %s is not available", productID)
	}
	
	logger.WithFields(logrus.Fields{
		"name":     product.Name,
		"price":    product.Price,
		"category": product.Category,
	}).Info("✅ Product retrieved successfully")
	
	return product, nil
}

// GetProducts retrieves all products with filtering and pagination
func (s *ProductService) GetProducts(ctx context.Context, filters repository.ProductFilters) (*models.ProductCatalogResponse, error) {
	s.requests++
	
	logger := s.logger.WithFields(logrus.Fields{
		"operation": "GetProducts",
		"filters":   fmt.Sprintf("%+v", filters),
		"requestId": ctx.Value("requestId"),
	})
	
	logger.Info("📋 Getting product catalog")
	
	// Get products from repository
	products, err := s.repo.GetAll(ctx, filters)
	if err != nil {
		s.errors++
		logger.WithError(err).Error("💥 Failed to get products")
		return nil, fmt.Errorf("failed to retrieve products: %w", err)
	}
	
	// Get total count for pagination
	totalCount, err := s.repo.Count(ctx, repository.ProductFilters{
		Category: filters.Category,
		Active:   filters.Active,
		MinPrice: filters.MinPrice,
		MaxPrice: filters.MaxPrice,
	})
	if err != nil {
		logger.WithError(err).Warn("⚠️ Failed to get total count")
		totalCount = len(products) // Fallback to current page count
	}
	
	// Convert to summary format (exclude sensitive data)
	summaries := make([]models.ProductSummary, 0, len(products))
	for _, product := range products {
		// Only include active products in listings
		if product.Active {
			summaries = append(summaries, models.ProductSummary{
				ProductID: product.ProductID,
				Name:      product.Name,
				Price:     product.Price,
			})
		}
	}
	
	response := &models.ProductCatalogResponse{
		Products: summaries,
		Total:    totalCount,
	}
	
	// Include pagination info if applicable
	if filters.PageSize > 0 {
		response.Page = filters.Page
		response.PageSize = filters.PageSize
	}
	
	logger.WithFields(logrus.Fields{
		"count": len(summaries),
		"total": totalCount,
	}).Info("✅ Product catalog retrieved successfully")
	
	return response, nil
}

// CreateProduct adds a new product with validation
func (s *ProductService) CreateProduct(ctx context.Context, product *models.Product) error {
	s.requests++
	
	logger := s.logger.WithFields(logrus.Fields{
		"operation": "CreateProduct",
		"productId": product.ProductID,
		"requestId": ctx.Value("requestId"),
	})
	
	logger.Info("➕ Creating new product")
	
	// Business validation
	if err := s.validateProduct(product); err != nil {
		s.errors++
		logger.WithError(err).Error("💥 Product validation failed")
		return fmt.Errorf("validation failed: %w", err)
	}
	
	// Create in repository
	if err := s.repo.Create(ctx, product); err != nil {
		if err == repository.ErrProductExists {
			logger.WithField("reason", "already_exists").Warn("⚠️ Product already exists")
			return fmt.Errorf("product %s already exists", product.ProductID)
		}
		s.errors++
		logger.WithError(err).Error("💥 Failed to create product")
		return fmt.Errorf("failed to create product: %w", err)
	}
	
	logger.WithFields(logrus.Fields{
		"name":  product.Name,
		"price": product.Price,
	}).Info("✅ Product created successfully")
	
	return nil
}

// GetHealthStatus returns the service health status
func (s *ProductService) GetHealthStatus(ctx context.Context) (*models.HealthResponse, error) {
	logger := s.logger.WithFields(logrus.Fields{
		"operation": "HealthCheck",
		"requestId": ctx.Value("requestId"),
	})
	
	// Check repository health
	dependencies := make(map[string]string)
	if err := s.repo.HealthCheck(ctx); err != nil {
		dependencies["repository"] = "unhealthy: " + err.Error()
		logger.WithError(err).Error("💥 Repository health check failed")
	} else {
		dependencies["repository"] = "healthy"
	}
	
	// Calculate metrics
	uptime := time.Since(s.startTime)
	
	metrics := map[string]int{
		"total_requests": int(s.requests),
		"total_errors":   int(s.errors),
		"products_count": s.getProductCount(ctx),
	}
	
	if s.requests > 0 {
		metrics["error_rate_percent"] = int((s.errors * 100) / s.requests)
	}
	
	status := "healthy"
	if s.errors > 0 && s.requests > 0 && (s.errors*100/s.requests) > 10 {
		status = "degraded"
	}
	
	response := &models.HealthResponse{
		Status:       status,
		Service:      "product-api",
		Version:      s.config.Server.Version,
		Timestamp:    time.Now(),
		Uptime:       uptime.String(),
		Environment:  s.config.Server.Environment,
		Metrics:      metrics,
		Dependencies: dependencies,
	}
	
	logger.WithFields(logrus.Fields{
		"status":   status,
		"uptime":   uptime,
		"requests": s.requests,
		"errors":   s.errors,
	}).Info("🏥 Health check completed")
	
	return response, nil
}

// validateProduct performs business validation on product data
func (s *ProductService) validateProduct(product *models.Product) error {
	if product.ProductID == "" {
		return fmt.Errorf("product ID is required")
	}
	
	if product.Name == "" {
		return fmt.Errorf("product name is required")
	}
	
	if product.Price <= 0 {
		return fmt.Errorf("product price must be greater than 0")
	}
	
	if len(product.Name) > 255 {
		return fmt.Errorf("product name cannot exceed 255 characters")
	}
	
	return nil
}

// getProductCount returns the total number of products for metrics
func (s *ProductService) getProductCount(ctx context.Context) int {
	count, err := s.repo.Count(ctx, repository.ProductFilters{})
	if err != nil {
		s.logger.WithError(err).Warn("⚠️ Failed to get product count for metrics")
		return 0
	}
	return count
}

// GetMetrics returns service metrics
func (s *ProductService) GetMetrics() map[string]interface{} {
	uptime := time.Since(s.startTime)
	
	metrics := map[string]interface{}{
		"service":         "product-api",
		"version":         s.config.Server.Version,
		"environment":     s.config.Server.Environment,
		"uptime_seconds":  int(uptime.Seconds()),
		"total_requests":  s.requests,
		"total_errors":    s.errors,
		"timestamp":       time.Now(),
	}
	
	if s.requests > 0 {
		metrics["error_rate"] = float64(s.errors) / float64(s.requests)
		metrics["success_rate"] = 1.0 - (float64(s.errors) / float64(s.requests))
	}
	
	return metrics
}