package services

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/customer-api-v2/configs"
	"github.com/customer-api-v2/internal/models"
	"github.com/customer-api-v2/internal/repository"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockCustomerRepository is a mock implementation of CustomerRepository
type MockCustomerRepository struct {
	mock.Mock
}

func (m *MockCustomerRepository) GetByID(ctx context.Context, customerID string) (*models.Customer, error) {
	args := m.Called(ctx, customerID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Customer), args.Error(1)
}

func (m *MockCustomerRepository) GetAll(ctx context.Context, filters repository.CustomerFilters) ([]*models.Customer, error) {
	args := m.Called(ctx, filters)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*models.Customer), args.Error(1)
}

func (m *MockCustomerRepository) Create(ctx context.Context, customer *models.Customer) error {
	args := m.Called(ctx, customer)
	return args.Error(0)
}

func (m *MockCustomerRepository) Update(ctx context.Context, customer *models.Customer) error {
	args := m.Called(ctx, customer)
	return args.Error(0)
}

func (m *MockCustomerRepository) Delete(ctx context.Context, customerID string) error {
	args := m.Called(ctx, customerID)
	return args.Error(0)
}

func (m *MockCustomerRepository) Count(ctx context.Context, filters repository.CustomerFilters) (int, error) {
	args := m.Called(ctx, filters)
	return args.Int(0), args.Error(1)
}

