package handler

import (
	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.Engine, handler *CustomerHandler) {
	customerGroup := router.Group("/customer")
	{
		customerGroup.POST("", handler.CreateCustomer)
		customerGroup.GET("/:cpf", handler.GetCustomerByCPF)
		customerGroup.PATCH("/:id", handler.UpdateCustomer)
		customerGroup.DELETE("/:id", handler.DeleteCustomer)
	}
}
