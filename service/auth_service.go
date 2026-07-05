package service

import (
	"errors"
	"os"
	"time"

	"spotsync-api/dto"
	"spotsync-api/models"
	"spotsync-api/repository"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

// AuthService defines the business operations for authentication.
type AuthService interface {
	Register(req dto.RegisterRequest) (*models.User, error)
	Login(req dto.LoginRequest) (string, *models.User, error)
}

type authService struct {
	userRepo repository.UserRepository
}

// NewAuthService creates an AuthService, injecting the repository it depends on.
func NewAuthService(userRepo repository.UserRepository) AuthService {
	return &authService{userRepo: userRepo}
}

func (s *authService) Register(req dto.RegisterRequest) (*models.User, error) {
	// 1. Check if email is already taken
	existing, _ := s.userRepo.FindByEmail(req.Email)
	if existing != nil {
		return nil, errors.New("email already registered")
	}

	// 2. Hash the password (never store plain text!)
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), 10)
	if err != nil {
		return nil, err
	}

	// 3. Default role to "driver" if not provided
	role := req.Role
	if role == "" {
		role = "driver"
	}

	// 4. Build the user model and save it
	user := &models.User{
		Name:     req.Name,
		Email:    req.Email,
		Password: string(hashedPassword),
		Role:     role,
	}

	if err := s.userRepo.Create(user); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *authService) Login(req dto.LoginRequest) (string, *models.User, error) {
	// 1. Find user by email
	user, err := s.userRepo.FindByEmail(req.Email)
	if err != nil {
		return "", nil, errors.New("invalid email or password")
	}

	// 2. Compare provided password with stored hash
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		return "", nil, errors.New("invalid email or password")
	}

	// 3. Generate JWT token with user id + role in the payload
	claims := jwt.MapClaims{
		"user_id": user.ID,
		"role":    user.Role,
		"exp":     time.Now().Add(time.Hour * 24).Unix(), // expires in 24h
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		return "", nil, err
	}

	return signedToken, user, nil
}