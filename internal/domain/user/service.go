package user

import (
	"errors"

	"spotsync/internal/apperror"
	"spotsync/internal/auth"
	"spotsync/internal/domain/user/dto"

	"gorm.io/gorm"
)

const defaultRole = "driver"

type Service interface {
	Register(req dto.RegisterRequest) (*dto.UserResponse, error)
	Login(req dto.LoginRequest) (*dto.LoginResponse, error)
}

type service struct {
	repository Repository
	jwtService *auth.JWTService
}

func NewService(repository Repository, jwtService *auth.JWTService) Service {
	return &service{
		repository: repository,
		jwtService: jwtService,
	}
}

func (s *service) Register(req dto.RegisterRequest) (*dto.UserResponse, error) {
	exists, err := s.repository.ExistsByEmail(req.Email)
	if err != nil {
		return nil, apperror.Internal("Failed to check user email")
	}
	if exists {
		return nil, apperror.BadRequest("Email already registered", nil)
	}

	hashedPassword, err := auth.HashPassword(req.Password)
	if err != nil {
		return nil, apperror.Internal("Failed to hash password")
	}

	role := req.Role
	if role == "" {
		role = defaultRole
	}

	user := &User{
		Name:     req.Name,
		Email:    req.Email,
		Password: hashedPassword,
		Role:     role,
	}

	if err := s.repository.Create(user); err != nil {
		return nil, apperror.Internal("Failed to register user")
	}

	response := toUserResponse(user)
	return &response, nil
}

func (s *service) Login(req dto.LoginRequest) (*dto.LoginResponse, error) {
	user, err := s.repository.FindByEmail(req.Email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperror.Unauthorized("Invalid email or password")
		}

		return nil, apperror.Internal("Failed to find user")
	}

	if !auth.CheckPassword(user.Password, req.Password) {
		return nil, apperror.Unauthorized("Invalid email or password")
	}

	token, err := s.jwtService.GenerateToken(user.ID, user.Role)
	if err != nil {
		return nil, apperror.Internal("Failed to generate token")
	}

	return &dto.LoginResponse{
		Token: token,
		User: dto.UserResponse{
			ID:    user.ID,
			Name:  user.Name,
			Email: user.Email,
			Role:  user.Role,
		},
	}, nil
}

func toUserResponse(user *User) dto.UserResponse {
	return dto.UserResponse{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		Role:      user.Role,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
}
