package repository

import (
	"context"
	"customer-service/internal/domain"
	"customer-service/pkg/errors"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoDBCustomerRepository struct {
	collection *mongo.Collection
}

func NewMongoDBCustomerRepository(db *mongo.Database) *MongoDBCustomerRepository {
	collection := db.Collection("customers")

	// Create unique indexes for CPF and Email
	_, _ = collection.Indexes().CreateMany(context.Background(), []mongo.IndexModel{
		{
			Keys:    bson.D{{Key: "cpf", Value: 1}},
			Options: options.Index().SetUnique(true),
		},
		{
			Keys:    bson.D{{Key: "email", Value: 1}},
			Options: options.Index().SetUnique(true),
		},
	})

	return &MongoDBCustomerRepository{
		collection: collection,
	}
}

func (r *MongoDBCustomerRepository) Create(ctx context.Context, customer *domain.Customer) error {
	_, err := r.collection.InsertOne(ctx, customer)
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			return errors.NewConflictError("Customer already exists", "CUSTOMER_ALREADY_EXISTS")
		}
		return errors.WrapError(err, "Failed to create customer")
	}
	return nil
}

func (r *MongoDBCustomerRepository) FindByID(ctx context.Context, id string) (*domain.Customer, error) {
	var customer domain.Customer
	err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&customer)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, errors.WrapError(err, "Failed to find customer by ID")
	}
	return &customer, nil
}

func (r *MongoDBCustomerRepository) FindByCPF(ctx context.Context, cpf string) (*domain.Customer, error) {
	var customer domain.Customer
	err := r.collection.FindOne(ctx, bson.M{"cpf": cpf}).Decode(&customer)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, errors.WrapError(err, "Failed to find customer by CPF")
	}
	return &customer, nil
}

func (r *MongoDBCustomerRepository) FindByCPFOrEmail(ctx context.Context, cpf, email string) (*domain.Customer, error) {
	var customer domain.Customer
	filter := bson.M{
		"$or": []bson.M{
			{"cpf": cpf},
			{"email": email},
		},
	}
	err := r.collection.FindOne(ctx, filter).Decode(&customer)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, errors.WrapError(err, "Failed to find customer by CPF or Email")
	}
	return &customer, nil
}

func (r *MongoDBCustomerRepository) Update(ctx context.Context, customer *domain.Customer) error {
	update := bson.M{
		"$set": bson.M{
			"name":      customer.Name,
			"email":     customer.Email,
			"updatedAt": customer.UpdatedAt,
		},
	}

	result, err := r.collection.UpdateOne(ctx, bson.M{"_id": customer.ID}, update)
	if err != nil {
		return errors.WrapError(err, "Failed to update customer")
	}

	if result.MatchedCount == 0 {
		return errors.NewNotFoundError("Customer not found", "CUSTOMER_NOT_FOUND")
	}

	return nil
}

func (r *MongoDBCustomerRepository) Delete(ctx context.Context, id string) error {
	result, err := r.collection.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		return errors.WrapError(err, "Failed to delete customer")
	}

	if result.DeletedCount == 0 {
		return errors.NewNotFoundError("Customer not found", "CUSTOMER_NOT_FOUND")
	}

	return nil
}

func (r *MongoDBCustomerRepository) GetEmailByID(ctx context.Context, id string) (string, error) {
	var result struct {
		Email string `bson:"email"`
	}

	opts := options.FindOne().SetProjection(bson.M{"email": 1})
	err := r.collection.FindOne(ctx, bson.M{"_id": id}, opts).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return "", errors.NewNotFoundError("Customer not found", "CUSTOMER_NOT_FOUND")
		}
		return "", errors.WrapError(err, "Failed to get customer email")
	}

	return result.Email, nil
}
