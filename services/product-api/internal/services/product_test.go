package services

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/product-api-v2/configs"
	"github.com/product-api-v2/internal/models"
	"github.com/product-api-v2/internal/repository"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockProductRepository is a mock implementation of ProductRepository
type MockProductRepository struct {
	mock.Mock
}

func (m *MockProductRepository) GetByID(ctx context.Context, productID string) (*models.Product, error) {
	args := m.Called(ctx, productID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Product), args.Error(1)
}

func (m *MockProductRepository) GetAll(ctx context.Context, filters repository.ProductFilters) ([]*models.Product, error) {
	args := m.Called(ctx, filters)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*models.Product), args.Error(1)
}

func (m *MockProductRepository) Create(ctx context.Context, product *models.Product) error {
	args := m.Called(ctx, product)
	return args.Error(0)
}

func (m *MockProductRepository) Update(ctx context.Context, product *models.Product) error {
	args := m.Called(ctx, product)
	return args.Error(0)
}

func (m *MockProductRepository) Delete(ctx context.Context, productID string) error {
	args := m.Called(ctx, productID)
	return args.Error(0)
}

func (m *MockProductRepository) Count(ctx context.Context, filters repository.ProductFilters) (int, error) {
	args := m.Called(ctx, filters)
	return args.Int(0), args.Error(1)
}

