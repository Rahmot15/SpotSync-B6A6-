package reservation

import (
	"errors"

	"spotsync/internal/apperror"
	"spotsync/internal/domain/parkingzone"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

const (
	StatusActive    = "active"
	StatusCompleted = "completed"
	StatusCancelled = "cancelled"

	zoneFullMessage = "Parking zone is full"
)

type Repository interface {
	ReserveSpot(userID, zoneID uint, licensePlate string) (*Reservation, error)
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

func (r *repository) ReserveSpot(userID, zoneID uint, licensePlate string) (*Reservation, error) {
	var createdReservation Reservation

	err := r.db.Transaction(func(tx *gorm.DB) error {
		var zone parkingzone.ParkingZone
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&zone, zoneID).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return apperror.NotFound("Parking zone not found")
			}

			return apperror.Internal("Failed to load parking zone")
		}

		var activeCount int64
		if err := tx.Model(&Reservation{}).
			Where("zone_id = ? AND status = ?", zoneID, StatusActive).
			Count(&activeCount).Error; err != nil {
			return apperror.Internal("Failed to count active reservations")
		}

		if activeCount >= int64(zone.TotalCapacity) {
			return apperror.Conflict(zoneFullMessage)
		}

		reservation := &Reservation{
			UserID:       userID,
			ZoneID:       zoneID,
			LicensePlate: licensePlate,
			Status:       StatusActive,
		}

		if err := tx.Create(reservation).Error; err != nil {
			return apperror.Internal("Failed to create reservation")
		}

		createdReservation = *reservation
		return nil
	})
	if err != nil {
		return nil, err
	}

	return &createdReservation, nil
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
