package reservation

import (
	"time"

	"spotsync/internal/domain/parkingzone"
	"spotsync/internal/domain/user"
)

type Reservation struct {
	ID           uint                    `gorm:"primaryKey" json:"id"`
	UserID       uint                    `gorm:"not null;index" json:"user_id"`
	ZoneID       uint                    `gorm:"not null;index" json:"zone_id"`
	LicensePlate string                  `gorm:"type:varchar(15);not null" json:"license_plate"`
	Status       string                  `gorm:"type:varchar(20);not null;default:active" json:"status"`
	User         user.User               `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT" json:"user,omitempty"`
	Zone         parkingzone.ParkingZone `gorm:"foreignKey:ZoneID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT" json:"zone,omitempty"`
	CreatedAt    time.Time               `json:"created_at"`
	UpdatedAt    time.Time               `json:"updated_at"`
}

func (Reservation) TableName() string {
	return "reservations"
}