func (m *MockProductRepository) HealthCheck(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

// Helper function to create a test ProductService
func createTestProductService(repo repository.ProductRepository) *ProductService {
	config := &configs.Config{
		Server: configs.ServerConfig{
			Version:     "test-v1.0.0",
			Environment: "test",
		},
		Features: configs.FeatureFlags{
			SimulateLatency: false,
			SimulateErrors:  false,
			MaxLatencyMs:    100,
			ErrorRate:       0.1,
		},
	}
	
	logger := logrus.New()
	logger.SetLevel(logrus.ErrorLevel) // Reduce log noise in tests
	
	return NewProductService(repo, config, logger)
}

// Helper function to create a test product
func createTestProduct() *models.Product {
	return &models.Product{
		ProductID:   "test-product-1",
		Name:        "Test Product",
		Description: "A test product",
		Price:       99.99,
		Category:    "electronics",
		Stock:       10,
		Active:      true,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
}

func TestNewProductService(t *testing.T) {
	mockRepo := &MockProductRepository{}
	service := createTestProductService(mockRepo)
	
	assert.NotNil(t, service)
	assert.Equal(t, mockRepo, service.repo)
	assert.NotNil(t, service.config)
	assert.NotNil(t, service.logger)
	assert.False(t, service.startTime.IsZero())
}

func TestProductService_GetProduct_Success(t *testing.T) {
	mockRepo := &MockProductRepository{}
	service := createTestProductService(mockRepo)
	ctx := context.Background()
	
	expectedProduct := createTestProduct()
	mockRepo.On("GetByID", ctx, "test-product-1").Return(expectedProduct, nil)
	
	product, err := service.GetProduct(ctx, "test-product-1")
	
	assert.NoError(t, err)
	assert.Equal(t, expectedProduct, product)
	assert.Equal(t, int64(1), service.requests)
	mockRepo.AssertExpectations(t)
}

func TestProductService_GetProduct_NotFound(t *testing.T) {
	mockRepo := &MockProductRepository{}
	service := createTestProductService(mockRepo)
	ctx := context.Background()
	
	mockRepo.On("GetByID", ctx, "nonexistent").Return(nil, repository.ErrProductNotFound)
	
	product, err := service.GetProduct(ctx, "nonexistent")
	
	assert.Error(t, err)
	assert.Nil(t, product)
	assert.Contains(t, err.Error(), "not found")
	assert.Equal(t, int64(1), service.requests)
	mockRepo.AssertExpectations(t)
}

func TestProductService_GetProduct_InactiveProduct(t *testing.T) {
	mockRepo := &MockProductRepository{}
	service := createTestProductService(mockRepo)
	ctx := context.Background()
	
	inactiveProduct := createTestProduct()
	inactiveProduct.Active = false
	
	mockRepo.On("GetByID", ctx, "inactive-product").Return(inactiveProduct, nil)
	
	product, err := service.GetProduct(ctx, "inactive-product")
	
	assert.Error(t, err)
	assert.Nil(t, product)
	assert.Contains(t, err.Error(), "not available")
	mockRepo.AssertExpectations(t)
}

func TestProductService_GetProduct_ErrorCase(t *testing.T) {
	mockRepo := &MockProductRepository{}
	service := createTestProductService(mockRepo)
	ctx := context.Background()
	
	// No mock needed - product-error has special handling that doesn't call repository
	
	product, err := service.GetProduct(ctx, "product-error")
	
	assert.Error(t, err)
	assert.Nil(t, product)
	assert.Contains(t, err.Error(), "this product always returns an error")
	assert.Equal(t, int64(1), service.errors)
	// No need to assert expectations since no repository call is made
}

func TestProductService_GetProduct_TestErrorProduct(t *testing.T) {
	mockRepo := &MockProductRepository{}
	service := createTestProductService(mockRepo)
	ctx := context.Background()
	
	// The service has special logic for "product-error" ID
	product, err := service.GetProduct(ctx, "product-error")
	
	assert.Error(t, err)
	assert.Nil(t, product)
	assert.Contains(t, err.Error(), "always returns an error")
	assert.Equal(t, int64(1), service.errors)
}

func TestProductService_GetProducts_Success(t *testing.T) {
	mockRepo := &MockProductRepository{}
	service := createTestProductService(mockRepo)
	ctx := context.Background()
	
	products := []*models.Product{createTestProduct()}
	filters := repository.ProductFilters{Category: "electronics"}
	
	mockRepo.On("GetAll", ctx, filters).Return(products, nil)
	mockRepo.On("Count", ctx, mock.AnythingOfType("repository.ProductFilters")).Return(1, nil)
	
	response, err := service.GetProducts(ctx, filters)
	
	assert.NoError(t, err)
	assert.NotNil(t, response)
	assert.Equal(t, 1, len(response.Products))
	assert.Equal(t, 1, response.Total)
	assert.Equal(t, "test-product-1", response.Products[0].ProductID)
	mockRepo.AssertExpectations(t)
}

func TestProductService_GetProducts_Error(t *testing.T) {
	mockRepo := &MockProductRepository{}
	service := createTestProductService(mockRepo)
	ctx := context.Background()
	
	filters := repository.ProductFilters{Category: "electronics"}
	mockRepo.On("GetAll", ctx, filters).Return(nil, errors.New("database error"))
	
	response, err := service.GetProducts(ctx, filters)
	
	assert.Error(t, err)
	assert.Nil(t, response)
	assert.Contains(t, err.Error(), "failed to retrieve products")
	assert.Equal(t, int64(1), service.errors)
	mockRepo.AssertExpectations(t)
}

func TestProductService_CreateProduct_Success(t *testing.T) {
	mockRepo := &MockProductRepository{}
	service := createTestProductService(mockRepo)
	ctx := context.Background()
	
	product := createTestProduct()
	mockRepo.On("Create", ctx, product).Return(nil)
	
	err := service.CreateProduct(ctx, product)
	
	assert.NoError(t, err)
	assert.Equal(t, int64(1), service.requests)
	mockRepo.AssertExpectations(t)
}

func TestProductService_CreateProduct_ValidationError(t *testing.T) {
	mockRepo := &MockProductRepository{}
	service := createTestProductService(mockRepo)
	ctx := context.Background()
	
	// Invalid product (missing required fields)
	invalidProduct := &models.Product{Price: -10}
	
	err := service.CreateProduct(ctx, invalidProduct)
	
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "validation failed")
	assert.Equal(t, int64(1), service.errors)
}

func TestProductService_CreateProduct_AlreadyExists(t *testing.T) {
	mockRepo := &MockProductRepository{}
	service := createTestProductService(mockRepo)
	ctx := context.Background()
	
	product := createTestProduct()
	mockRepo.On("Create", ctx, product).Return(repository.ErrProductExists)
	
	err := service.CreateProduct(ctx, product)
	
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "already exists")
	mockRepo.AssertExpectations(t)
}

func TestProductService_ValidateProduct_Success(t *testing.T) {
	service := createTestProductService(&MockProductRepository{})
	
	validProduct := createTestProduct()
	err := service.validateProduct(validProduct)
	
	assert.NoError(t, err)
}

