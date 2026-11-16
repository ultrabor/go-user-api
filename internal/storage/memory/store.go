package memory

import (
	"log/slog"
	"sort"
	"sync"

	m "github.com/ultrabor/go-user-api/internal/domain"
)

/// create dirictory for user struct -- entity

type Store struct {
	users  []m.User
	mu     sync.Mutex
	nextID int
}

func New(logger *slog.Logger) *Store {
	logger.Info("local storage was succsecfully opened")
	return &Store{mu: sync.Mutex{}, nextID: 1}
}

/// added database, add errors to func,

func (s *Store) CreateUser(name string, age int) (m.User, error) { /// unic id
	s.mu.Lock()
	defer s.mu.Unlock()

	user := m.User{ID: s.nextID, Name: name, Age: age}

	s.nextID++
	s.users = append(s.users, user)

	return user, nil
}

//// added only name or age changing

func (s *Store) UpdateUser(u m.User) (m.User, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	index := -1
	for i, v := range s.users {
		if u.ID == v.ID {
			index = i
			break
		}
	}

	if index == -1 {
		return m.User{}, m.ErrUserNotFound
	}

	if u.Name != "" {
		s.users[index].Name = u.Name
	}

	if u.Age > 0 {
		s.users[index].Age = u.Age
	}

	return s.users[index], nil
}

func (s *Store) GetUser(id int) (m.User, error) {

	s.mu.Lock()
	defer s.mu.Unlock()

	for _, v := range s.users {
		if v.ID == id {
			return v, nil

		}
	}

	return m.User{}, m.ErrUserNotFound
}

func (s *Store) DeleteUser(id int) error {
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
		return m.ErrUserNotFound
	}

	s.users = append(s.users[:index], s.users[index+1:]...)

	return nil
}

func NewStore() *Store {
	return &Store{mu: sync.Mutex{}}
}

func (s *Store) GetAll(limit, page int, name string, age int) ([]m.User, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	users := []m.User{}

	for _, v := range s.users {
		if name != "" && name != v.Name {
			continue
		}
		if age > 0 && age != v.Age {
			continue
		}

		users = append(users, v)
	}

	// keep the same order as the postgres store (ORDER BY id)
	sort.Slice(users, func(i, j int) bool { return users[i].ID < users[j].ID })

	offset := (page - 1) * limit

	if offset >= len(users) {
		return []m.User{}, nil
	}

	end := offset + limit
	if end > len(users) {
		end = len(users)
	}

	return users[offset:end], nil
}
