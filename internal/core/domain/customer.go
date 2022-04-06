package domain

type Customer struct {
	UserID      string `gorm:"primaryKey"`
	User        User
	ServiceArea int
}

func NewCustomer(user User, serviceArea int) Customer {
	return Customer{
		UserID:      user.ID,
		User:        user,
		ServiceArea: serviceArea,
	}
}
