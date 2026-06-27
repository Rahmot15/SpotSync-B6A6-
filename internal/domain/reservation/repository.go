package reservation

import "gorm.io/gorm"

const (
	StatusActive    = "active"
	StatusCompleted = "completed"
	StatusCancelled = "cancelled"
)

type Repository interface {
	FindByID(id uint) (*Reservation, error)
	FindByIDWithRelations(id uint) (*Reservation, error)
	FindByUserID(userID uint) ([]Reservation, error)
	FindAll() ([]Reservation, error)
	Update(reservation *Reservation) error
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) FindByID(id uint) (*Reservation, error) {
	var reservation Reservation
	if err := r.db.First(&reservation, id).Error; err != nil {
		return nil, err
	}

	return &reservation, nil
}

func (r *repository) FindByIDWithRelations(id uint) (*Reservation, error) {
	var reservation Reservation
	if err := r.db.Preload("User").Preload("Zone").First(&reservation, id).Error; err != nil {
		return nil, err
	}

	return &reservation, nil
}

func (r *repository) FindByUserID(userID uint) ([]Reservation, error) {
	var reservations []Reservation
	if err := r.db.Preload("Zone").
		Where("user_id = ?", userID).
		Order("id DESC").
		Find(&reservations).Error; err != nil {
		return nil, err
	}

	return reservations, nil
}

func (r *repository) FindAll() ([]Reservation, error) {
	var reservations []Reservation
	if err := r.db.Preload("User").
		Preload("Zone").
		Order("id DESC").
		Find(&reservations).Error; err != nil {
		return nil, err
	}

	return reservations, nil
}

func (r *repository) Update(reservation *Reservation) error {
	return r.db.Save(reservation).Error
}
