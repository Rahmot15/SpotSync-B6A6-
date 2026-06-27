package parkingzone

import (
	"gorm.io/gorm"
)

const activeReservationStatus = "active"

type ActiveReservationCount struct {
	ZoneID uint
	Count  int64
}

type Repository interface {
	Create(zone *ParkingZone) error
	FindAll() ([]ParkingZone, error)
	FindByID(id uint) (*ParkingZone, error)
	Update(zone *ParkingZone) error
	Delete(zone *ParkingZone) error
	CountActiveReservations(zoneID uint) (int64, error)
	CountActiveReservationsByZone() (map[uint]int64, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) Create(zone *ParkingZone) error {
	return r.db.Create(zone).Error
}

func (r *repository) FindAll() ([]ParkingZone, error) {
	var zones []ParkingZone
	if err := r.db.Order("id ASC").Find(&zones).Error; err != nil {
		return nil, err
	}

	return zones, nil
}

func (r *repository) FindByID(id uint) (*ParkingZone, error) {
	var zone ParkingZone
	if err := r.db.First(&zone, id).Error; err != nil {
		return nil, err
	}

	return &zone, nil
}

func (r *repository) Update(zone *ParkingZone) error {
	return r.db.Save(zone).Error
}

func (r *repository) Delete(zone *ParkingZone) error {
	return r.db.Delete(zone).Error
}

func (r *repository) CountActiveReservations(zoneID uint) (int64, error) {
	var count int64
	err := r.db.Table("reservations").
		Where("zone_id = ? AND status = ?", zoneID, activeReservationStatus).
		Count(&count).Error
	if err != nil {
		return 0, err
	}

	return count, nil
}

func (r *repository) CountActiveReservationsByZone() (map[uint]int64, error) {
	var rows []ActiveReservationCount
	err := r.db.Table("reservations").
		Select("zone_id, COUNT(*) AS count").
		Where("status = ?", activeReservationStatus).
		Group("zone_id").
		Scan(&rows).Error
	if err != nil {
		return nil, err
	}

	counts := make(map[uint]int64, len(rows))
	for _, row := range rows {
		counts[row.ZoneID] = row.Count
	}

	return counts, nil
}
