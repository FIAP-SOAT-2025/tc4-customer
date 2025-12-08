package handler

import (
	"customer-service/internal/usecase"
	"customer-service/pkg/errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

type CustomerHandler struct {
	createUseCase      *usecase.CreateCustomerUseCase
	getByCPFUseCase    *usecase.GetCustomerByCPFUseCase
	updateUseCase      *usecase.UpdateCustomerUseCase
	deleteUseCase      *usecase.DeleteCustomerUseCase
}

func NewCustomerHandler(
	createUC *usecase.CreateCustomerUseCase,
	getByCP *usecase.GetCustomerByCPFUseCase,
	updateUC *usecase.UpdateCustomerUseCase,
	deleteUC *usecase.DeleteCustomerUseCase,
) *CustomerHandler {
	return &CustomerHandler{
		createUseCase:   createUC,
		getByCPFUseCase: getByCP,
		updateUseCase:   updateUC,
		deleteUseCase:   deleteUC,
	}
}

type CreateCustomerRequest struct {
	Name  string `json:"name" binding:"required"`
	CPF   string `json:"cpf" binding:"required"`
	Email string `json:"email" binding:"required,email"`
}

type UpdateCustomerRequest struct {
	Name  *string `json:"name,omitempty"`
	Email *string `json:"email,omitempty"`
}

func (h *CustomerHandler) CreateCustomer(c *gin.Context) {
	var req CreateCustomerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message":    "Invalid request body",
			"statusCode": 400,
			"error":      "INVALID_REQUEST",
		})
		return
	}

	customer, err := h.createUseCase.Execute(c.Request.Context(), req.Name, req.CPF, req.Email)
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusCreated, customer)
}

func (h *CustomerHandler) GetCustomerByCPF(c *gin.Context) {
	cpf := c.Param("cpf")

	customer, err := h.getByCPFUseCase.Execute(c.Request.Context(), cpf)
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, customer)
}

func (h *CustomerHandler) UpdateCustomer(c *gin.Context) {
	id := c.Param("id")

	var req UpdateCustomerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message":    "Invalid request body",
			"statusCode": 400,
			"error":      "INVALID_REQUEST",
		})
		return
	}

	customer, err := h.updateUseCase.Execute(c.Request.Context(), id, req.Name, req.Email)
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, customer)
}

func (h *CustomerHandler) DeleteCustomer(c *gin.Context) {
	id := c.Param("id")

	err := h.deleteUseCase.Execute(c.Request.Context(), id)
	if err != nil {
		handleError(c, err)
		return
	}

	c.Status(http.StatusNoContent)
}

func handleError(c *gin.Context, err error) {
	if appErr, ok := err.(*errors.AppError); ok {
		c.JSON(appErr.StatusCode, gin.H{
			"message":    appErr.Message,
			"statusCode": appErr.StatusCode,
			"error":      appErr.Code,
		})
		return
	}

	c.JSON(http.StatusInternalServerError, gin.H{
		"message":    "Internal server error",
		"statusCode": 500,
		"error":      "INTERNAL_ERROR",
	})
}
