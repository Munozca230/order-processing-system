package repository

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/product-api-v2/internal/models"
)

var (
	ErrProductNotFound = errors.New("product not found")
	ErrProductExists   = errors.New("product already exists")
)

// ProductRepository defines the interface for product data operations
type ProductRepository interface {
	GetByID(ctx context.Context, productID string) (*models.Product, error)
	GetAll(ctx context.Context, filters ProductFilters) ([]*models.Product, error)
	Create(ctx context.Context, product *models.Product) error
	Update(ctx context.Context, product *models.Product) error
	Delete(ctx context.Context, productID string) error
	Count(ctx context.Context, filters ProductFilters) (int, error)
	HealthCheck(ctx context.Context) error
}

// ProductFilters holds filtering options for product queries
type ProductFilters struct {
	Category string
	Active   *bool
	MinPrice *float64
	MaxPrice *float64
	Page     int
	PageSize int
}

// MemoryProductRepository implements ProductRepository using in-memory storage
type MemoryProductRepository struct {
	products map[string]*models.Product
	mutex    sync.RWMutex
}

// NewMemoryProductRepository creates a new in-memory product repository
func NewMemoryProductRepository() *MemoryProductRepository {
	repo := &MemoryProductRepository{
		products: make(map[string]*models.Product),
		mutex:    sync.RWMutex{},
	}
	
	// Initialize with sample data
	repo.initializeSampleData()
	
	return repo
}

// GetByID retrieves a product by its ID
func (r *MemoryProductRepository) GetByID(ctx context.Context, productID string) (*models.Product, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()
	
	product, exists := r.products[productID]
	if !exists {
		return nil, ErrProductNotFound
	}
	
	// Return a copy to prevent external modifications
	productCopy := *product
	return &productCopy, nil
}

// GetAll retrieves all products with optional filtering
func (r *MemoryProductRepository) GetAll(ctx context.Context, filters ProductFilters) ([]*models.Product, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()
	
	var products []*models.Product
	
	for _, product := range r.products {
		if r.matchesFilters(product, filters) {
			productCopy := *product
			products = append(products, &productCopy)
		}
	}
	
	// Apply pagination
	if filters.PageSize > 0 {
		start := filters.Page * filters.PageSize
		end := start + filters.PageSize
		
		if start >= len(products) {
			return []*models.Product{}, nil
		}
		
		if end > len(products) {
			end = len(products)
		}
		
		products = products[start:end]
	}
	
	return products, nil
}

// Create adds a new product
func (r *MemoryProductRepository) Create(ctx context.Context, product *models.Product) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	
	if _, exists := r.products[product.ProductID]; exists {
		return ErrProductExists
	}
	
	now := time.Now()
	product.CreatedAt = now
	product.UpdatedAt = now
	
	productCopy := *product
	r.products[product.ProductID] = &productCopy
	
	return nil
}

// Update modifies an existing product
func (r *MemoryProductRepository) Update(ctx context.Context, product *models.Product) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	
	existing, exists := r.products[product.ProductID]
	if !exists {
		return ErrProductNotFound
	}
	
	product.CreatedAt = existing.CreatedAt
	product.UpdatedAt = time.Now()
	
	productCopy := *product
	r.products[product.ProductID] = &productCopy
	
	return nil
}

// Delete removes a product
func (r *MemoryProductRepository) Delete(ctx context.Context, productID string) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	
	if _, exists := r.products[productID]; !exists {
		return ErrProductNotFound
	}
	
	delete(r.products, productID)
	return nil
}

// Count returns the total number of products matching the filters
func (r *MemoryProductRepository) Count(ctx context.Context, filters ProductFilters) (int, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()
	
	count := 0
	for _, product := range r.products {
		if r.matchesFilters(product, filters) {
			count++
		}
	}
	
	return count, nil
}

// HealthCheck verifies the repository is working
func (r *MemoryProductRepository) HealthCheck(ctx context.Context) error {
	r.mutex.RLock()
	defer r.mutex.RUnlock()
	
	// Simple check - repository is healthy if we can access the map
	_ = len(r.products)
	return nil
}

// matchesFilters checks if a product matches the given filters
func (r *MemoryProductRepository) matchesFilters(product *models.Product, filters ProductFilters) bool {
	if filters.Category != "" && product.Category != filters.Category {
		return false
	}
	
	if filters.Active != nil && product.Active != *filters.Active {
		return false
	}
	
	if filters.MinPrice != nil && product.Price < *filters.MinPrice {
		return false
	}
	
	if filters.MaxPrice != nil && product.Price > *filters.MaxPrice {
		return false
	}
	
	return true
}

// initializeSampleData is deprecated - data now comes from MongoDB initialization scripts
func (r *MemoryProductRepository) initializeSampleData() {
	// No longer needed - MongoDB provides the expanded catalog (9 products) via init scripts
	// including: product-1 to product-9 with categories: laptops, peripherals, monitors, 
	// storage, components, furniture
	// This is kept for backward compatibility but does nothing
}