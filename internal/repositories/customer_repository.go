package repositories

import (
	"context"
	"customer-service/internal/core/domain"
	"errors"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type customerRepository struct {
	Connection *gorm.DB
}

func NewCustomerRepository(db *gorm.DB) (*customerRepository, error) {
	err := db.AutoMigrate(&domain.Customer{})

	if err != nil {
		return nil, err
	}

	database := customerRepository{
		Connection: db,
	}

	return &database, nil
}

func (repository *customerRepository) Get(ctx context.Context, id string) (domain.Customer, error) {
	var customer domain.Customer

	repository.Connection.WithContext(ctx).Preload(clause.Associations).First(&customer, "user_id = ?", id)

	if (customer == domain.Customer{}) {
		return customer, errors.New("customer not found")
	}

	return customer, nil
}

func (repository *customerRepository) GetAll(ctx context.Context) ([]domain.Customer, error) {
	var customers []domain.Customer

	repository.Connection.WithContext(ctx).Find(&customers)

	return customers, nil
}

func (repository *customerRepository) Save(ctx context.Context, customer domain.Customer) (domain.Customer, error) {
	result := repository.Connection.WithContext(ctx).Omit("User").Create(&customer)

	if result.Error != nil {
		return domain.Customer{}, result.Error
	}

	return customer, nil
}

func (repository *customerRepository) Update(ctx context.Context, customer domain.Customer) (domain.Customer, error) {
	result := repository.Connection.WithContext(ctx).Model(&customer).Updates(customer)

	if result.Error != nil {
		return domain.Customer{}, result.Error
	}

	return customer, nil
}

func (repository *customerRepository) SaveOrUpdateUser(ctx context.Context, user domain.User) error {
	updateResult := repository.Connection.WithContext(ctx).Model(&user).Where("id = ?", user.ID).Updates(&user)

	if updateResult.RowsAffected == 0 {
		createResult := repository.Connection.WithContext(ctx).Create(&user)

		if createResult.Error != nil {
			return errors.New("could not create user")
		}
	}

	if updateResult.Error != nil {
		return errors.New("could not update user")
	}

	return nil
}

func (repository *customerRepository) GetUser(ctx context.Context, id string) (domain.User, error) {
	var user domain.User

	repository.Connection.WithContext(ctx).Preload(clause.Associations).First(&user, "id = ?", id)

	if (user == domain.User{}) {
		return user, errors.New("user not found")
	}

	return user, nil
}
