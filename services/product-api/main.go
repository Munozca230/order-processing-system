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

type Product struct {
	ProductID string  `json:"productId"`
	Name      string  `json:"name"`
	Price     float64 `json:"price"`
}

type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message"`
}

// In-memory product catalog
var productCatalog = map[string]Product{
	"product-1": {
		ProductID: "product-1",
		Name:      "Laptop Gaming MSI",
		Price:     1299.99,
	},
	"product-2": {
		ProductID: "product-2",
		Name:      "Mouse Gamer Logitech",
		Price:     59.99,
	},
	"product-3": {
		ProductID: "product-3",
		Name:      "Teclado Mecánico RGB",
		Price:     129.99,
	},
	"product-4": {
		ProductID: "product-4",
		Name:      "Monitor 4K 27 pulgadas",
		Price:     399.99,
	},
	"product-5": {
		ProductID: "product-5",
		Name:      "Auriculares Gaming",
		Price:     89.99,
	},
	"product-error": {
		ProductID: "product-error",
		Name:      "Producto que causa error",
		Price:     999.99,
	},
}

func getProduct(c echo.Context) error {
	productID := c.Param("id")
	
	// Simulate network latency for testing
	if latency := os.Getenv("SIMULATE_LATENCY"); latency == "true" {
		delay := time.Duration(rand.Intn(200)+100) * time.Millisecond
		c.Logger().Infof("Simulating latency: %v for product %s", delay, productID)
		time.Sleep(delay)
	}
	
	// Simulate random errors for testing retry mechanism
	if errorRate := os.Getenv("ERROR_RATE"); errorRate != "" {
		if rate, err := strconv.ParseFloat(errorRate, 64); err == nil {
			if rand.Float64() < rate {
				c.Logger().Errorf("Simulating error for product %s", productID)
				return c.JSON(http.StatusInternalServerError, ErrorResponse{
					Error:   "internal_server_error",
					Message: "Simulated error for testing",
				})
			}
		}
	}
	
	// Special case to always return error (for testing)
	if productID == "product-error" {
		return c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "product_error",
			Message: "This product always returns an error",
		})
	}
	
	// Look up product in catalog
	product, exists := productCatalog[productID]
	if !exists {
		c.Logger().Warnf("Product not found: %s", productID)
		return c.JSON(http.StatusNotFound, ErrorResponse{
			Error:   "product_not_found",
			Message: "Product with ID " + productID + " not found",
		})
	}
	
	c.Logger().Infof("Successfully retrieved product: %s", productID)
	return c.JSON(http.StatusOK, product)
}

func getAllProducts(c echo.Context) error {
	c.Logger().Info("Getting all products")
	products := make([]Product, 0, len(productCatalog))
	
	for _, product := range productCatalog {
		if product.ProductID != "product-error" { // Don't include error product in listing
			products = append(products, product)
		}
	}
	
	return c.JSON(http.StatusOK, map[string]interface{}{
		"products": products,
		"total":    len(products),
	})
}

func healthCheck(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]interface{}{
		"status":    "healthy",
		"service":   "product-api",
		"timestamp": time.Now().UTC(),
		"version":   "1.0.0",
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
	e.GET("/products/:id", getProduct)
	e.GET("/products", getAllProducts)
	e.GET("/health", healthCheck)
	
	// Seed random number generator
	rand.Seed(time.Now().UnixNano())
	
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	
	e.Logger.Infof("🚀 Product API starting on port %s", port)
	e.Logger.Infof("📦 Loaded %d products in catalog", len(productCatalog))
	e.Logger.Fatal(e.Start(":" + port))
}