package models

import "time"

// Reservation represents a row in the "reservations" table.
type Reservation struct {
	ID            uint        `gorm:"primaryKey" json:"id"`
	UserID        uint        `gorm:"not null" json:"user_id"`
	ZoneID        uint        `gorm:"not null" json:"zone_id"`
	LicensePlate  string      `gorm:"type:varchar(15);not null" json:"license_plate"`
	Status        string      `gorm:"type:varchar(20);default:'active'" json:"status"` // active, completed, cancelled

	// These are GORM "associations" — not actual columns.
	// They let us fetch related data using Preload (like a JOIN).
	User User        `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Zone ParkingZone `gorm:"foreignKey:ZoneID" json:"zone,omitempty"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}