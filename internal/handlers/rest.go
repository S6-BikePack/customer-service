package handlers

import (
	"customer-service/internal/core/ports"
	"customer-service/pkg/authorization"
	"customer-service/pkg/dto"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"
	"net/http"

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
	api.PUT("/customers/:id/service-area", handler.UpdateServiceArea)
}

func (handler *HTTPHandler) SetupSwagger() {
	docs.SwaggerInfo.Title = "Customer service API"
	docs.SwaggerInfo.Description = "The customer service manages all customers for the BikePack system."

	handler.router.GET("/swagger/customer/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
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
	if authorization.NewRest(c).AuthorizeAdmin() {

		customers, err := handler.customerService.GetAll()

		if err != nil {
			c.AbortWithStatus(http.StatusNotFound)
			return
		}

		c.JSON(http.StatusOK, customers)
	}

	c.AbortWithStatus(http.StatusUnauthorized)
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
	auth := authorization.NewRest(c)

	if auth.AuthorizeAdmin() || auth.AuthorizeMatchingId(c.Param("id")) {
		customer, err := handler.customerService.Get(c.Param("id"))

		if err != nil {
			c.AbortWithStatus(http.StatusNotFound)
			return
		}

		c.JSON(http.StatusOK, customer)
	}

	c.AbortWithStatus(http.StatusUnauthorized)
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
		c.AbortWithStatus(http.StatusBadRequest)
	}

	auth := authorization.NewRest(c)

	if auth.AuthorizeAdmin() || auth.AuthorizeMatchingId(body.ID) {

		customer, err := handler.customerService.Create(body.ID, body.ServiceArea)

		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			return
		}

		c.JSON(http.StatusOK, dto.BuildResponseCreateCustomer(customer))
	}

	c.AbortWithStatus(http.StatusUnauthorized)
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
	auth := authorization.NewRest(c)

	if auth.AuthorizeAdmin() || auth.AuthorizeMatchingId(c.Param("id")) {

		body := dto.BodyUpdateServiceArea{}
		err := c.BindJSON(&body)

		if err != nil {
			c.AbortWithStatus(http.StatusBadRequest)
		}

		customer, err := handler.customerService.UpdateServiceArea(c.Param("id"), body.ServiceArea)

		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			return
		}

		c.JSON(http.StatusOK, dto.BuildResponseUpdateServiceArea(customer))
	}

	c.AbortWithStatus(http.StatusUnauthorized)
}
