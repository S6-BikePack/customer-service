package riderhdl

import (
	"github.com/google/uuid"
	"rider-service/internal/core/domain"
	"rider-service/internal/core/ports"
)
import "github.com/gin-gonic/gin"

type HTTPHandler struct {
	riderService ports.RiderService
}

func NewHTTPHandler(riderService ports.RiderService) *HTTPHandler {
	return &HTTPHandler{
		riderService: riderService,
	}
}

// GetAll godoc
// @Summary  get all riders
// @Schemes
// @Description  gets all riders in the system
// @Accept       json
// @Produce      json
// @Success      200  {object}  []domain.Rider
// @Router       /riders [get]
func (hdl *HTTPHandler) GetAll(c *gin.Context) {
	riders, err := hdl.riderService.GetAll()

	if err != nil {
		c.AbortWithStatus(404)
		return
	}

	c.JSON(200, riders)
}

// Get godoc
// @Summary  get rider
// @Schemes
// @Param        id     path  string           true  "Rider id"
// @Description  gets a rider from the system by its ID
// @Produce      json
// @Success      200  {object}  domain.Rider
// @Router       /riders/{id} [get]
func (hdl *HTTPHandler) Get(c *gin.Context) {
	uid, err := uuid.Parse(c.Param("id"))

	if err != nil {
		c.AbortWithStatus(400)
		return
	}

	rider, err := hdl.riderService.Get(uid)

	if err != nil {
		c.AbortWithStatus(404)
		return
	}

	c.JSON(200, rider)
}

// Create godoc
// @Summary  create rider
// @Schemes
// @Description  creates a new rider
// @Accept       json
// @Param        rider  body  BodyCreate  true  "Add rider"
// @Produce      json
// @Success      200  {object}  ResponseCreate
// @Router       /riders [post]
func (hdl *HTTPHandler) Create(c *gin.Context) {
	body := BodyCreate{}
	err := c.BindJSON(&body)

	if err != nil {
		c.AbortWithStatus(500)
	}

	rider, err := hdl.riderService.Create(body.Name, body.Status)

	if err != nil {
		c.AbortWithStatusJSON(500, gin.H{"message": err.Error()})
		return
	}

	c.JSON(200, BuildResponseCreate(rider))
}

// UpdateRider godoc
// @Summary  update rider
// @Schemes
// @Description  updates a rider's information
// @Accept       json
// @Param        rider  body  BodyUpdate  true  "Update rider"
// @Param        id     path  string      true  "Rider id"
// @Produce      json
// @Success      200  {object}  ResponseUpdate
// @Router       /riders/{id} [put]
func (hdl *HTTPHandler) UpdateRider(c *gin.Context) {
	body := BodyUpdate{}
	err := c.BindJSON(&body)

	if err != nil {
		c.AbortWithStatus(500)
	}

	uid, err := uuid.Parse(c.Param("id"))

	if err != nil {
		c.AbortWithStatus(400)
		return
	}

	rider, err := hdl.riderService.Update(uid, body.Name, body.Status)

	if err != nil {
		c.AbortWithStatusJSON(500, gin.H{"message": err.Error()})
		return
	}

	c.JSON(200, BuildResponseUpdate(rider))
}

// UpdateLocation godoc
// @Summary  update rider location
// @Schemes
// @Description  updates a rider's location
// @Accept       json
// @Param        rider  body  domain.Location  true  "Update rider"
// @Param        id  path  string  true  "Rider id"
// @Produce      json
// @Success      200  {object}  ResponseUpdate
// @Router       /riders/{id}/location [put]
func (hdl *HTTPHandler) UpdateLocation(c *gin.Context) {
	body := domain.Location{}
	err := c.BindJSON(&body)

	if err != nil {
		c.AbortWithStatus(500)
	}

	uid, err := uuid.Parse(c.Param("id"))

	if err != nil {
		c.AbortWithStatus(400)
		return
	}

	rider, err := hdl.riderService.UpdateLocation(uid, body)

	if err != nil {
		c.AbortWithStatusJSON(500, gin.H{"message": err.Error()})
		return
	}

	c.JSON(200, BuildResponseUpdate(rider))
}
