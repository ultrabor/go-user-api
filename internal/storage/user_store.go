package storage

import models "github.com/ultrabor/go-user-api/internal/models"

type UserStore interface {
	CreateUser(name string, age int) (models.User, error)
	UpdateUser(u models.User) (models.User, error)
	DeleteUser(id int) error
	GetUser(id int) (models.User, error)
	GetAll(limit, page int, name *string, age *int) ([]models.User, error)
}
