package handlers

import (
	"customer-service/config"
	"customer-service/internal/core/interfaces"
	"customer-service/pkg/authorization"
	"customer-service/pkg/dto"
	"customer-service/pkg/logging"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"
	"go.opentelemetry.io/otel/trace"
	"net/http"

	"customer-service/docs"
	_ "customer-service/docs"
)
import "github.com/gin-gonic/gin"

type HTTPHandler struct {
	customerService interfaces.CustomerService
	router          *gin.Engine
	logger          logging.Logger
	config          *config.Config
}

func NewRest(customerService interfaces.CustomerService, router *gin.Engine, logger logging.Logger, config *config.Config) *HTTPHandler {
	return &HTTPHandler{
		customerService: customerService,
		router:          router,
		config:          config,
		logger:          logger,
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
	docs.SwaggerInfo.Title = handler.config.Server.Service + " API"
	docs.SwaggerInfo.Description = handler.config.Server.Description

	handler.router.GET("/swagger/customer/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}

// GetAll godoc
// @Summary  get all customers
// @Schemes
// @Description  gets all customers in the system
// @Accept       json
// @Produce      json
// @Success      200  {object}  dto.CustomerListResponse
// @Router       /api/customers [get]
func (handler *HTTPHandler) GetAll(c *gin.Context) {
	ctx := c.Request.Context()
	span := trace.SpanFromContext(ctx)
	defer span.End()

	if authorization.NewRest(c).AuthorizeAdmin() {

		customers, err := handler.customerService.GetAll(ctx)

		if err != nil {
			c.AbortWithStatus(http.StatusNotFound)
			return
		}

		c.JSON(http.StatusOK, dto.CreateCustomerListResponse(customers))
		return
	}

	c.AbortWithStatus(http.StatusUnauthorized)
}

// Get godoc
// @Summary  get customer
// @Schemes
// @Param        id     path  string           true  "Customer id"
// @Description  gets a customer from the system by its ID
// @Produce      json
// @Success      200  {object}  dto.CustomerResponse
// @Router       /api/customers/{id} [get]
func (handler *HTTPHandler) Get(c *gin.Context) {
	ctx := c.Request.Context()
	span := trace.SpanFromContext(ctx)
	defer span.End()

	auth := authorization.NewRest(c)

	if auth.AuthorizeAdmin() || auth.AuthorizeMatchingId(c.Param("id")) {
		customer, err := handler.customerService.Get(ctx, c.Param("id"))

		if err != nil {
			c.AbortWithStatus(http.StatusNotFound)
			return
		}

		c.JSON(http.StatusOK, dto.CreateCustomerResponse(customer))
		return
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
// @Success      200  {object}  dto.CustomerResponse
// @Router       /api/customers [post]
func (handler *HTTPHandler) Create(c *gin.Context) {
	ctx := c.Request.Context()
	span := trace.SpanFromContext(ctx)
	defer span.End()

	body := dto.BodyCreateCustomer{}
	err := c.BindJSON(&body)

	if err != nil || body == (dto.BodyCreateCustomer{}) {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	auth := authorization.NewRest(c)

	if auth.AuthorizeAdmin() || auth.AuthorizeMatchingId(body.ID) {

		customer, err := handler.customerService.Create(ctx, body.ID, body.ServiceArea)

		if err != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
			handler.logger.Error(ctx, err.Error())
			return
		}

		c.JSON(http.StatusOK, dto.CreateCustomerResponse(customer))
		return
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
// @Success      200  {object}  dto.CustomerResponse
// @Router       /api/customers/{id}/service-area [put]
func (handler *HTTPHandler) UpdateServiceArea(c *gin.Context) {
	ctx := c.Request.Context()
	span := trace.SpanFromContext(ctx)
	defer span.End()

	auth := authorization.NewRest(c)

	if auth.AuthorizeAdmin() || auth.AuthorizeMatchingId(c.Param("id")) {

		body := dto.BodyUpdateServiceArea{}
		err := c.BindJSON(&body)

		if err != nil {
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}

		customer, err := handler.customerService.UpdateServiceArea(ctx, c.Param("id"), body.ServiceArea)

		if err != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
			handler.logger.Error(ctx, err.Error())
			return
		}

		c.JSON(http.StatusOK, dto.CreateCustomerResponse(customer))
		return
	}

	c.AbortWithStatus(http.StatusUnauthorized)
}