func (m *MockCustomerRepository) HealthCheck(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

// Helper function to create a test CustomerService
func createTestCustomerService(repo repository.CustomerRepository) *CustomerService {
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
	
	return NewCustomerService(repo, config, logger)
}

// Helper function to create a test customer
func createTestCustomer() *models.Customer {
	return &models.Customer{
		CustomerID: "test-customer-1",
		Name:       "John Doe",
		Email:      "john.doe@example.com",
		Phone:      "+1-555-123-4567",
		Address: models.Address{
			Street:     "123 Main St",
			City:       "Anytown",
			PostalCode: "12345",
			Country:    "USA",
		},
		Active:        true,
		CustomerTier: "premium",
		Preferences: models.Preferences{
			Newsletter:    true,
			Notifications: true,
		},
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

func TestNewCustomerService(t *testing.T) {
	mockRepo := &MockCustomerRepository{}
	service := createTestCustomerService(mockRepo)
	
	assert.NotNil(t, service)
	assert.Equal(t, mockRepo, service.repo)
	assert.NotNil(t, service.config)
	assert.NotNil(t, service.logger)
	assert.False(t, service.startTime.IsZero())
}

func TestCustomerService_GetCustomer_Success(t *testing.T) {
	mockRepo := &MockCustomerRepository{}
	service := createTestCustomerService(mockRepo)
	ctx := context.Background()
	
	expectedCustomer := createTestCustomer()
	mockRepo.On("GetByID", ctx, "test-customer-1").Return(expectedCustomer, nil)
	
	customer, err := service.GetCustomer(ctx, "test-customer-1")
	
	assert.NoError(t, err)
	assert.Equal(t, expectedCustomer, customer)
	assert.Equal(t, int64(1), service.requests)
	mockRepo.AssertExpectations(t)
}

func TestCustomerService_GetCustomer_NotFound(t *testing.T) {
	mockRepo := &MockCustomerRepository{}
	service := createTestCustomerService(mockRepo)
	ctx := context.Background()
	
	mockRepo.On("GetByID", ctx, "nonexistent").Return(nil, repository.ErrCustomerNotFound)
	
	customer, err := service.GetCustomer(ctx, "nonexistent")
	
	assert.Error(t, err)
	assert.Nil(t, customer)
	assert.Contains(t, err.Error(), "not found")
	assert.Equal(t, int64(1), service.requests)
	mockRepo.AssertExpectations(t)
}

func TestCustomerService_GetCustomer_InactiveCustomer(t *testing.T) {
	mockRepo := &MockCustomerRepository{}
	service := createTestCustomerService(mockRepo)
	ctx := context.Background()
	
	inactiveCustomer := createTestCustomer()
	inactiveCustomer.Active = false
	
	mockRepo.On("GetByID", ctx, "inactive-customer").Return(inactiveCustomer, nil)
	
	customer, err := service.GetCustomer(ctx, "inactive-customer")
	
	assert.Error(t, err)
	assert.Nil(t, customer)
	assert.Contains(t, err.Error(), "not active")
	mockRepo.AssertExpectations(t)
}

func TestCustomerService_GetCustomer_ErrorCase(t *testing.T) {
	mockRepo := &MockCustomerRepository{}
	service := createTestCustomerService(mockRepo)
	ctx := context.Background()
	
	mockRepo.On("GetByID", ctx, "customer-error").Return(nil, errors.New("database error"))
	
	customer, err := service.GetCustomer(ctx, "customer-error")
	
	assert.Error(t, err)
	assert.Nil(t, customer)
	assert.Contains(t, err.Error(), "failed to retrieve customer")
	assert.Equal(t, int64(1), service.errors)
	mockRepo.AssertExpectations(t)
}

func TestCustomerService_GetCustomer_TestErrorCustomer(t *testing.T) {
	mockRepo := &MockCustomerRepository{}
	service := createTestCustomerService(mockRepo)
	ctx := context.Background()
	
	// The service has special logic for "customer-error" ID
	customer, err := service.GetCustomer(ctx, "customer-error")
	
	assert.Error(t, err)
	assert.Nil(t, customer)
	assert.Contains(t, err.Error(), "always returns an error")
	assert.Equal(t, int64(1), service.errors)
}

func TestCustomerService_GetCustomers_Success(t *testing.T) {
	mockRepo := &MockCustomerRepository{}
	service := createTestCustomerService(mockRepo)
	ctx := context.Background()
	
	customers := []*models.Customer{createTestCustomer()}
	filters := repository.CustomerFilters{Email: "john.doe@example.com"}
	
	mockRepo.On("GetAll", ctx, filters).Return(customers, nil)
	mockRepo.On("Count", ctx, mock.AnythingOfType("repository.CustomerFilters")).Return(1, nil)
	
	response, err := service.GetCustomers(ctx, filters)
	
	assert.NoError(t, err)
	assert.NotNil(t, response)
	assert.Equal(t, 1, len(response.Customers))
	assert.Equal(t, 1, response.Total)
	assert.Equal(t, "test-customer-1", response.Customers[0].CustomerID)
	mockRepo.AssertExpectations(t)
}

func TestCustomerService_GetCustomers_Error(t *testing.T) {
	mockRepo := &MockCustomerRepository{}
	service := createTestCustomerService(mockRepo)
	ctx := context.Background()
	
	filters := repository.CustomerFilters{Email: "test@example.com"}
	mockRepo.On("GetAll", ctx, filters).Return(nil, errors.New("database error"))
	
	response, err := service.GetCustomers(ctx, filters)
	
	assert.Error(t, err)
	assert.Nil(t, response)
	assert.Contains(t, err.Error(), "failed to retrieve customers")
	assert.Equal(t, int64(1), service.errors)
	mockRepo.AssertExpectations(t)
}

func TestCustomerService_GetActiveCustomers_Success(t *testing.T) {
	mockRepo := &MockCustomerRepository{}
	service := createTestCustomerService(mockRepo)
	ctx := context.Background()
	
	activeCustomer := createTestCustomer()
	customers := []*models.Customer{activeCustomer}
	
	// For active customers, we need to set up both GetAll and Count calls
	activeFilter := true
	expectedFilters := repository.CustomerFilters{Active: &activeFilter}
	
	mockRepo.On("GetAll", ctx, mock.MatchedBy(func(f repository.CustomerFilters) bool {
		return f.Active != nil && *f.Active == true
	})).Return(customers, nil)
	
	mockRepo.On("Count", ctx, mock.MatchedBy(func(f repository.CustomerFilters) bool {
		return f.Active != nil && *f.Active == true
	})).Return(1, nil)
	
	response, err := service.GetActiveCustomers(ctx, expectedFilters)
	
	assert.NoError(t, err)
	assert.NotNil(t, response)
	assert.Equal(t, 1, len(response.Customers))
	assert.Equal(t, 1, response.Total)
	assert.Equal(t, "test-customer-1", response.Customers[0].CustomerID)
	mockRepo.AssertExpectations(t)
}

func TestCustomerService_CreateCustomer_Success(t *testing.T) {
	mockRepo := &MockCustomerRepository{}
	service := createTestCustomerService(mockRepo)
	ctx := context.Background()
	
	customer := createTestCustomer()
	mockRepo.On("Create", ctx, customer).Return(nil)
	
	err := service.CreateCustomer(ctx, customer)
	
	assert.NoError(t, err)
	assert.Equal(t, int64(1), service.requests)
	mockRepo.AssertExpectations(t)
}

func TestCustomerService_CreateCustomer_ValidationError(t *testing.T) {
	mockRepo := &MockCustomerRepository{}
	service := createTestCustomerService(mockRepo)
	ctx := context.Background()
	
	// Invalid customer (missing required fields)
	invalidCustomer := &models.Customer{CustomerID: ""}
	
	err := service.CreateCustomer(ctx, invalidCustomer)
	
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "validation failed")
	assert.Equal(t, int64(1), service.errors)
}

func TestCustomerService_CreateCustomer_AlreadyExists(t *testing.T) {
	mockRepo := &MockCustomerRepository{}
	service := createTestCustomerService(mockRepo)
	ctx := context.Background()
	
	customer := createTestCustomer()
	mockRepo.On("Create", ctx, customer).Return(repository.ErrCustomerExists)
	
	err := service.CreateCustomer(ctx, customer)
	
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "already exists")
	mockRepo.AssertExpectations(t)
}

func TestCustomerService_ValidateCustomer_Success(t *testing.T) {
	service := createTestCustomerService(&MockCustomerRepository{})
	
	validCustomer := createTestCustomer()
	err := service.validateCustomer(validCustomer)
	
	assert.NoError(t, err)
}

func TestCustomerService_ValidateCustomer_MissingCustomerID(t *testing.T) {
	service := createTestCustomerService(&MockCustomerRepository{})
	
	invalidCustomer := createTestCustomer()
	invalidCustomer.CustomerID = ""
	
	err := service.validateCustomer(invalidCustomer)
	
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "customer ID is required")
}

