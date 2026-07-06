package service

import (
	"spotsync-api/dto"
	"spotsync-api/models"
	"spotsync-api/repository"
)

type ZoneService interface {
	CreateZone(req dto.CreateZoneRequest) (*models.ParkingZone, error)
	GetAllZones() ([]dto.ZoneResponse, error)
	GetZoneByID(id uint) (*dto.ZoneResponse, error)
}

type zoneService struct {
	zoneRepo repository.ZoneRepository
}

func NewZoneService(zoneRepo repository.ZoneRepository) ZoneService {
	return &zoneService{zoneRepo: zoneRepo}
}

func (s *zoneService) CreateZone(req dto.CreateZoneRequest) (*models.ParkingZone, error) {
	zone := &models.ParkingZone{
		Name:          req.Name,
		Type:          req.Type,
		TotalCapacity: req.TotalCapacity,
		PricePerHour:  req.PricePerHour,
	}
	if err := s.zoneRepo.Create(zone); err != nil {
		return nil, err
	}
	return zone, nil
}

// toZoneResponse is a small helper that converts a model + active count
// into the response shape with available_spots calculated.
func (s *zoneService) toZoneResponse(zone models.ParkingZone) (dto.ZoneResponse, error) {
	activeCount, err := s.zoneRepo.CountActiveReservations(zone.ID)
	if err != nil {
		return dto.ZoneResponse{}, err
	}

	return dto.ZoneResponse{
		ID:             zone.ID,
		Name:           zone.Name,
		Type:           zone.Type,
		TotalCapacity:  zone.TotalCapacity,
		AvailableSpots: zone.TotalCapacity - int(activeCount),
		PricePerHour:   zone.PricePerHour,
		CreatedAt:      zone.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}, nil
}

func (s *zoneService) GetAllZones() ([]dto.ZoneResponse, error) {
	zones, err := s.zoneRepo.FindAll()
	if err != nil {
		return nil, err
	}

	responses := make([]dto.ZoneResponse, 0)
	for _, zone := range zones {
		res, err := s.toZoneResponse(zone)
		if err != nil {
			return nil, err
		}
		responses = append(responses, res)
	}
	return responses, nil
}

func (s *zoneService) GetZoneByID(id uint) (*dto.ZoneResponse, error) {
	zone, err := s.zoneRepo.FindByID(id)
	if err != nil {
		return nil, err
	}

	res, err := s.toZoneResponse(*zone)
	if err != nil {
		return nil, err
	}
	return &res, nil
}