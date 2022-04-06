package repositories

import (
	"customer-service/internal/core/domain"
	"errors"
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

func (repository *cockroachdb) Get(id string) (domain.Customer, error) {
	var customer domain.Customer

	repository.Connection.Preload(clause.Associations).First(&customer, "user_id = ?", id)

	if (customer == domain.Customer{}) {
		return customer, errors.New("customer not found")
	}

	return customer, nil
}

func (repository *cockroachdb) GetAll() ([]domain.Customer, error) {
	var customers []domain.Customer

	repository.Connection.Find(&customers)

	return customers, nil
}

func (repository *cockroachdb) Save(customer domain.Customer) (domain.Customer, error) {
	result := repository.Connection.Omit("User").Create(&customer)

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

func (repository *cockroachdb) SaveOrUpdateCustomer(user domain.User) error {
	updateResult := repository.Connection.Model(&user).Where("id = ?", user.ID).Updates(&user)

	if updateResult.RowsAffected == 0 {
		createResult := repository.Connection.Create(&user)

		if createResult.Error != nil {
			return errors.New("could not create user")
		}
	}

	if updateResult.Error != nil {
		return errors.New("could not update user")
	}

	return nil
}

func (repository *cockroachdb) GetUser(id string) (domain.User, error) {
	var user domain.User

	repository.Connection.Preload(clause.Associations).First(&user, "id = ?", id)

	if (user == domain.User{}) {
		return user, errors.New("user not found")
	}

	return user, nil
}
