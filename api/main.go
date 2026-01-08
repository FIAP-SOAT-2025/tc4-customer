// @title Customer Service API
// @version 1.0
// @description API para gerenciamento de clientes
// @host localhost:8080
// @BasePath /
package main

import (
	"context"
	"customer-service/internal/handler"
	"customer-service/internal/repository"
	"customer-service/internal/usecase"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	ginSwagger "github.com/swaggo/gin-swagger"
	swaggerFiles "github.com/swaggo/files"
	_ "customer-service/docs"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	// Get configuration from environment variables
	mongoURI := getEnv("MONGODB_URI", "mongodb://localhost:27017")
	dbName := getEnv("MONGODB_DATABASE", "customer_db")
	port := getEnv("PORT", "8080")

	// Connect to MongoDB
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoURI))
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}
	defer func() {
		if err := client.Disconnect(context.Background()); err != nil {
			log.Printf("Error disconnecting from MongoDB: %v", err)
		}
	}()

	// Ping MongoDB to verify connection
	if err := client.Ping(ctx, nil); err != nil {
		log.Fatalf("Failed to ping MongoDB: %v", err)
	}

	log.Println("Successfully connected to MongoDB")

	// Initialize database
	db := client.Database(dbName)

	// Check if running seed command
	if len(os.Args) > 1 && os.Args[1] == "seed" {
		if err := runSeed(db); err != nil {
			log.Fatalf("Seed failed: %v", err)
		}
		log.Println("Seed completed successfully")
		return
	}

	// Initialize repository
	customerRepo := repository.NewMongoDBCustomerRepository(db)

	// Initialize use cases
	createUC := usecase.NewCreateCustomerUseCase(customerRepo)
	getByCPFUC := usecase.NewGetCustomerByCPFUseCase(customerRepo)
	updateUC := usecase.NewUpdateCustomerUseCase(customerRepo)
	deleteUC := usecase.NewDeleteCustomerUseCase(customerRepo)

	// Initialize handler
	customerHandler := handler.NewCustomerHandler(createUC, getByCPFUC, updateUC, deleteUC)

	// Setup Gin router
	router := gin.Default()

	// Health check endpoint
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "healthy",
			"service": "customer-service",
		})
	})

	// Setup routes
	handler.SetupRoutes(router, customerHandler)

	// Swagger UI
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Start server
	addr := fmt.Sprintf(":%s", port)
	log.Printf("Starting server on %s", addr)
	if err := router.Run(addr); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
