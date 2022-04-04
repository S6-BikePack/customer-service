package handlers

import (
	"customer-service/internal/core/ports"
	"customer-service/pkg/dto"
	"github.com/google/uuid"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"

	"customer-service/docs"
	_ "customer-service/docs"
)
import "github.com/gin-gonic/gin"

type HTTPHandler struct {
	customerService ports.CustomerService
	router          *gin.Engine
}

func NewRest(customerService ports.CustomerService, router *gin.Engine) *HTTPHandler {
	return &HTTPHandler{
		customerService: customerService,
		router:          router,
	}
}

func (handler *HTTPHandler) SetupEndpoints() {
	api := handler.router.Group("/api")
	api.GET("/customers", handler.GetAll)
	api.GET("/customers/:id", handler.Get)
	api.POST("/customers", handler.Create)
	api.PUT("/customers/:id", handler.UpdateDetails)
	api.PUT("/customers/:id/service-area", handler.UpdateServiceArea)
}

func (handler *HTTPHandler) SetupSwagger() {
	docs.SwaggerInfo.Title = "Customer service API"
	docs.SwaggerInfo.Description = "The customer service manages all customers for the BikePack system."

	handler.router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}

// GetAll godoc
// @Summary  get all customers
// @Schemes
// @Description  gets all customers in the system
// @Accept       json
// @Produce      json
// @Success      200  {object}  []domain.Customer
// @Router       /api/customers [get]
func (handler *HTTPHandler) GetAll(c *gin.Context) {
	customers, err := handler.customerService.GetAll()

	if err != nil {
		c.AbortWithStatus(404)
		return
	}

	c.JSON(200, customers)
}

// Get godoc
// @Summary  get customer
// @Schemes
// @Param        id     path  string           true  "Customer id"
// @Description  gets a customer from the system by its ID
// @Produce      json
// @Success      200  {object}  domain.Customer
// @Router       /api/customers/{id} [get]
func (handler *HTTPHandler) Get(c *gin.Context) {
	uid, err := uuid.Parse(c.Param("id"))

	if err != nil {
		c.AbortWithStatus(400)
		return
	}

	customer, err := handler.customerService.Get(uid)

	if err != nil {
		c.AbortWithStatus(404)
		return
	}

	c.JSON(200, customer)
}

// Create godoc
// @Summary  create customer
// @Schemes
// @Description  creates a new customer
// @Accept       json
// @Param        customer  body  dto.BodyCreateCustomer  true  "Add customer"
// @Produce      json
// @Success      200  {object}  dto.ResponseCreateCustomer
// @Router       /api/customers [post]
func (handler *HTTPHandler) Create(c *gin.Context) {
	body := dto.BodyCreateCustomer{}
	err := c.BindJSON(&body)

	if err != nil {
		c.AbortWithStatus(500)
	}

	customer, err := handler.customerService.Create(body.Name, body.LastName, body.Email)

	if err != nil {
		c.AbortWithStatusJSON(500, gin.H{"message": err.Error()})
		return
	}

	c.JSON(200, dto.BuildResponseCreateCustomer(customer))
}

// UpdateDetails godoc
// @Summary  update customer details
// @Schemes
// @Description  updates a customers name, last name and/or email
// @Accept       json
// @Param        customer  body  dto.BodyUpdateCustomer  true  "Update customer"
// @Param        id  path  string  true  "Customer id"
// @Produce      json
// @Success      200  {object}  dto.ResponseUpdateCustomer
// @Router       /api/customers/{id} [put]
func (handler *HTTPHandler) UpdateDetails(c *gin.Context) {
	body := dto.BodyUpdateCustomer{}
	err := c.BindJSON(&body)

	if err != nil {
		c.AbortWithStatus(500)
	}

	uid, err := uuid.Parse(c.Param("id"))

	if err != nil {
		c.AbortWithStatus(400)
		return
	}

	customer, err := handler.customerService.UpdateCustomerDetails(uid, body.Name, body.LastName, body.Email)

	if err != nil {
		c.AbortWithStatusJSON(500, gin.H{"message": err.Error()})
		return
	}

	c.JSON(200, dto.BuildResponseUpdateCustomer(customer))
}

// UpdateServiceArea godoc
// @Summary  update customers service area
// @Schemes
// @Description  updates a customers service area
// @Accept       json
// @Param        customer  body  dto.BodyUpdateServiceArea  true  "Update customer"
// @Param        id  path  string  true  "Customer id"
// @Produce      json
// @Success      200  {object}  dto.ResponseUpdateServiceArea
// @Router       /api/customers/{id}/service-area [put]
func (handler *HTTPHandler) UpdateServiceArea(c *gin.Context) {
	body := dto.BodyUpdateServiceArea{}
	err := c.BindJSON(&body)

	if err != nil {
		c.AbortWithStatus(500)
	}

	uid, err := uuid.Parse(c.Param("id"))

	if err != nil {
		c.AbortWithStatus(400)
		return
	}

	customer, err := handler.customerService.UpdateServiceArea(uid, body.ServiceArea)

	if err != nil {
		c.AbortWithStatusJSON(500, gin.H{"message": err.Error()})
		return
	}

	c.JSON(200, dto.BuildResponseUpdateServiceArea(customer))
}