func TestProductService_ValidateProduct_MissingProductID(t *testing.T) {
	service := createTestProductService(&MockProductRepository{})
	
	invalidProduct := createTestProduct()
	invalidProduct.ProductID = ""
	
	err := service.validateProduct(invalidProduct)
	
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "product ID is required")
}

func TestProductService_ValidateProduct_MissingName(t *testing.T) {
	service := createTestProductService(&MockProductRepository{})
	
	invalidProduct := createTestProduct()
	invalidProduct.Name = ""
	
	err := service.validateProduct(invalidProduct)
	
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "product name is required")
}

func TestProductService_ValidateProduct_InvalidPrice(t *testing.T) {
	service := createTestProductService(&MockProductRepository{})
	
	invalidProduct := createTestProduct()
	invalidProduct.Price = -10
	
	err := service.validateProduct(invalidProduct)
	
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "price must be greater than 0")
}

func TestProductService_ValidateProduct_NameTooLong(t *testing.T) {
	service := createTestProductService(&MockProductRepository{})
	
	invalidProduct := createTestProduct()
	invalidProduct.Name = string(make([]byte, 256)) // 256 characters
	
	err := service.validateProduct(invalidProduct)
	
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "cannot exceed 255 characters")
}

func TestProductService_GetHealthStatus_Healthy(t *testing.T) {
	mockRepo := &MockProductRepository{}
	service := createTestProductService(mockRepo)
	ctx := context.Background()
	
	mockRepo.On("HealthCheck", ctx).Return(nil)
	mockRepo.On("Count", ctx, repository.ProductFilters{}).Return(5, nil)
	
	health, err := service.GetHealthStatus(ctx)
	
	assert.NoError(t, err)
	assert.NotNil(t, health)
	assert.Equal(t, "healthy", health.Status)
	assert.Equal(t, "product-api", health.Service)
	assert.Equal(t, "test-v1.0.0", health.Version)
	assert.Equal(t, "test", health.Environment)
	assert.Contains(t, health.Dependencies, "repository")
	assert.Equal(t, "healthy", health.Dependencies["repository"])
	mockRepo.AssertExpectations(t)
}

func TestProductService_GetHealthStatus_UnhealthyRepository(t *testing.T) {
	mockRepo := &MockProductRepository{}
	service := createTestProductService(mockRepo)
	ctx := context.Background()
	
	mockRepo.On("HealthCheck", ctx).Return(errors.New("database connection failed"))
	mockRepo.On("Count", ctx, repository.ProductFilters{}).Return(0, nil)
	
	health, err := service.GetHealthStatus(ctx)
	
	assert.NoError(t, err)
	assert.NotNil(t, health)
	assert.Contains(t, health.Dependencies["repository"], "unhealthy")
	mockRepo.AssertExpectations(t)
}

func TestProductService_GetMetrics(t *testing.T) {
	service := createTestProductService(&MockProductRepository{})
	
	// Simulate some requests and errors
	service.requests = 100
	service.errors = 5
	
	metrics := service.GetMetrics()
	
	assert.NotNil(t, metrics)
	assert.Equal(t, "product-api", metrics["service"])
	assert.Equal(t, "test-v1.0.0", metrics["version"])
	assert.Equal(t, "test", metrics["environment"])
	assert.Equal(t, int64(100), metrics["total_requests"])
	assert.Equal(t, int64(5), metrics["total_errors"])
	assert.InDelta(t, 0.05, metrics["error_rate"], 0.001)
	assert.InDelta(t, 0.95, metrics["success_rate"], 0.001)
}

func TestProductService_GetMetrics_NoRequests(t *testing.T) {
	service := createTestProductService(&MockProductRepository{})
	
	metrics := service.GetMetrics()
	
	assert.NotNil(t, metrics)
	assert.Equal(t, int64(0), metrics["total_requests"])
	assert.Equal(t, int64(0), metrics["total_errors"])
	// error_rate and success_rate should not be present when there are no requests
	_, hasErrorRate := metrics["error_rate"]
	_, hasSuccessRate := metrics["success_rate"]
	assert.False(t, hasErrorRate)
	assert.False(t, hasSuccessRate)
}
