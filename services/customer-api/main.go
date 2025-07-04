package main

import (
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Customer struct {
	CustomerID string `json:"customerId"`
	Name       string `json:"name"`
	Active     bool   `json:"active"`
}

type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message"`
}

// In-memory customer database
var customerDatabase = map[string]Customer{
	"customer-1": {
		CustomerID: "customer-1",
		Name:       "Juan Pérez García",
		Active:     true,
	},
	"customer-2": {
		CustomerID: "customer-2",
		Name:       "María González López",
		Active:     true,
	},
	"customer-3": {
		CustomerID: "customer-3",
		Name:       "Carlos Rodríguez Silva",
		Active:     false, // Inactive customer for testing
	},
	"customer-inactive": {
		CustomerID: "customer-inactive",
		Name:       "Cliente Inactivo",
		Active:     false,
	},
	"customer-premium": {
		CustomerID: "customer-premium",
		Name:       "Ana Premium VIP",
		Active:     true,
	},
	"customer-error": {
		CustomerID: "customer-error",
		Name:       "Cliente que causa error",
		Active:     true,
	},
}

func getCustomer(c echo.Context) error {
	customerID := c.Param("id")
	
	// Simulate network latency for testing
	if latency := os.Getenv("SIMULATE_LATENCY"); latency == "true" {
		delay := time.Duration(rand.Intn(150)+50) * time.Millisecond
		c.Logger().Infof("Simulating latency: %v for customer %s", delay, customerID)
		time.Sleep(delay)
	}
	
	// Simulate random errors for testing retry mechanism
	if errorRate := os.Getenv("ERROR_RATE"); errorRate != "" {
		if rate, err := strconv.ParseFloat(errorRate, 64); err == nil {
			if rand.Float64() < rate {
				c.Logger().Errorf("Simulating error for customer %s", customerID)
				return c.JSON(http.StatusInternalServerError, ErrorResponse{
					Error:   "internal_server_error",
					Message: "Simulated error for testing",
				})
			}
		}
	}
	
	// Special case to always return error (for testing)
	if customerID == "customer-error" {
		return c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "customer_error",
			Message: "This customer always returns an error",
		})
	}
	
	// Look up customer in database
	customer, exists := customerDatabase[customerID]
	if !exists {
		c.Logger().Warnf("Customer not found: %s", customerID)
		return c.JSON(http.StatusNotFound, ErrorResponse{
			Error:   "customer_not_found",
			Message: "Customer with ID " + customerID + " not found",
		})
	}
	
	// Log customer status for monitoring
	if customer.Active {
		c.Logger().Infof("Successfully retrieved active customer: %s", customerID)
	} else {
		c.Logger().Warnf("Retrieved inactive customer: %s", customerID)
	}
	
	return c.JSON(http.StatusOK, customer)
}

func getAllCustomers(c echo.Context) error {
	c.Logger().Info("Getting all customers")
	customers := make([]Customer, 0, len(customerDatabase))
	activeCount := 0
	
	for _, customer := range customerDatabase {
		if customer.CustomerID != "customer-error" { // Don't include error customer in listing
			customers = append(customers, customer)
			if customer.Active {
				activeCount++
			}
		}
	}
	
	return c.JSON(http.StatusOK, map[string]interface{}{
		"customers":     customers,
		"total":         len(customers),
		"active_count":  activeCount,
		"inactive_count": len(customers) - activeCount,
	})
}

func getActiveCustomers(c echo.Context) error {
	c.Logger().Info("Getting active customers only")
	activeCustomers := make([]Customer, 0)
	
	for _, customer := range customerDatabase {
		if customer.Active && customer.CustomerID != "customer-error" {
			activeCustomers = append(activeCustomers, customer)
		}
	}
	
	return c.JSON(http.StatusOK, map[string]interface{}{
		"customers": activeCustomers,
		"total":     len(activeCustomers),
		"active":    true,
	})
}

func healthCheck(c echo.Context) error {
	activeCount := 0
	totalCount := 0
	
	for _, customer := range customerDatabase {
		if customer.CustomerID != "customer-error" {
			totalCount++
			if customer.Active {
				activeCount++
			}
		}
	}
	
	return c.JSON(http.StatusOK, map[string]interface{}{
		"status":         "healthy",
		"service":        "customer-api",
		"timestamp":      time.Now().UTC(),
		"version":        "1.0.0",
		"total_customers": totalCount,
		"active_customers": activeCount,
	})
}

func main() {
	e := echo.New()
	
	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())
	
	// Custom middleware for request ID
	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			c.Response().Header().Set("X-Request-ID", strconv.FormatInt(time.Now().UnixNano(), 36))
			return next(c)
		}
	})
	
	// Routes
	e.GET("/customers/:id", getCustomer)
	e.GET("/customers", getAllCustomers)
	e.GET("/customers/active", getActiveCustomers)
	e.GET("/health", healthCheck)
	
	// Seed random number generator
	rand.Seed(time.Now().UnixNano())
	
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	
	e.Logger.Infof("🚀 Customer API starting on port %s", port)
	e.Logger.Infof("👥 Loaded %d customers in database", len(customerDatabase))
	e.Logger.Fatal(e.Start(":" + port))
}