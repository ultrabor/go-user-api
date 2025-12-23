package storage

import "sync"

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

func NewStore() *Store {
	return &Store{mu: sync.Mutex{}}
}
