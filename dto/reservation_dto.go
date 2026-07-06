package dto

// CreateReservationRequest is what a driver sends to reserve a spot.
type CreateReservationRequest struct {
	ZoneID       uint   `json:"zone_id" validate:"required"`
	LicensePlate string `json:"license_plate" validate:"required,max=15"`
}

// ReservationResponse is the shape returned for a single reservation.
type ReservationResponse struct {
	ID           uint              `json:"id"`
	UserID       uint              `json:"user_id"`
	ZoneID       uint              `json:"zone_id"`
	LicensePlate string            `json:"license_plate"`
	Status       string            `json:"status"`
	Zone         *ZoneSummary      `json:"zone,omitempty"`
	User         *UserSummary      `json:"user,omitempty"`
	CreatedAt    string            `json:"created_at"`
}

// ZoneSummary is a lightweight zone shape used inside reservation responses.
type ZoneSummary struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
	Type string `json:"type"`
}

// UserSummary is a lightweight user shape used inside reservation responses (admin view).
type UserSummary struct {
	ID    uint   `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}