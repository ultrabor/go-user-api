package postgres

import (
	"database/sql"
	"errors"
	"fmt"
	"log/slog"

	_ "github.com/lib/pq"
	m "github.com/ultrabor/go-user-api/internal/models"
)

type Store struct {
	db *sql.DB
}

func New(conn string, logger *slog.Logger) (*Store, error) {
	db, err := sql.Open("postgres", conn)
	if err != nil {
		logger.Error("error in opening postgres", slog.String("error", err.Error()))
		return nil, err
	}
	if err = db.Ping(); err != nil {
		logger.Error("error on pinging postgres", slog.String("error", err.Error()))
		return nil, err
	}

	logger.Info("database was succsecfully opened", "status", db.Stats().OpenConnections)
	return &Store{db: db}, nil
}

func (s *Store) CreateUser(name string, age int) (m.User, error) {
	var user m.User

	err := s.db.QueryRow(
		`INSERT INTO users (name, age)
		 VALUES ($1, $2)
		 RETURNING id, name, age`,
		name, age,
	).Scan(&user.ID, &user.Name, &user.Age)

	return user, err
}

func (s *Store) DeleteUser(id int) error {

	res, err := s.db.Exec(`DELETE FROM users where id = $1`, id)

	if err != nil {
		return err
	}

	n, _ := res.RowsAffected()

	if n == 0 {
		return errors.New("user not found")
	}

	return nil
}

func (s *Store) UpdateUser(u m.User) (m.User, error) {

	sqlStatement := `
UPDATE users
SET name = COALESCE(NULLIF($2, ''), name),
    age  = COALESCE(NULLIF($3, 0), age)
WHERE id = $1;`

	res, err := s.db.Exec(sqlStatement, u.ID, u.Name, u.Age)

	if err != nil {
		return m.User{}, err
	}

	n, err := res.RowsAffected()

	if err != nil {
		return m.User{}, err
	}

	if n == 0 {
		return m.User{}, errors.New("user not found")
	}

	user, err := s.GetUser(u.ID)
	if err != nil {
		return m.User{}, err
	}

	return user, nil
}

func (s *Store) GetUser(id int) (m.User, error) {
	row := s.db.QueryRow("SELECT id, name, age FROM users WHERE id = $1", id)

	var user m.User

	err := row.Scan(&user.ID, &user.Name, &user.Age)
	if err != nil {
		return m.User{}, err
	}
	return user, nil
}

func (s *Store) GetAll(limit, page int, name *string, age *int) ([]m.User, error) {

	offset := (page - 1) * limit

	query := "SELECT id, name, age FROM users WHERE 1=1"

	args := []interface{}{}
	argId := 1

	if name != nil && *name != "" {
		query += fmt.Sprintf(" AND name = $%d", argId)
		argId++
		args = append(args, *name)
	}

	if age != nil && *age > 0 {
		query += fmt.Sprintf(" AND age = $%d", argId)
		argId++
		args = append(args, *name)
	}

	query += fmt.Sprintf(" ORDER BY id LIMIT $%d OFFSET $%d", argId, argId+1)

	args = append(args, limit, offset)

	rows, err := s.db.Query(query, args...)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []m.User
	for rows.Next() {
		var user m.User
		err = rows.Scan(&user.ID, &user.Name, &user.Age)
		if err != nil {
			return nil, err
		}

		users = append(users, user)
	}

	return users, nil
}
