package domain

type User struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Age  int    `json:"age"`
}

type UserStore interface {
	CreateUser(name string, age int) (User, error)
	UpdateUser(u User) (User, error)
	DeleteUser(id int) error
	GetUser(id int) (User, error)
	GetAll(limit, page int, name string, age int) ([]User, error)
}

type UserService interface {
	GetAllUsers(limit, page int, name string, age int) ([]User, error)
	CreateUser(name string, age int) (User, error)
	DeleteUser(id int) error
	GetUser(id int) (User, error)
	UpdateUser(u User) (User, error)
}
