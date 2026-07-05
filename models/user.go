package models

import "time"

// User represents a row in the "users" table.
// GORM uses struct tags to know column constraints.
type User struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Name      string    `gorm:"not null" json:"name"`
	Email     string    `gorm:"unique;not null" json:"email"`
	Password  string    `gorm:"not null" json:"-"` // json:"-" means never include in JSON response
	Role      string    `gorm:"type:varchar(20);default:'driver'" json:"role"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}