func TestCustomerService_ValidateCustomer_MissingName(t *testing.T) {
	service := createTestCustomerService(&MockCustomerRepository{})
	
	invalidCustomer := createTestCustomer()
	invalidCustomer.Name = ""
	
	err := service.validateCustomer(invalidCustomer)
	
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "customer name is required")
}

func TestCustomerService_ValidateCustomer_NameTooLong(t *testing.T) {
	service := createTestCustomerService(&MockCustomerRepository{})
	
	invalidCustomer := createTestCustomer()
	invalidCustomer.Name = string(make([]byte, 256)) // 256 characters
	
	err := service.validateCustomer(invalidCustomer)
	
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "cannot exceed 255 characters")
}

func TestCustomerService_ValidateCustomer_EmailTooLong(t *testing.T) {
	service := createTestCustomerService(&MockCustomerRepository{})
	
	invalidCustomer := createTestCustomer()
	invalidCustomer.Email = string(make([]byte, 321)) // 321 characters
	
	err := service.validateCustomer(invalidCustomer)
	
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "email cannot exceed 320 characters")
}

func TestCustomerService_GetHealthStatus_Healthy(t *testing.T) {
	mockRepo := &MockCustomerRepository{}
	service := createTestCustomerService(mockRepo)
	ctx := context.Background()
	
	mockRepo.On("HealthCheck", ctx).Return(nil)
	mockRepo.On("Count", ctx, repository.CustomerFilters{}).Return(10, nil)
	
	activeFilter := true
	mockRepo.On("Count", ctx, repository.CustomerFilters{Active: &activeFilter}).Return(8, nil)
	
	health, err := service.GetHealthStatus(ctx)
	
	assert.NoError(t, err)
	assert.NotNil(t, health)
	assert.Equal(t, "healthy", health.Status)
	assert.Equal(t, "customer-api", health.Service)
	assert.Equal(t, "test-v1.0.0", health.Version)
	assert.Equal(t, "test", health.Environment)
	assert.Contains(t, health.Dependencies, "repository")
	assert.Equal(t, "healthy", health.Dependencies["repository"])
	assert.Equal(t, 10, health.TotalCustomers)
	assert.Equal(t, 8, health.ActiveCustomers)
	mockRepo.AssertExpectations(t)
}

func TestCustomerService_GetHealthStatus_UnhealthyRepository(t *testing.T) {
	mockRepo := &MockCustomerRepository{}
	service := createTestCustomerService(mockRepo)
	ctx := context.Background()
	
	mockRepo.On("HealthCheck", ctx).Return(errors.New("database connection failed"))
	mockRepo.On("Count", ctx, repository.CustomerFilters{}).Return(0, nil)
	
	activeFilter := true
	mockRepo.On("Count", ctx, repository.CustomerFilters{Active: &activeFilter}).Return(0, nil)
	
	health, err := service.GetHealthStatus(ctx)
	
	assert.NoError(t, err)
	assert.NotNil(t, health)
	assert.Contains(t, health.Dependencies["repository"], "unhealthy")
	mockRepo.AssertExpectations(t)
}

func TestCustomerService_GetMetrics(t *testing.T) {
	service := createTestCustomerService(&MockCustomerRepository{})
	
	// Simulate some requests and errors
	service.requests = 150
	service.errors = 3
	
	metrics := service.GetMetrics()
	
	assert.NotNil(t, metrics)
	assert.Equal(t, "customer-api", metrics["service"])
	assert.Equal(t, "test-v1.0.0", metrics["version"])
	assert.Equal(t, "test", metrics["environment"])
	assert.Equal(t, int64(150), metrics["total_requests"])
	assert.Equal(t, int64(3), metrics["total_errors"])
	assert.InDelta(t, 0.02, metrics["error_rate"], 0.001)
	assert.InDelta(t, 0.98, metrics["success_rate"], 0.001)
}

func TestCustomerService_GetMetrics_NoRequests(t *testing.T) {
	service := createTestCustomerService(&MockCustomerRepository{})
	
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
