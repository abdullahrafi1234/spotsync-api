package dto

// RegisterRequest is what the client sends when registering.
// validator tags define the rules automatically checked by Echo.
type RegisterRequest struct {
	Name     string `json:"name" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
	Role     string `json:"role" validate:"omitempty,oneof=driver admin"`
}

// LoginRequest is what the client sends when logging in.
type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

// UserResponse is the safe shape of a user we send back (no password!).
type UserResponse struct {
	ID        uint   `json:"id"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	Role      string `json:"role"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

// LoginResponse wraps the token + user info returned after login.
type LoginResponse struct {
	Token string       `json:"token"`
	User  UserResponse `json:"user"`
}