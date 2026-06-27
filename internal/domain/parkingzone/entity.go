package parkingzone

import "time"

type ParkingZone struct {
	ID            uint      `gorm:"primaryKey" json:"id"`
	Name          string    `gorm:"type:varchar(150);not null" json:"name"`
	Type          string    `gorm:"type:varchar(30);not null" json:"type"`
	TotalCapacity int       `gorm:"not null" json:"total_capacity"`
	PricePerHour  float64   `gorm:"type:decimal(10,2);not null" json:"price_per_hour"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

func (ParkingZone) TableName() string {
	return "parking_zones"
}
