package services

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/pksingh21/go-echo-htmx/db"
)

func NewServicesUser(u User, uStore db.UserStore) *ServicesUser {

	return &ServicesUser{
		User:      u,
		UserStore: uStore,
	}
}

type User struct {
	ID        int       `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at,omitempty"`
}

type ServicesUser struct {
	User      User
	UserStore db.UserStore
}

func (su *ServicesUser) GetAllUsers() ([]User, error) {
	n := 10
	for i := 0; i < n; i++ {
		username := fmt.Sprintf("user%d", rand.Intn(1000))
		email := fmt.Sprintf("%s@example.com", username)
		createdAt := time.Now().Add(-time.Duration(rand.Intn(720)) * time.Hour) // Random date within last 30 days

		query := `INSERT INTO users (username, email, created_at) VALUES (?, ?, ?)`
		_, err := su.UserStore.Db.Query(query, username, email, createdAt)
		if err != nil {
			return []User{}, err
		}
	}
	query := `SELECT id, username, email, created_at FROM users ORDER BY created_at DESC`
	rows, err := su.UserStore.Db.Query(query)
	if err != nil {
		return []User{}, err
	}
	// We close the resource
	defer rows.Close()

	users := []User{}
	for rows.Next() {
		rows.Scan(
			&su.User.ID,
			&su.User.Username,
			&su.User.Email,
			&su.User.CreatedAt,
		)

		users = append(users, su.User)
	}

	return users, nil
}

func (su *ServicesUser) GetUserById(id int) (User, error) {

	query := `SELECT id, username, email, created_at FROM users
		WHERE id = ?`

	stmt, err := su.UserStore.Db.Prepare(query)
	if err != nil {
		return User{}, err
	}

	defer stmt.Close()

	su.User.ID = id
	err = stmt.QueryRow(
		su.User.ID,
	).Scan(
		&su.User.ID,
		&su.User.Username,
		&su.User.Email,
		&su.User.CreatedAt,
	)
	if err != nil {
		return User{}, err
	}

	return su.User, nil
}

func ConverDateTime(tz string, dt time.Time) string {
	loc, _ := time.LoadLocation(tz)

	return dt.In(loc).Format(time.RFC822Z)
}
