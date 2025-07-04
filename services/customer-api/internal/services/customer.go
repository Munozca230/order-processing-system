package services

import (
	"context"
	"fmt"
	"math/rand"
	"strconv"
	"time"

	"github.com/customer-api-v2/configs"
	"github.com/customer-api-v2/internal/models"
	"github.com/customer-api-v2/internal/repository"
	"github.com/sirupsen/logrus"
)

// CustomerService handles business logic for customers
type CustomerService struct {
	repo   repository.CustomerRepository
	config *configs.Config
	logger *logrus.Logger
	
	// Metrics
	startTime time.Time
	requests  int64
	errors    int64
}

// NewCustomerService creates a new customer service
func NewCustomerService(repo repository.CustomerRepository, config *configs.Config, logger *logrus.Logger) *CustomerService {
	return &CustomerService{
		repo:      repo,
		config:    config,
		logger:    logger,
		startTime: time.Now(),
	}
}

// GetCustomer retrieves a customer by ID with business logic and error simulation
func (s *CustomerService) GetCustomer(ctx context.Context, customerID string) (*models.Customer, error) {
	s.requests++
	
	logger := s.logger.WithFields(logrus.Fields{
		"operation":  "GetCustomer",
		"customerId": customerID,
		"requestId":  ctx.Value("requestId"),
	})
	
	logger.Info("🔍 Starting customer lookup")
	
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
	if customerID == "customer-error" {
		s.errors++
		logger.WithField("reason", "test_customer").Error("💥 Test customer error")
		return nil, fmt.Errorf("this customer always returns an error")
	}
	
	// Get customer from repository
	customer, err := s.repo.GetByID(ctx, customerID)
	if err != nil {
		if err == repository.ErrCustomerNotFound {
			logger.WithField("reason", "not_found").Warn("⚠️ Customer not found")
			return nil, fmt.Errorf("customer with ID %s not found", customerID)
		}
		s.errors++
		logger.WithError(err).Error("💥 Repository error")
		return nil, fmt.Errorf("failed to retrieve customer: %w", err)
	}
	
	// Business logic: Check if customer is active
	if !customer.Active {
		logger.WithField("reason", "inactive").Warn("⚠️ Customer is inactive")
		return nil, fmt.Errorf("customer %s is not active", customerID)
	}
	
	logger.WithFields(logrus.Fields{
		"name":  customer.Name,
		"email": customer.Email,
		"phone": customer.Phone,
	}).Info("✅ Customer retrieved successfully")
	
	return customer, nil
}

// GetCustomers retrieves all customers with filtering and pagination
func (s *CustomerService) GetCustomers(ctx context.Context, filters repository.CustomerFilters) (*models.CustomerResponse, error) {
	s.requests++
	
	logger := s.logger.WithFields(logrus.Fields{
		"operation": "GetCustomers",
		"filters":   fmt.Sprintf("%+v", filters),
		"requestId": ctx.Value("requestId"),
	})
	
	logger.Info("📋 Getting customer list")
	
	// Get customers from repository
	customers, err := s.repo.GetAll(ctx, filters)
	if err != nil {
		s.errors++
		logger.WithError(err).Error("💥 Failed to get customers")
		return nil, fmt.Errorf("failed to retrieve customers: %w", err)
	}
	
	// Get total count for pagination
	totalCount, err := s.repo.Count(ctx, repository.CustomerFilters{
		Active: filters.Active,
		Email:  filters.Email,
	})
	if err != nil {
		logger.WithError(err).Warn("⚠️ Failed to get total count")
		totalCount = len(customers) // Fallback to current page count
	}
	
	// Convert to summary format and count active/inactive
	summaries := make([]models.CustomerSummary, 0, len(customers))
	activeCount := 0
	inactiveCount := 0
	
	for _, customer := range customers {
		// Include all customers in listings but don't include error customer
		if customer.CustomerID != "customer-error" {
			summaries = append(summaries, models.CustomerSummary{
				CustomerID: customer.CustomerID,
				Name:       customer.Name,
				Active:     customer.Active,
			})
			
			if customer.Active {
				activeCount++
			} else {
				inactiveCount++
			}
		}
	}
	
	response := &models.CustomerResponse{
		Customers: summaries,
		Total:     totalCount,
		Active:    activeCount,
		Inactive:  inactiveCount,
	}
	
	// Include pagination info if applicable
	if filters.PageSize > 0 {
		response.Page = filters.Page
		response.PageSize = filters.PageSize
	}
	
	logger.WithFields(logrus.Fields{
		"count":         len(summaries),
		"total":         totalCount,
		"active_count":  activeCount,
		"inactive_count": inactiveCount,
	}).Info("✅ Customer list retrieved successfully")
	
	return response, nil
}

