package storage

import (
	"sync"
)

type User struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Age  int    `json:"age"`
}

type Store struct {
	users []User
	mu    sync.Mutex
}

func (s *Store) CreateUser(name string, age int) User {
	s.mu.Lock()
	defer s.mu.Unlock()
	id := 1
	if len(s.users) > 0 {
		id = s.users[len(s.users)-1].ID + 1
	}

	user := User{ID: id, Name: name, Age: age}

	s.users = append(s.users, user)

	return user
}

func (s *Store) UpdateUser(id int, name string, age int) (User, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()

	index := -1
	for i, v := range s.users {
		if id == v.ID {
			index = i
			break
		}
	}

	if index == -1 {
		return User{}, false
	}

	s.users[index].Name = name
	s.users[index].Age = age

	return s.users[index], true
}

func (s *Store) GetUser(id int) (User, bool) {
	var user User
	var find bool = false
	s.mu.Lock()
	defer s.mu.Unlock()

	for _, v := range s.users {
		if v.ID == id {
			user = v
			find = true
			break
		}
	}

	return user, find
}

func (s *Store) DeleteUser(id int) (User, bool) {
	var user User
	var find bool = true
	s.mu.Lock()
	defer s.mu.Unlock()

	index := -1

	for i, v := range s.users {
		if v.ID == id {
			index = i
			break
		}
	}
	if index == -1 {
		return user, false
	}

	s.users = append(s.users[:index], s.users[index+1:]...)

	return user, find
}

func NewStore() *Store {
	return &Store{mu: sync.Mutex{}}
}

func (s *Store) GetAllUser(limit, page int, name string, age int) []User {
	s.mu.Lock()
	defer s.mu.Unlock()

	var users []User

	for _, v := range s.users {
		if name != "" && name != v.Name {
			continue
		}
		if age != -1 && age != v.Age {
			continue
		}

		users = append(users, v)
	}

	if limit <= 0 {
		limit = 10
	}
	if page <= 0 {
		page = 1
	}

	offset := (page - 1) * limit

	if offset >= len(users) {
		return []User{}
	}

	end := offset + limit
	if end > len(users) {
		end = len(users)
	}

	return users[offset:end]
}
