package handler

import (
	"customer-service/internal/usecase"
	"customer-service/pkg/errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

type CustomerHandler struct {
	createUseCase   *usecase.CreateCustomerUseCase
	getByCPFUseCase *usecase.GetCustomerByCPFUseCase
	updateUseCase   *usecase.UpdateCustomerUseCase
	deleteUseCase   *usecase.DeleteCustomerUseCase
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

// CreateCustomer godoc
// @Summary Create a new customer
// @Description Create a new customer with name, cpf and email
// @Tags customers
// @Accept json
// @Produce json
// @Param customer body CreateCustomerRequest true "Customer to create"
// @Success 201 {object} domain.Customer
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /customer [post]
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

// GetCustomerByCPF godoc
// @Summary Get customer by CPF
// @Description Returns a customer identified by CPF
// @Tags customers
// @Produce json
// @Param cpf path string true "CPF"
// @Success 200 {object} domain.Customer
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /customer/{cpf} [get]
func (h *CustomerHandler) GetCustomerByCPF(c *gin.Context) {
	cpf := c.Param("cpf")

	customer, err := h.getByCPFUseCase.Execute(c.Request.Context(), cpf)
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, customer)
}

// UpdateCustomer godoc
// @Summary Update a customer
// @Description Update customer's name and/or email
// @Tags customers
// @Accept json
// @Produce json
// @Param id path string true "Customer ID"
// @Param customer body UpdateCustomerRequest true "Customer fields to update"
// @Success 200 {object} domain.Customer
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /customer/{id} [patch]
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

// DeleteCustomer godoc
// @Summary Delete a customer
// @Description Delete a customer by ID
// @Tags customers
// @Param id path string true "Customer ID"
// @Success 204
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /customer/{id} [delete]
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
