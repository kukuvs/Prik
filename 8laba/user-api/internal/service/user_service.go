package service

import (
	"user-api/internal/models"
	"user-api/internal/repository"
)

// UserService интерфейс бизнес-логики
type UserService interface {
	CreateUser(req *models.CreateUserRequest) (*models.User, error)
	GetUser(id int) (*models.User, error)
	GetUsers(page, pageSize int, filters map[string]interface{}) (*models.UserListResponse, error)
	UpdateUser(id int, req *models.UpdateUserRequest) (*models.User, error)
	DeleteUser(id int) error
}

type userService struct {
	repo repository.UserRepository
}

// NewUserService создает новый сервис пользователей
func NewUserService(repo repository.UserRepository) UserService {
	return &userService{repo: repo}
}

func (s *userService) CreateUser(req *models.CreateUserRequest) (*models.User, error) {
	return s.repo.Create(req)
}

func (s *userService) GetUser(id int) (*models.User, error) {
	return s.repo.GetByID(id)
}

func (s *userService) GetUsers(page, pageSize int, filters map[string]interface{}) (*models.UserListResponse, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	users, total, err := s.repo.GetAll(page, pageSize, filters)
	if err != nil {
		return nil, err
	}

	totalPages := (total + pageSize - 1) / pageSize

	return &models.UserListResponse{
		Users:      users,
		Total:      total,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: totalPages,
	}, nil
}

func (s *userService) UpdateUser(id int, req *models.UpdateUserRequest) (*models.User, error) {
	return s.repo.Update(id, req)
}

func (s *userService) DeleteUser(id int) error {
	return s.repo.Delete(id)
}
