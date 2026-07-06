package service

import (
	"errors"

	"spotsync-api/dto"
	"spotsync-api/models"
	"spotsync-api/repository"
)

var ErrForbidden = errors.New("you are not allowed to perform this action")
var ErrNotFound = errors.New("reservation not found")

type ReservationService interface {
	Reserve(userID uint, req dto.CreateReservationRequest) (*models.Reservation, error)
	GetMyReservations(userID uint) ([]dto.ReservationResponse, error)
	GetAllReservations() ([]dto.ReservationResponse, error)
	Cancel(reservationID, requesterID uint, requesterRole string) error
}

type reservationService struct {
	reservationRepo repository.ReservationRepository
	zoneRepo        repository.ZoneRepository
}

func NewReservationService(reservationRepo repository.ReservationRepository, zoneRepo repository.ZoneRepository) ReservationService {
	return &reservationService{
		reservationRepo: reservationRepo,
		zoneRepo:        zoneRepo,
	}
}

// Reserve validates the zone exists, then delegates to the repository's
// locked transaction to safely create the reservation.
func (s *reservationService) Reserve(userID uint, req dto.CreateReservationRequest) (*models.Reservation, error) {
	// 1. Make sure the zone actually exists first (avoid locking a non-existent row)
	_, err := s.zoneRepo.FindByID(req.ZoneID)
	if err != nil {
		return nil, errors.New("parking zone not found")
	}

	// 2. Delegate to the repository, which handles the transaction + row lock
	reservation, err := s.reservationRepo.CreateWithLock(userID, req.ZoneID, req.LicensePlate)
	if err != nil {
		return nil, err
	}

	return reservation, nil
}

func (s *reservationService) toResponse(r models.Reservation) dto.ReservationResponse {
	res := dto.ReservationResponse{
		ID:           r.ID,
		UserID:       r.UserID,
		ZoneID:       r.ZoneID,
		LicensePlate: r.LicensePlate,
		Status:       r.Status,
		CreatedAt:    r.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}

	// Only include Zone/User if they were preloaded (ID != 0 means loaded)
	if r.Zone.ID != 0 {
		res.Zone = &dto.ZoneSummary{
			ID:   r.Zone.ID,
			Name: r.Zone.Name,
			Type: r.Zone.Type,
		}
	}
	if r.User.ID != 0 {
		res.User = &dto.UserSummary{
			ID:    r.User.ID,
			Name:  r.User.Name,
			Email: r.User.Email,
		}
	}

	return res
}

func (s *reservationService) GetMyReservations(userID uint) ([]dto.ReservationResponse, error) {
	reservations, err := s.reservationRepo.FindByUserID(userID)
	if err != nil {
		return nil, err
	}

	responses := make([]dto.ReservationResponse, 0)
	for _, r := range reservations {
		responses = append(responses, s.toResponse(r))
	}
	return responses, nil
}

func (s *reservationService) GetAllReservations() ([]dto.ReservationResponse, error) {
	reservations, err := s.reservationRepo.FindAll()
	if err != nil {
		return nil, err
	}

	responses := make([]dto.ReservationResponse, 0)
	for _, r := range reservations {
		responses = append(responses, s.toResponse(r))
	}
	return responses, nil
}

// Cancel enforces ownership: a driver can only cancel their OWN reservation,
// but an admin can cancel any reservation.
func (s *reservationService) Cancel(reservationID, requesterID uint, requesterRole string) error {
	reservation, err := s.reservationRepo.FindByID(reservationID)
	if err != nil {
		return ErrNotFound
	}

	if requesterRole != "admin" && reservation.UserID != requesterID {
		return ErrForbidden
	}

	return s.reservationRepo.UpdateStatus(reservationID, "cancelled")
}