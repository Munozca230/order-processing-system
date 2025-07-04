package repository

import (
	"context"
	"time"

	"github.com/product-api-v2/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// MongoProductRepository implements ProductRepository using MongoDB
type MongoProductRepository struct {
	collection *mongo.Collection
	client     *mongo.Client
}

// NewMongoProductRepository creates a new MongoDB product repository
func NewMongoProductRepository(mongoURL string) (*MongoProductRepository, error) {
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

	collection := client.Database("catalog").Collection("products")

	return &MongoProductRepository{
		collection: collection,
		client:     client,
	}, nil
}

// GetByID retrieves a product by its ID from MongoDB
func (r *MongoProductRepository) GetByID(ctx context.Context, productID string) (*models.Product, error) {
	var product models.Product
	
	filter := bson.M{"productId": productID}
	err := r.collection.FindOne(ctx, filter).Decode(&product)
	
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, ErrProductNotFound
		}
		return nil, err
	}
	
	return &product, nil
}

// GetAll retrieves all products with optional filtering from MongoDB
func (r *MongoProductRepository) GetAll(ctx context.Context, filters ProductFilters) ([]*models.Product, error) {
	// Build MongoDB filter
	filter := bson.M{}
	
	if filters.Category != "" {
		filter["category"] = filters.Category
	}
	
	if filters.Active != nil {
		filter["active"] = *filters.Active
	}
	
	if filters.MinPrice != nil || filters.MaxPrice != nil {
		priceFilter := bson.M{}
		if filters.MinPrice != nil {
			priceFilter["$gte"] = *filters.MinPrice
		}
		if filters.MaxPrice != nil {
			priceFilter["$lte"] = *filters.MaxPrice
		}
		filter["price"] = priceFilter
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
	
	var products []*models.Product
	for cursor.Next(ctx) {
		var product models.Product
		if err := cursor.Decode(&product); err != nil {
			return nil, err
		}
		products = append(products, &product)
	}
	
	if err := cursor.Err(); err != nil {
		return nil, err
	}
	
	return products, nil
}

// Create adds a new product to MongoDB
func (r *MongoProductRepository) Create(ctx context.Context, product *models.Product) error {
	// Check if product already exists
	existing, err := r.GetByID(ctx, product.ProductID)
	if err == nil && existing != nil {
		return ErrProductExists
	}
	
	product.CreatedAt = time.Now()
	product.UpdatedAt = time.Now()
	
	_, err = r.collection.InsertOne(ctx, product)
	return err
}

// Update modifies an existing product in MongoDB
func (r *MongoProductRepository) Update(ctx context.Context, product *models.Product) error {
	filter := bson.M{"productId": product.ProductID}
	
	product.UpdatedAt = time.Now()
	update := bson.M{"$set": product}
	
	result, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}
	
	if result.MatchedCount == 0 {
		return ErrProductNotFound
	}
	
	return nil
}

// Delete removes a product from MongoDB
func (r *MongoProductRepository) Delete(ctx context.Context, productID string) error {
	filter := bson.M{"productId": productID}
	
	result, err := r.collection.DeleteOne(ctx, filter)
	if err != nil {
		return err
	}
	
	if result.DeletedCount == 0 {
		return ErrProductNotFound
	}
	
	return nil
}

// Count returns the total number of products matching the filters
func (r *MongoProductRepository) Count(ctx context.Context, filters ProductFilters) (int, error) {
	filter := bson.M{}
	
	if filters.Category != "" {
		filter["category"] = filters.Category
	}
	
	if filters.Active != nil {
		filter["active"] = *filters.Active
	}
	
	if filters.MinPrice != nil || filters.MaxPrice != nil {
		priceFilter := bson.M{}
		if filters.MinPrice != nil {
			priceFilter["$gte"] = *filters.MinPrice
		}
		if filters.MaxPrice != nil {
			priceFilter["$lte"] = *filters.MaxPrice
		}
		filter["price"] = priceFilter
	}
	
	count, err := r.collection.CountDocuments(ctx, filter)
	return int(count), err
}

// HealthCheck verifies the MongoDB connection is working
func (r *MongoProductRepository) HealthCheck(ctx context.Context) error {
	return r.client.Ping(ctx, nil)
}

// Close closes the MongoDB connection
func (r *MongoProductRepository) Close(ctx context.Context) error {
	return r.client.Disconnect(ctx)
}