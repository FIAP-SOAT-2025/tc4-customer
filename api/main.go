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
	"strings"
	"time"

	docs "customer-service/docs"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
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

	// Configure Swagger defaults from environment (can be overridden per-request)
	docs.SwaggerInfo.BasePath = getEnv("SWAGGER_BASEPATH", "/")
	defaultSchemes := []string{"http"}
	if s := os.Getenv("SWAGGER_SCHEMES"); s != "" {
		defaultSchemes = strings.Split(s, ",")
	} else if os.Getenv("ENV") == "production" {
		defaultSchemes = []string{"https"}
	}
	docs.SwaggerInfo.Schemes = defaultSchemes

	// Serve Swagger UI and set host/scheme dynamically from the incoming request
	router.GET("/swagger/*any", func(c *gin.Context) {
		// Prefer an explicit env override, otherwise use request Host
		host := os.Getenv("SWAGGER_HOST")
		if host == "" {
			host = c.Request.Host
		}
		docs.SwaggerInfo.Host = host

		// Determine scheme from X-Forwarded-Proto (common in AWS) or TLS
		scheme := "http"
		if proto := c.Request.Header.Get("X-Forwarded-Proto"); proto != "" {
			parts := strings.Split(proto, ",")
			scheme = strings.TrimSpace(parts[0])
		} else if c.Request.TLS != nil {
			scheme = "https"
		} else if len(docs.SwaggerInfo.Schemes) > 0 {
			scheme = docs.SwaggerInfo.Schemes[0]
		}
		docs.SwaggerInfo.Schemes = []string{scheme}

		ginSwagger.WrapHandler(swaggerFiles.Handler)(c)
	})

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
