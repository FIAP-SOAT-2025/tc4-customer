package repository

import (
	"context"
	"customer-service/internal/domain"
	"customer-service/pkg/errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/integration/mtest"
)

func TestNewMongoDBCustomerRepository(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))

	mt.Run("Create repository", func(mt *mtest.T) {
		repo := NewMongoDBCustomerRepository(mt.DB)
		assert.NotNil(t, repo)
		assert.NotNil(t, repo.collection)
	})
}

func TestCreate(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))

	mt.Run("Successfully create customer", func(mt *mtest.T) {
		mt.AddMockResponses(mtest.CreateSuccessResponse())

		repo := &MongoDBCustomerRepository{collection: mt.Coll}
		customer, _ := domain.NewCustomer("John Doe", "11144477735", "john@example.com")

		err := repo.Create(context.Background(), customer)
		assert.NoError(t, err)
	})

	mt.Run("Duplicate key error", func(mt *mtest.T) {
		mt.AddMockResponses(mtest.CreateWriteErrorsResponse(mtest.WriteError{
			Index:   0,
			Code:    11000, // Duplicate key error code
			Message: "duplicate key error",
		}))

		repo := &MongoDBCustomerRepository{collection: mt.Coll}
		customer, _ := domain.NewCustomer("John Doe", "11144477735", "john@example.com")

		err := repo.Create(context.Background(), customer)
		assert.Error(t, err)
		appErr, ok := err.(*errors.AppError)
		assert.True(t, ok)
		assert.Equal(t, "CUSTOMER_ALREADY_EXISTS", appErr.Code)
	})

	mt.Run("Generic error", func(mt *mtest.T) {
		mt.AddMockResponses(mtest.CreateWriteErrorsResponse(mtest.WriteError{
			Index:   0,
			Code:    500,
			Message: "internal error",
		}))

		repo := &MongoDBCustomerRepository{collection: mt.Coll}
		customer, _ := domain.NewCustomer("John Doe", "11144477735", "john@example.com")

		err := repo.Create(context.Background(), customer)
		assert.Error(t, err)
	})
}

func TestFindByID(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))

	mt.Run("Successfully find customer", func(mt *mtest.T) {
		customer, _ := domain.NewCustomer("John Doe", "11144477735", "john@example.com")
		mt.AddMockResponses(mtest.CreateCursorResponse(1, "customer_db.customers", mtest.FirstBatch, bson.D{
			{Key: "_id", Value: customer.ID},
			{Key: "name", Value: customer.Name},
			{Key: "cpf", Value: customer.CPF},
			{Key: "email", Value: customer.Email},
		}))

		repo := &MongoDBCustomerRepository{collection: mt.Coll}
		result, err := repo.FindByID(context.Background(), customer.ID)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, customer.Name, result.Name)
	})

	mt.Run("Customer not found", func(mt *mtest.T) {
		mt.AddMockResponses(mtest.CreateCursorResponse(0, "customer_db.customers", mtest.FirstBatch))

		repo := &MongoDBCustomerRepository{collection: mt.Coll}
		result, err := repo.FindByID(context.Background(), "nonexistent")

		assert.NoError(t, err)
		assert.Nil(t, result)
	})

	mt.Run("Database error", func(mt *mtest.T) {
		mt.AddMockResponses(mtest.CreateCommandErrorResponse(mtest.CommandError{
			Code:    500,
			Message: "database error",
		}))

		repo := &MongoDBCustomerRepository{collection: mt.Coll}
		result, err := repo.FindByID(context.Background(), "123")

		assert.Error(t, err)
		assert.Nil(t, result)
	})
}

func TestFindByCPF(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))

	mt.Run("Successfully find customer", func(mt *mtest.T) {
		customer, _ := domain.NewCustomer("John Doe", "11144477735", "john@example.com")
		mt.AddMockResponses(mtest.CreateCursorResponse(1, "customer_db.customers", mtest.FirstBatch, bson.D{
			{Key: "_id", Value: customer.ID},
			{Key: "name", Value: customer.Name},
			{Key: "cpf", Value: customer.CPF},
			{Key: "email", Value: customer.Email},
		}))

		repo := &MongoDBCustomerRepository{collection: mt.Coll}
		result, err := repo.FindByCPF(context.Background(), "11144477735")

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, customer.CPF, result.CPF)
	})

	mt.Run("Customer not found", func(mt *mtest.T) {
		mt.AddMockResponses(mtest.CreateCursorResponse(0, "customer_db.customers", mtest.FirstBatch))

		repo := &MongoDBCustomerRepository{collection: mt.Coll}
		result, err := repo.FindByCPF(context.Background(), "11144477735")

		assert.NoError(t, err)
		assert.Nil(t, result)
	})

	mt.Run("Database error", func(mt *mtest.T) {
		mt.AddMockResponses(mtest.CreateCommandErrorResponse(mtest.CommandError{
			Code:    500,
			Message: "database error",
		}))

		repo := &MongoDBCustomerRepository{collection: mt.Coll}
		result, err := repo.FindByCPF(context.Background(), "11144477735")

		assert.Error(t, err)
		assert.Nil(t, result)
	})
}

