// +build integration

package repository

import (
	"context"
	"customer-service/internal/domain"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func setupTestDB(t *testing.T) (*mongo.Database, func()) {
	mongoURI := os.Getenv("TEST_MONGODB_URI")
	if mongoURI == "" {
		mongoURI = "mongodb://localhost:27017"
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoURI))
	require.NoError(t, err)

	db := client.Database("customer_test_db")

	cleanup := func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		_ = db.Drop(ctx)
		_ = client.Disconnect(ctx)
	}

	return db, cleanup
}

func TestMongoDBCustomerRepository_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	db, cleanup := setupTestDB(t)
	defer cleanup()

	repo := NewMongoDBCustomerRepository(db)
	ctx := context.Background()

	t.Run("Create and Find Customer", func(t *testing.T) {
		customer, err := domain.NewCustomer("John Doe", "11144477735", "john@example.com")
		require.NoError(t, err)

		// Create
		err = repo.Create(ctx, customer)
		assert.NoError(t, err)

		// Find by ID
		found, err := repo.FindByID(ctx, customer.ID)
		assert.NoError(t, err)
		assert.NotNil(t, found)
		assert.Equal(t, customer.Name, found.Name)
		assert.Equal(t, customer.CPF, found.CPF)
		assert.Equal(t, customer.Email, found.Email)

		// Find by CPF
		foundByCPF, err := repo.FindByCPF(ctx, customer.CPF)
		assert.NoError(t, err)
		assert.NotNil(t, foundByCPF)
		assert.Equal(t, customer.ID, foundByCPF.ID)
	})

	t.Run("Duplicate CPF should return error", func(t *testing.T) {
		customer1, _ := domain.NewCustomer("Jane Doe", "12345678909", "jane@example.com")
		err := repo.Create(ctx, customer1)
		require.NoError(t, err)

		customer2, _ := domain.NewCustomer("John Smith", "12345678909", "john.smith@example.com")
		err = repo.Create(ctx, customer2)
		assert.Error(t, err)
	})

	t.Run("Update Customer", func(t *testing.T) {
		customer, _ := domain.NewCustomer("Alice", "98765432100", "alice@example.com")
		err := repo.Create(ctx, customer)
		require.NoError(t, err)

		newName := "Alice Smith"
		newEmail := "alice.smith@example.com"
		err = customer.Update(&newName, &newEmail)
		require.NoError(t, err)

		err = repo.Update(ctx, customer)
		assert.NoError(t, err)

		updated, err := repo.FindByID(ctx, customer.ID)
		assert.NoError(t, err)
		assert.Equal(t, newName, updated.Name)
		assert.Equal(t, newEmail, updated.Email)
	})

	t.Run("Delete Customer", func(t *testing.T) {
		customer, _ := domain.NewCustomer("Bob", "45678912300", "bob@example.com")
		err := repo.Create(ctx, customer)
		require.NoError(t, err)

		err = repo.Delete(ctx, customer.ID)
		assert.NoError(t, err)

		found, err := repo.FindByID(ctx, customer.ID)
		assert.NoError(t, err)
		assert.Nil(t, found)
	})

	t.Run("FindByCPFOrEmail", func(t *testing.T) {
		customer, _ := domain.NewCustomer("Charlie", "15935745600", "charlie@example.com")
		err := repo.Create(ctx, customer)
		require.NoError(t, err)

		// Find by CPF
		found, err := repo.FindByCPFOrEmail(ctx, customer.CPF, "other@example.com")
		assert.NoError(t, err)
		assert.NotNil(t, found)

		// Find by Email
		found, err = repo.FindByCPFOrEmail(ctx, "00000000000", customer.Email)
		assert.NoError(t, err)
		assert.NotNil(t, found)

		// Not found
		found, err = repo.FindByCPFOrEmail(ctx, "00000000000", "notfound@example.com")
		assert.NoError(t, err)
		assert.Nil(t, found)
	})

	t.Run("GetEmailByID", func(t *testing.T) {
		customer, _ := domain.NewCustomer("David", "75315935700", "david@example.com")
		err := repo.Create(ctx, customer)
		require.NoError(t, err)

		email, err := repo.GetEmailByID(ctx, customer.ID)
		assert.NoError(t, err)
		assert.Equal(t, customer.Email, email)

		// Not found
		_, err = repo.GetEmailByID(ctx, "non-existent-id")
		assert.Error(t, err)
	})
}
