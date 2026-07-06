package repository

import (
	"errors"

	"spotsync-api/models"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// ErrZoneFull is returned when a zone has no available capacity left.
var ErrZoneFull = errors.New("zone is full, no spots available")

type ReservationRepository interface {
	CreateWithLock(userID, zoneID uint, licensePlate string) (*models.Reservation, error)
	FindByUserID(userID uint) ([]models.Reservation, error)
	FindAll() ([]models.Reservation, error)
	FindByID(id uint) (*models.Reservation, error)
	UpdateStatus(id uint, status string) error
}

type reservationRepository struct {
	db *gorm.DB
}

func NewReservationRepository(db *gorm.DB) ReservationRepository {
	return &reservationRepository{db: db}
}

// CreateWithLock is the heart of the "EV Spot Bottleneck" solution.
// It uses a database transaction + row-level lock (SELECT ... FOR UPDATE)
// so that concurrent requests for the same zone are processed one at a time.
func (r *reservationRepository) CreateWithLock(userID, zoneID uint, licensePlate string) (*models.Reservation, error) {
	var reservation models.Reservation

	err := r.db.Transaction(func(tx *gorm.DB) error {
		var zone models.ParkingZone

		// 1. Lock the zone row. Any other transaction trying to lock the
		//    SAME row will wait here until this transaction commits/rollbacks.
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			First(&zone, zoneID).Error; err != nil {
			return err
		}

		// 2. Count current active reservations for this zone
		//    (this happens INSIDE the lock, so the count is guaranteed accurate)
		var activeCount int64
		if err := tx.Model(&models.Reservation{}).
			Where("zone_id = ? AND status = ?", zoneID, "active").
			Count(&activeCount).Error; err != nil {
			return err
		}

		// 3. Check capacity
		if int(activeCount) >= zone.TotalCapacity {
			return ErrZoneFull
		}

		// 4. Create the reservation (still inside the same transaction/lock)
		reservation = models.Reservation{
			UserID:       userID,
			ZoneID:       zoneID,
			LicensePlate: licensePlate,
			Status:       "active",
		}
		if err := tx.Create(&reservation).Error; err != nil {
			return err
		}

		return nil // commits the transaction
	})

	if err != nil {
		return nil, err
	}

	return &reservation, nil
}

func (r *reservationRepository) FindByUserID(userID uint) ([]models.Reservation, error) {
	var reservations []models.Reservation
	err := r.db.Preload("Zone").Where("user_id = ?", userID).Find(&reservations).Error
	return reservations, err
}

func (r *reservationRepository) FindAll() ([]models.Reservation, error) {
	var reservations []models.Reservation
	err := r.db.Preload("Zone").Preload("User").Find(&reservations).Error
	return reservations, err
}

func (r *reservationRepository) FindByID(id uint) (*models.Reservation, error) {
	var reservation models.Reservation
	err := r.db.First(&reservation, id).Error
	if err != nil {
		return nil, err
	}
	return &reservation, nil
}

func (r *reservationRepository) UpdateStatus(id uint, status string) error {
	return r.db.Model(&models.Reservation{}).Where("id = ?", id).Update("status", status).Error
}