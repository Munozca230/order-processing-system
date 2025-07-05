package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/customer-api-v2/configs"
	"github.com/customer-api-v2/internal/handlers"
	custommiddleware "github.com/customer-api-v2/internal/middleware"
	"github.com/customer-api-v2/internal/repository"
	"github.com/customer-api-v2/internal/services"
	"github.com/sirupsen/logrus"
)

func main() {
	// Load configuration
	config := configs.LoadConfig()
	
	// Setup logger
	logger := setupLogger(config)
	
	logger.WithFields(logrus.Fields{
		"version":     config.Server.Version,
		"environment": config.Server.Environment,
		"port":        config.Server.Port,
	}).Info("🚀 Starting Customer API")
	
	// Initialize dependencies
	var customerRepo repository.CustomerRepository
	var err error
	
	if config.Database.Type == "mongodb" && config.Database.URL != "" {
		logger.Info("🔌 Connecting to MongoDB...")
		customerRepo, err = repository.NewMongoCustomerRepository(config.Database.URL)
		if err != nil {
			logger.WithError(err).Warn("⚠️ Failed to connect to MongoDB, falling back to memory repository")
			customerRepo = repository.NewMemoryCustomerRepository()
		} else {
			logger.Info("✅ Connected to MongoDB successfully")
		}
	} else {
		logger.Info("💾 Using in-memory repository")
		customerRepo = repository.NewMemoryCustomerRepository()
	}
	
	customerService := services.NewCustomerService(customerRepo, config, logger)
	customerHandler := handlers.NewCustomerHandler(customerService, logger)
	
	// Setup Echo server
	e := echo.New()
	
	// Configure Echo
	e.HideBanner = true
	e.Debug = config.Server.Environment == "development"
	
	// Global middleware
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())
	e.Use(custommiddleware.RequestIDMiddleware())
	
	if config.Logging.RequestLog {
		e.Use(custommiddleware.StructuredLoggingMiddleware(logger))
	}
	
	if config.Features.EnableMetrics {
		e.Use(custommiddleware.MetricsMiddleware())
	}
	
	e.Use(custommiddleware.ErrorHandlingMiddleware(logger))
	
	// Setup routes
	setupRoutes(e, customerHandler)
	
	// Setup server with timeouts
	server := &http.Server{
		Addr:         fmt.Sprintf("%s:%s", config.Server.Host, config.Server.Port),
		ReadTimeout:  config.Server.ReadTimeout,
		WriteTimeout: config.Server.WriteTimeout,
	}
	
	// Start server in goroutine
	go func() {
		logger.WithFields(logrus.Fields{
			"address": server.Addr,
			"env":     config.Server.Environment,
		}).Info("🌐 Server starting")
		
		if err := e.StartServer(server); err != nil && err != http.ErrServerClosed {
			logger.WithError(err).Fatal("💥 Failed to start server")
		}
	}()
	
	// Wait for interrupt signal for graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit
	
	logger.Info("🛑 Shutting down server...")
	
	// Graceful shutdown with timeout
	ctx, cancel := context.WithTimeout(context.Background(), config.Server.ShutdownTimeout)
	defer cancel()
	
	if err := e.Shutdown(ctx); err != nil {
		logger.WithError(err).Error("💥 Server forced to shutdown")
	} else {
		logger.Info("✅ Server shutdown completed")
	}
}

// setupLogger configures the application logger
func setupLogger(config *configs.Config) *logrus.Logger {
	logger := logrus.New()
	
	// Set log level
	level, err := logrus.ParseLevel(config.Logging.Level)
	if err != nil {
		level = logrus.InfoLevel
	}
	logger.SetLevel(level)
	
	// Set log format
	if config.Logging.Format == "json" {
		logger.SetFormatter(&logrus.JSONFormatter{
			TimestampFormat: time.RFC3339,
			FieldMap: logrus.FieldMap{
				logrus.FieldKeyTime:  "timestamp",
				logrus.FieldKeyLevel: "level",
				logrus.FieldKeyMsg:   "message",
			},
		})
	} else {
		logger.SetFormatter(&logrus.TextFormatter{
			FullTimestamp:   true,
			TimestampFormat: time.RFC3339,
		})
	}
	
	// Set output
	if config.Logging.Output == "file" {
		// In a real application, you'd configure file output here
		// For now, we'll stick with stdout
		logger.SetOutput(os.Stdout)
	} else {
		logger.SetOutput(os.Stdout)
	}
	
	return logger
}

// setupRoutes configures all API routes
func setupRoutes(e *echo.Echo, customerHandler *handlers.CustomerHandler) {
	// API v1 routes
	v1 := e.Group("/api/v1")
	{
		// Customer routes
		v1.GET("/customers", customerHandler.GetCustomers)
		v1.GET("/customers/active", customerHandler.GetActiveCustomers)
		v1.GET("/customers/:id", customerHandler.GetCustomer)
		v1.POST("/customers", customerHandler.CreateCustomer)
	}
	
	// Legacy routes for backward compatibility
	e.GET("/customers", customerHandler.GetCustomers)
	e.GET("/customers/active", customerHandler.GetActiveCustomers)
	e.GET("/customers/:id", customerHandler.GetCustomer)
	e.POST("/customers", customerHandler.CreateCustomer)
	
	// Health and monitoring routes (support both GET and HEAD for Docker health checks)
	e.GET("/health", customerHandler.GetHealth)
	e.HEAD("/health", customerHandler.GetHealth)
	e.GET("/metrics", customerHandler.GetMetrics)
	
	// Root endpoint
	e.GET("/", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]interface{}{
			"service":     "customer-api",
			"version":     "2.0.0",
			"environment": os.Getenv("ENVIRONMENT"),
			"timestamp":   time.Now(),
			"endpoints": map[string]string{
				"health":           "/health",
				"metrics":          "/metrics",
				"customers":        "/customers",
				"active_customers": "/customers/active",
				"api_v1":           "/api/v1",
			},
		})
	})
}