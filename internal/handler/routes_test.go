package handler

import (
	"customer-service/internal/usecase"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestSetupRoutes(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	mockRepo := new(MockRepository)
	handler := NewCustomerHandler(
		usecase.NewCreateCustomerUseCase(mockRepo),
		usecase.NewGetCustomerByCPFUseCase(mockRepo),
		usecase.NewUpdateCustomerUseCase(mockRepo),
		usecase.NewDeleteCustomerUseCase(mockRepo),
	)

	SetupRoutes(router, handler)

	routes := router.Routes()

	// Verify that all customer routes are registered
	expectedRoutes := map[string]string{
		"POST /customer":        "POST",
		"GET /customer/:cpf":    "GET",
		"PATCH /customer/:id":   "PATCH",
		"DELETE /customer/:id":  "DELETE",
	}

	routeMap := make(map[string]string)
	for _, route := range routes {
		key := route.Method + " " + route.Path
		routeMap[key] = route.Method
	}

	for expectedRoute, expectedMethod := range expectedRoutes {
		method, exists := routeMap[expectedRoute]
		assert.True(t, exists, "Route %s should exist", expectedRoute)
		assert.Equal(t, expectedMethod, method, "Route %s should have method %s", expectedRoute, expectedMethod)
	}

	assert.Equal(t, len(expectedRoutes), len(routes), "Should have exactly %d routes", len(expectedRoutes))
}
