package main

import (
	"context"
	"customer-service/internal/domain"
	"customer-service/internal/repository"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
)

func runSeed(db *mongo.Database) error {
	log.Println("Starting database seeding...")

	customerRepo := repository.NewMongoDBCustomerRepository(db)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	seedCustomers := []struct {
		name  string
		cpf   string
		email string
	}{
		{
			name:  "Ana Souza",
			cpf:   "52998224725",
			email: "ana.souza@email.com",
		},
		{
			name:  "Carlos Mendes",
			cpf:   "11144477735",
			email: "carlos.mendes@email.com",
		},
		{
			name:  "Beatriz Lima",
			cpf:   "98765432100",
			email: "beatriz.lima@email.com",
		},
	}

	successCount := 0
	for _, sc := range seedCustomers {
		// Check if customer already exists
		existing, _ := customerRepo.FindByCPF(ctx, sc.cpf)
		if existing != nil {
			log.Printf("Customer with CPF %s already exists, skipping...", sc.cpf)
			continue
		}

		// Create new customer
		customer, err := domain.NewCustomer(sc.name, sc.cpf, sc.email)
		if err != nil {
			log.Printf("Error creating customer %s: %v", sc.name, err)
			continue
		}

		// Save to database
		if err := customerRepo.Create(ctx, customer); err != nil {
			log.Printf("Error saving customer %s: %v", sc.name, err)
			continue
		}

		log.Printf("Successfully seeded customer: %s (CPF: %s)", sc.name, sc.cpf)
		successCount++
	}

	log.Printf("Database seeding completed. %d customers created.", successCount)
	return nil
}
