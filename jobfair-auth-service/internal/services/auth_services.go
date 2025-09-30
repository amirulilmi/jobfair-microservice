package services

import (
	"errors"
	// "time"

	"jobfair-auth-service/internal/models"
	"jobfair-auth-service/internal/repository"
	"jobfair-auth-service/internal/utils"
)

type AuthService struct {
	userRepo  *repository.UserRepository
	jwtSecret string
}

func NewAuthService(userRepo *repository.UserRepository, jwtSecret string) *AuthService {
	return &AuthService{
		userRepo:  userRepo,
		jwtSecret: jwtSecret,
	}
}

func (s *AuthService) Register(req *models.RegisterRequest) (*models.User, error) {
	// Check if user exists
	existingUser, _ := s.userRepo.GetByEmail(req.Email)
	if existingUser != nil {
		return nil, errors.New("user already exists")
	}

	// Hash password
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	user := &models.User{
		Email:    req.Email,
		Password: hashedPassword,
		UserType: req.UserType,
		IsActive: true,
	}

	return s.userRepo.Create(user)
}

func (s *AuthService) Login(req *models.LoginRequest) (*models.LoginResponse, error) {
	user, err := s.userRepo.GetByEmail(req.Email)
	if err != nil {
		return nil, errors.New("invalid credentials")
	}

	if !utils.CheckPassword(req.Password, user.Password) {
		return nil, errors.New("invalid credentials")
	}

	token, err := utils.GenerateToken(user.ID, string(user.UserType), s.jwtSecret)
	if err != nil {
		return nil, err
	}

	refreshToken, err := utils.GenerateRefreshToken(user.ID, s.jwtSecret)
	if err != nil {
		return nil, err
	}

	return &models.LoginResponse{
		Token:        token,
		RefreshToken: refreshToken,
		User:         *user,
	}, nil
}

func (s *AuthService) RefreshToken(refreshToken string) (string, error) {
	userID, err := utils.ValidateRefreshToken(refreshToken, s.jwtSecret)
	if err != nil {
		return "", err
	}

	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return "", err
	}

	return utils.GenerateToken(user.ID, string(user.UserType), s.jwtSecret)
}

func (s *AuthService) GetUserByID(id uint) (*models.User, error) {
	return s.userRepo.GetByID(id)
}

func (s *AuthService) GetAllUsers() ([]models.User, error) {
	return s.userRepo.GetAll()
}
