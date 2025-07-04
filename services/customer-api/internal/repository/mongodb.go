package repository

import (
	"context"
	"time"

	"github.com/customer-api-v2/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// MongoCustomerRepository implements CustomerRepository using MongoDB
type MongoCustomerRepository struct {
	collection *mongo.Collection
	client     *mongo.Client
}

// NewMongoCustomerRepository creates a new MongoDB customer repository
func NewMongoCustomerRepository(mongoURL string) (*MongoCustomerRepository, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoURL))
	if err != nil {
		return nil, err
	}

	// Test connection
	if err := client.Ping(ctx, nil); err != nil {
		return nil, err
	}

	collection := client.Database("catalog").Collection("customers")

	return &MongoCustomerRepository{
		collection: collection,
		client:     client,
	}, nil
}

// GetByID retrieves a customer by their ID from MongoDB
func (r *MongoCustomerRepository) GetByID(ctx context.Context, customerID string) (*models.Customer, error) {
	var customer models.Customer
	
	filter := bson.M{"customerId": customerID}
	err := r.collection.FindOne(ctx, filter).Decode(&customer)
	
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, ErrCustomerNotFound
		}
		return nil, err
	}
	
	return &customer, nil
}

// GetAll retrieves all customers with optional filtering from MongoDB
func (r *MongoCustomerRepository) GetAll(ctx context.Context, filters CustomerFilters) ([]*models.Customer, error) {
	// Build MongoDB filter
	filter := bson.M{}
	
	if filters.Active != nil {
		filter["active"] = *filters.Active
	}
	
	if filters.CustomerTier != "" {
		filter["customerTier"] = filters.CustomerTier
	}
	
	if filters.Email != "" {
		filter["email"] = bson.M{"$regex": filters.Email, "$options": "i"}
	}
	
	// Build options for pagination
	opts := options.Find()
	if filters.PageSize > 0 {
		opts.SetLimit(int64(filters.PageSize))
		opts.SetSkip(int64(filters.Page * filters.PageSize))
	}
	
	cursor, err := r.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	
	var customers []*models.Customer
	for cursor.Next(ctx) {
		var customer models.Customer
		if err := cursor.Decode(&customer); err != nil {
			return nil, err
		}
		customers = append(customers, &customer)
	}
	
	if err := cursor.Err(); err != nil {
		return nil, err
	}
	
	return customers, nil
}

// Create adds a new customer to MongoDB
func (r *MongoCustomerRepository) Create(ctx context.Context, customer *models.Customer) error {
	// Check if customer already exists
	existing, err := r.GetByID(ctx, customer.CustomerID)
	if err == nil && existing != nil {
		return ErrCustomerExists
	}
	
	customer.CreatedAt = time.Now()
	customer.UpdatedAt = time.Now()
	
	_, err = r.collection.InsertOne(ctx, customer)
	return err
}

// Update modifies an existing customer in MongoDB
func (r *MongoCustomerRepository) Update(ctx context.Context, customer *models.Customer) error {
	filter := bson.M{"customerId": customer.CustomerID}
	
	customer.UpdatedAt = time.Now()
	update := bson.M{"$set": customer}
	
	result, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}
	
	if result.MatchedCount == 0 {
		return ErrCustomerNotFound
	}
	
	return nil
}

// Delete removes a customer from MongoDB
func (r *MongoCustomerRepository) Delete(ctx context.Context, customerID string) error {
	filter := bson.M{"customerId": customerID}
	
	result, err := r.collection.DeleteOne(ctx, filter)
	if err != nil {
		return err
	}
	
	if result.DeletedCount == 0 {
		return ErrCustomerNotFound
	}
	
	return nil
}

// Count returns the total number of customers matching the filters
func (r *MongoCustomerRepository) Count(ctx context.Context, filters CustomerFilters) (int, error) {
	filter := bson.M{}
	
	if filters.Active != nil {
		filter["active"] = *filters.Active
	}
	
	if filters.CustomerTier != "" {
		filter["customerTier"] = filters.CustomerTier
	}
	
	if filters.Email != "" {
		filter["email"] = bson.M{"$regex": filters.Email, "$options": "i"}
	}
	
	count, err := r.collection.CountDocuments(ctx, filter)
	return int(count), err
}

// HealthCheck verifies the MongoDB connection is working
func (r *MongoCustomerRepository) HealthCheck(ctx context.Context) error {
	return r.client.Ping(ctx, nil)
}

// Close closes the MongoDB connection
func (r *MongoCustomerRepository) Close(ctx context.Context) error {
	return r.client.Disconnect(ctx)
}