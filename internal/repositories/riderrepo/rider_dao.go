package riderrepo

import (
	"gorm.io/gorm"
	"rider-service/internal/core/domain"
)

type RiderDao struct {
	gorm.Model
	Name      string
	Status    string
	Latitude  float64
	Longitude float64
}

func (RiderDao) TableName() string {
	return "rider"
}

func (dao RiderDao) ToDomain() domain.Rider {
	rider := domain.NewRider(dao.Name, dao.Status, domain.NewLocation(dao.Latitude, dao.Longitude))
	rider.ID = dao.Model.ID

	return rider
}

func (dao *RiderDao) FromDomain(rider domain.Rider) {
	dao.ID = rider.ID
	dao.Name = rider.Name
	dao.Status = rider.Status
	dao.Latitude = rider.Location.Latitude
	dao.Longitude = rider.Location.Longitude
}