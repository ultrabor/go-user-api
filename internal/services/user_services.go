package services

import (
	"github.com/ultrabor/go-user-api/internal/domain"
)

type UserService struct {
	store domain.UserStore
}

func NewUserService(store domain.UserStore) *UserService {
	return &UserService{store: store}
}

func (s *UserService) GetAllUsers(limit, page int, name string, age int) ([]domain.User, error) {
	if limit <= 0 {
		limit = 10
	}
	if page <= 0 {
		page = 1
	}

	users, err := s.store.GetAll(limit, page, name, age)
	if err != nil {
		return nil, err
	}

	return users, nil
}

func (s *UserService) CreateUser(name string, age int) (domain.User, error) {
	return s.store.CreateUser(name, age)
}

func (s *UserService) DeleteUser(id int) error {
	return s.store.DeleteUser(id)
}

func (s *UserService) GetUser(id int) (domain.User, error) {
	return s.store.GetUser(id)
}

func (s *UserService) UpdateUser(u domain.User) (domain.User, error) {
	return s.store.UpdateUser(u)
}