func TestFindByCPFOrEmail(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))

	mt.Run("Successfully find customer", func(mt *mtest.T) {
		customer, _ := domain.NewCustomer("John Doe", "11144477735", "john@example.com")
		mt.AddMockResponses(mtest.CreateCursorResponse(1, "customer_db.customers", mtest.FirstBatch, bson.D{
			{Key: "_id", Value: customer.ID},
			{Key: "name", Value: customer.Name},
			{Key: "cpf", Value: customer.CPF},
			{Key: "email", Value: customer.Email},
		}))

		repo := &MongoDBCustomerRepository{collection: mt.Coll}
		result, err := repo.FindByCPFOrEmail(context.Background(), "11144477735", "john@example.com")

		assert.NoError(t, err)
		assert.NotNil(t, result)
	})

	mt.Run("Customer not found", func(mt *mtest.T) {
		mt.AddMockResponses(mtest.CreateCursorResponse(0, "customer_db.customers", mtest.FirstBatch))

		repo := &MongoDBCustomerRepository{collection: mt.Coll}
		result, err := repo.FindByCPFOrEmail(context.Background(), "11144477735", "john@example.com")

		assert.NoError(t, err)
		assert.Nil(t, result)
	})

	mt.Run("Database error", func(mt *mtest.T) {
		mt.AddMockResponses(mtest.CreateCommandErrorResponse(mtest.CommandError{
			Code:    500,
			Message: "database error",
		}))

		repo := &MongoDBCustomerRepository{collection: mt.Coll}
		result, err := repo.FindByCPFOrEmail(context.Background(), "11144477735", "john@example.com")

		assert.Error(t, err)
		assert.Nil(t, result)
	})
}

func TestUpdate(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))

	mt.Run("Successfully update customer", func(mt *mtest.T) {
		mt.AddMockResponses(mtest.CreateSuccessResponse(
			bson.E{Key: "n", Value: 1},
			bson.E{Key: "nModified", Value: 1},
		))

		repo := &MongoDBCustomerRepository{collection: mt.Coll}
		customer, _ := domain.NewCustomer("John Doe", "11144477735", "john@example.com")

		err := repo.Update(context.Background(), customer)
		assert.NoError(t, err)
	})

	mt.Run("Customer not found", func(mt *mtest.T) {
		mt.AddMockResponses(mtest.CreateSuccessResponse(
			bson.E{Key: "n", Value: 0},
		))

		repo := &MongoDBCustomerRepository{collection: mt.Coll}
		customer, _ := domain.NewCustomer("John Doe", "11144477735", "john@example.com")

		err := repo.Update(context.Background(), customer)
		assert.Error(t, err)
		appErr, ok := err.(*errors.AppError)
		assert.True(t, ok)
		assert.Equal(t, "CUSTOMER_NOT_FOUND", appErr.Code)
	})

	mt.Run("Database error", func(mt *mtest.T) {
		mt.AddMockResponses(mtest.CreateWriteErrorsResponse(mtest.WriteError{
			Index:   0,
			Code:    500,
			Message: "database error",
		}))

		repo := &MongoDBCustomerRepository{collection: mt.Coll}
		customer, _ := domain.NewCustomer("John Doe", "11144477735", "john@example.com")

		err := repo.Update(context.Background(), customer)
		assert.Error(t, err)
	})
}

func TestDelete(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))

	mt.Run("Successfully delete customer", func(mt *mtest.T) {
		mt.AddMockResponses(mtest.CreateSuccessResponse(
			bson.E{Key: "n", Value: 1},
		))

		repo := &MongoDBCustomerRepository{collection: mt.Coll}
		err := repo.Delete(context.Background(), "123")

		assert.NoError(t, err)
	})

	mt.Run("Customer not found", func(mt *mtest.T) {
		mt.AddMockResponses(mtest.CreateSuccessResponse(
			bson.E{Key: "n", Value: 0},
		))

		repo := &MongoDBCustomerRepository{collection: mt.Coll}
		err := repo.Delete(context.Background(), "123")

		assert.Error(t, err)
		appErr, ok := err.(*errors.AppError)
		assert.True(t, ok)
		assert.Equal(t, "CUSTOMER_NOT_FOUND", appErr.Code)
	})

	mt.Run("Database error", func(mt *mtest.T) {
		mt.AddMockResponses(mtest.CreateWriteErrorsResponse(mtest.WriteError{
			Index:   0,
			Code:    500,
			Message: "database error",
		}))

		repo := &MongoDBCustomerRepository{collection: mt.Coll}
		err := repo.Delete(context.Background(), "123")

		assert.Error(t, err)
	})
}

func TestGetEmailByID(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))

	mt.Run("Successfully get email", func(mt *mtest.T) {
		mt.AddMockResponses(mtest.CreateCursorResponse(1, "customer_db.customers", mtest.FirstBatch, bson.D{
			{Key: "email", Value: "john@example.com"},
		}))

		repo := &MongoDBCustomerRepository{collection: mt.Coll}
		email, err := repo.GetEmailByID(context.Background(), "123")

		assert.NoError(t, err)
		assert.Equal(t, "john@example.com", email)
	})

	mt.Run("Customer not found", func(mt *mtest.T) {
		mt.AddMockResponses(mtest.CreateCursorResponse(0, "customer_db.customers", mtest.FirstBatch))

		repo := &MongoDBCustomerRepository{collection: mt.Coll}
		email, err := repo.GetEmailByID(context.Background(), "123")

		assert.Error(t, err)
		assert.Empty(t, email)
		appErr, ok := err.(*errors.AppError)
		assert.True(t, ok)
		assert.Equal(t, "CUSTOMER_NOT_FOUND", appErr.Code)
	})

	mt.Run("Database error", func(mt *mtest.T) {
		mt.AddMockResponses(mtest.CreateCommandErrorResponse(mtest.CommandError{
			Code:    500,
			Message: "database error",
		}))

		repo := &MongoDBCustomerRepository{collection: mt.Coll}
		email, err := repo.GetEmailByID(context.Background(), "123")

		assert.Error(t, err)
		assert.Empty(t, email)
	})
}
