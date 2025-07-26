package store

import (
	"database/sql"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID           int    `json:"id"`
	Username     string `json:"username"`
	PasswordHash string `json:"-"` // Don't include in JSON responses
	CreatedAt    string `json:"created_at"`
}

type PostgresUserStore struct {
	db *sql.DB
}

func NewPostgresUserStore(db *sql.DB) *PostgresUserStore {
	return &PostgresUserStore{db: db}
}

type UserStore interface {
	CreateUser(user *User) error
	GetUserByID(id int) (*User, error)
	GetUserByUsername(username string) (*User, error)
	HashPassword(password string) (string, error)
	CheckPassword(hashedPassword, password string) error
	AuthenticateUser(username, password string) (*User, error)
}

func (s *PostgresUserStore) CreateUser(user *User) error {
	query := `INSERT INTO users (username, password_hash) VALUES ($1, $2) RETURNING id, created_at`

	err := s.db.QueryRow(query, user.Username, user.PasswordHash).Scan(&user.ID, &user.CreatedAt)
	if err != nil {
		return err
	}
	return nil
}

func (s *PostgresUserStore) GetUserByID(id int) (*User, error) {
	query := `SELECT id, username, password_hash, created_at FROM users WHERE id = $1`
	user := &User{}
	err := s.db.QueryRow(query, id).Scan(&user.ID, &user.Username, &user.PasswordHash, &user.CreatedAt)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (s *PostgresUserStore) GetUserByUsername(username string) (*User, error) {
	query := `SELECT id, username, password_hash, created_at FROM users WHERE username = $1`
	user := &User{}
	err := s.db.QueryRow(query, username).Scan(&user.ID, &user.Username, &user.PasswordHash, &user.CreatedAt)
	if err != nil {
		return nil, err
	}
	return user, nil
}

// HashPassword hashes a plain text password using bcrypt
func (s *PostgresUserStore) HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

// CheckPassword verifies a plain text password against a hashed password
func (s *PostgresUserStore) CheckPassword(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

// AuthenticateUser verifies username and password, returns user if valid
func (s *PostgresUserStore) AuthenticateUser(username, password string) (*User, error) {
	user, err := s.GetUserByUsername(username)
	if err != nil {
		return nil, err // User not found or database error
	}

	err = s.CheckPassword(user.PasswordHash, password)
	if err != nil {
		return nil, err // Invalid password
	}

	return user, nil
}