// GetActiveCustomers retrieves only active customers
func (s *CustomerService) GetActiveCustomers(ctx context.Context, filters repository.CustomerFilters) (*models.CustomerResponse, error) {
	s.requests++
	
	logger := s.logger.WithFields(logrus.Fields{
		"operation": "GetActiveCustomers",
		"requestId": ctx.Value("requestId"),
	})
	
	logger.Info("📋 Getting active customers only")
	
	// Force active filter
	activeFilter := true
	filters.Active = &activeFilter
	
	// Get customers from repository
	customers, err := s.repo.GetAll(ctx, filters)
	if err != nil {
		s.errors++
		logger.WithError(err).Error("💥 Failed to get active customers")
		return nil, fmt.Errorf("failed to retrieve active customers: %w", err)
	}
	
	// Convert to summary format
	summaries := make([]models.CustomerSummary, 0, len(customers))
	
	for _, customer := range customers {
		// Only include active customers and exclude error customer
		if customer.Active && customer.CustomerID != "customer-error" {
			summaries = append(summaries, models.CustomerSummary{
				CustomerID: customer.CustomerID,
				Name:       customer.Name,
				Active:     customer.Active,
			})
		}
	}
	
	response := &models.CustomerResponse{
		Customers: summaries,
		Total:     len(summaries),
		Active:    len(summaries),
		Inactive:  0,
	}
	
	logger.WithFields(logrus.Fields{
		"count": len(summaries),
	}).Info("✅ Active customers retrieved successfully")
	
	return response, nil
}

// CreateCustomer adds a new customer with validation
func (s *CustomerService) CreateCustomer(ctx context.Context, customer *models.Customer) error {
	s.requests++
	
	logger := s.logger.WithFields(logrus.Fields{
		"operation":  "CreateCustomer",
		"customerId": customer.CustomerID,
		"requestId":  ctx.Value("requestId"),
	})
	
	logger.Info("➕ Creating new customer")
	
	// Business validation
	if err := s.validateCustomer(customer); err != nil {
		s.errors++
		logger.WithError(err).Error("💥 Customer validation failed")
		return fmt.Errorf("validation failed: %w", err)
	}
	
	// Create in repository
	if err := s.repo.Create(ctx, customer); err != nil {
		if err == repository.ErrCustomerExists {
			logger.WithField("reason", "already_exists").Warn("⚠️ Customer already exists")
			return fmt.Errorf("customer %s already exists", customer.CustomerID)
		}
		s.errors++
		logger.WithError(err).Error("💥 Failed to create customer")
		return fmt.Errorf("failed to create customer: %w", err)
	}
	
	logger.WithFields(logrus.Fields{
		"name":  customer.Name,
		"email": customer.Email,
	}).Info("✅ Customer created successfully")
	
	return nil
}

// GetHealthStatus returns the service health status
func (s *CustomerService) GetHealthStatus(ctx context.Context) (*models.HealthResponse, error) {
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
		"customers_count": s.getCustomerCount(ctx),
	}
	
	if s.requests > 0 {
		metrics["error_rate_percent"] = int((s.errors * 100) / s.requests)
	}
	
	status := "healthy"
	if s.errors > 0 && s.requests > 0 && (s.errors*100/s.requests) > 10 {
		status = "degraded"
	}
	
	// Get customer counts for health response
	totalCustomers := s.getCustomerCount(ctx)
	activeCustomers := s.getActiveCustomerCount(ctx)
	
	response := &models.HealthResponse{
		Status:          status,
		Service:         "customer-api",
		Version:         s.config.Server.Version,
		Timestamp:       time.Now(),
		Uptime:          uptime.String(),
		Environment:     s.config.Server.Environment,
		Metrics:         metrics,
		Dependencies:    dependencies,
		TotalCustomers:  totalCustomers,
		ActiveCustomers: activeCustomers,
	}
	
	logger.WithFields(logrus.Fields{
		"status":           status,
		"uptime":           uptime,
		"requests":         s.requests,
		"errors":           s.errors,
		"total_customers":  totalCustomers,
		"active_customers": activeCustomers,
	}).Info("🏥 Health check completed")
	
	return response, nil
}

// validateCustomer performs business validation on customer data
func (s *CustomerService) validateCustomer(customer *models.Customer) error {
	if customer.CustomerID == "" {
		return fmt.Errorf("customer ID is required")
	}
	
	if customer.Name == "" {
		return fmt.Errorf("customer name is required")
	}
	
	if len(customer.Name) > 255 {
		return fmt.Errorf("customer name cannot exceed 255 characters")
	}
	
	// Optional email validation
	if customer.Email != "" && len(customer.Email) > 320 {
		return fmt.Errorf("customer email cannot exceed 320 characters")
	}
	
	return nil
}

// getCustomerCount returns the total number of customers for metrics
func (s *CustomerService) getCustomerCount(ctx context.Context) int {
	count, err := s.repo.Count(ctx, repository.CustomerFilters{})
	if err != nil {
		s.logger.WithError(err).Warn("⚠️ Failed to get customer count for metrics")
		return 0
	}
	return count
}

// getActiveCustomerCount returns the number of active customers
func (s *CustomerService) getActiveCustomerCount(ctx context.Context) int {
	activeFilter := true
	count, err := s.repo.Count(ctx, repository.CustomerFilters{Active: &activeFilter})
	if err != nil {
		s.logger.WithError(err).Warn("⚠️ Failed to get active customer count for metrics")
		return 0
	}
	return count
}

// GetMetrics returns service metrics
func (s *CustomerService) GetMetrics() map[string]interface{} {
	uptime := time.Since(s.startTime)
	
	metrics := map[string]interface{}{
		"service":         "customer-api",
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