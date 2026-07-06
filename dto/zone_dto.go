package dto

// CreateZoneRequest is what admin sends when creating a new parking zone.
type CreateZoneRequest struct {
	Name          string  `json:"name" validate:"required"`
	Type          string  `json:"type" validate:"required,oneof=general ev_charging covered"`
	TotalCapacity int     `json:"total_capacity" validate:"required,gt=0"`
	PricePerHour  float64 `json:"price_per_hour" validate:"required,gt=0"`
}

// ZoneResponse is the shape we return to clients, including calculated availability.
type ZoneResponse struct {
	ID             uint    `json:"id"`
	Name           string  `json:"name"`
	Type           string  `json:"type"`
	TotalCapacity  int     `json:"total_capacity"`
	AvailableSpots int     `json:"available_spots"`
	PricePerHour   float64 `json:"price_per_hour"`
	CreatedAt      string  `json:"created_at"`
}