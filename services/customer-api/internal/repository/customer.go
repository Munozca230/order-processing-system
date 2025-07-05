package repository

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/customer-api-v2/internal/models"
)

var (
	ErrCustomerNotFound = errors.New("customer not found")
	ErrCustomerExists   = errors.New("customer already exists")
)

// CustomerRepository defines the interface for customer data operations
type CustomerRepository interface {
	GetByID(ctx context.Context, customerID string) (*models.Customer, error)
	GetAll(ctx context.Context, filters CustomerFilters) ([]*models.Customer, error)
	Create(ctx context.Context, customer *models.Customer) error
	Update(ctx context.Context, customer *models.Customer) error
	Delete(ctx context.Context, customerID string) error
	Count(ctx context.Context, filters CustomerFilters) (int, error)
	HealthCheck(ctx context.Context) error
}

// CustomerFilters holds filtering options for customer queries
type CustomerFilters struct {
	Active       *bool
	Email        string
	CustomerTier string
	Page         int
	PageSize     int
}

// MemoryCustomerRepository implements CustomerRepository using in-memory storage
type MemoryCustomerRepository struct {
	customers map[string]*models.Customer
	mutex     sync.RWMutex
}

// NewMemoryCustomerRepository creates a new in-memory customer repository
func NewMemoryCustomerRepository() *MemoryCustomerRepository {
	repo := &MemoryCustomerRepository{
		customers: make(map[string]*models.Customer),
		mutex:     sync.RWMutex{},
	}
	
	// Initialize with sample data
	repo.initializeSampleData()
	
	return repo
}

// GetByID retrieves a customer by its ID
func (r *MemoryCustomerRepository) GetByID(ctx context.Context, customerID string) (*models.Customer, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()
	
	customer, exists := r.customers[customerID]
	if !exists {
		return nil, ErrCustomerNotFound
	}
	
	// Return a copy to prevent external modifications
	customerCopy := *customer
	return &customerCopy, nil
}

// GetAll retrieves all customers with optional filtering
func (r *MemoryCustomerRepository) GetAll(ctx context.Context, filters CustomerFilters) ([]*models.Customer, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()
	
	var customers []*models.Customer
	
	for _, customer := range r.customers {
		if r.matchesFilters(customer, filters) {
			customerCopy := *customer
			customers = append(customers, &customerCopy)
		}
	}
	
	// Apply pagination
	if filters.PageSize > 0 {
		start := filters.Page * filters.PageSize
		end := start + filters.PageSize
		
		if start >= len(customers) {
			return []*models.Customer{}, nil
		}
		
		if end > len(customers) {
			end = len(customers)
		}
		
		customers = customers[start:end]
	}
	
	return customers, nil
}

// Create adds a new customer
func (r *MemoryCustomerRepository) Create(ctx context.Context, customer *models.Customer) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	
	if _, exists := r.customers[customer.CustomerID]; exists {
		return ErrCustomerExists
	}
	
	now := time.Now()
	customer.CreatedAt = now
	customer.UpdatedAt = now
	
	customerCopy := *customer
	r.customers[customer.CustomerID] = &customerCopy
	
	return nil
}

// Update modifies an existing customer
func (r *MemoryCustomerRepository) Update(ctx context.Context, customer *models.Customer) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	
	existing, exists := r.customers[customer.CustomerID]
	if !exists {
		return ErrCustomerNotFound
	}
	
	customer.CreatedAt = existing.CreatedAt
	customer.UpdatedAt = time.Now()
	
	customerCopy := *customer
	r.customers[customer.CustomerID] = &customerCopy
	
	return nil
}

// Delete removes a customer
func (r *MemoryCustomerRepository) Delete(ctx context.Context, customerID string) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	
	if _, exists := r.customers[customerID]; !exists {
		return ErrCustomerNotFound
	}
	
	delete(r.customers, customerID)
	return nil
}

// Count returns the total number of customers matching the filters
func (r *MemoryCustomerRepository) Count(ctx context.Context, filters CustomerFilters) (int, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()
	
	count := 0
	for _, customer := range r.customers {
		if r.matchesFilters(customer, filters) {
			count++
		}
	}
	
	return count, nil
}

// HealthCheck verifies the repository is working
func (r *MemoryCustomerRepository) HealthCheck(ctx context.Context) error {
	r.mutex.RLock()
	defer r.mutex.RUnlock()
	
	// Simple check - repository is healthy if we can access the map
	_ = len(r.customers)
	return nil
}

// matchesFilters checks if a customer matches the given filters
func (r *MemoryCustomerRepository) matchesFilters(customer *models.Customer, filters CustomerFilters) bool {
	if filters.Active != nil && customer.Active != *filters.Active {
		return false
	}
	
	if filters.Email != "" && customer.Email != filters.Email {
		return false
	}
	
	if filters.CustomerTier != "" && customer.CustomerTier != filters.CustomerTier {
		return false
	}
	
	return true
}

// initializeSampleData is deprecated - data now comes from MongoDB initialization scripts
func (r *MemoryCustomerRepository) initializeSampleData() {
	// No longer needed - MongoDB provides the data via init scripts
	// This is kept for backward compatibility but does nothing
}