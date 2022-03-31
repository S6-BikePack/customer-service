package repositories

import (
	"customer-service/internal/core/domain"
	"github.com/google/uuid"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type cockroachdb struct {
	Connection *gorm.DB
}

func NewCockroachDB(connStr string) (*cockroachdb, error) {
	db, err := gorm.Open(postgres.Open(connStr))

	if err != nil {
		return nil, err
	}

	err = db.AutoMigrate(&domain.Customer{})

	if err != nil {
		return nil, err
	}

	database := cockroachdb{
		Connection: db,
	}

	return &database, nil
}

func (repository *cockroachdb) Get(uid uuid.UUID) (domain.Customer, error) {
	var customer domain.Customer

	repository.Connection.Preload(clause.Associations).First(&customer, uid)

	return customer, nil
}

func (repository *cockroachdb) GetAll() ([]domain.Customer, error) {
	var customers []domain.Customer

	repository.Connection.Find(&customers)

	return customers, nil
}

func (repository *cockroachdb) Save(customer domain.Customer) (domain.Customer, error) {
	result := repository.Connection.Create(&customer)

	if result.Error != nil {
		return domain.Customer{}, result.Error
	}

	return customer, nil
}

func (repository *cockroachdb) Update(customer domain.Customer) (domain.Customer, error) {
	result := repository.Connection.Model(&customer).Updates(customer)

	if result.Error != nil {
		return domain.Customer{}, result.Error
	}

	return customer, nil
}
