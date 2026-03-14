package auth

import (
	"crypto/subtle"
	"database/sql"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/ozsari/velour/internal/models"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrUserExists         = errors.New("user already exists")
	ErrNoUsers            = errors.New("no users found")
)

type AuthService struct {
	db        *sql.DB
	jwtSecret []byte
}

func New(db *sql.DB, jwtSecret string) *AuthService {
	return &AuthService{
		db:        db,
		jwtSecret: []byte(jwtSecret),
	}
}

func (a *AuthService) InitDB() error {
	_, err := a.db.Exec(`
		CREATE TABLE IF NOT EXISTS users (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			username TEXT UNIQUE NOT NULL,
			password TEXT NOT NULL,
			is_admin BOOLEAN DEFAULT 0,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)
	`)
	return err
}

func (a *AuthService) NeedsSetup() (bool, error) {
	var count int
	err := a.db.QueryRow("SELECT COUNT(*) FROM users").Scan(&count)
	if err != nil {
		return true, err
	}
	return count == 0, nil
}

func (a *AuthService) CreateUser(username, password string, isAdmin bool) (*models.User, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	result, err := a.db.Exec(
		"INSERT INTO users (username, password, is_admin) VALUES (?, ?, ?)",
		username, string(hash), isAdmin,
	)
	if err != nil {
		return nil, ErrUserExists
	}

	id, _ := result.LastInsertId()
	return &models.User{
		ID:        id,
		Username:  username,
		IsAdmin:   isAdmin,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}, nil
}

func (a *AuthService) Login(username, password string) (string, *models.User, error) {
	var user models.User
	var hash string

	err := a.db.QueryRow(
		"SELECT id, username, password, is_admin, created_at, updated_at FROM users WHERE username = ?",
		username,
	).Scan(&user.ID, &user.Username, &hash, &user.IsAdmin, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return "", nil, ErrInvalidCredentials
	}

	if err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)); err != nil {
		return "", nil, ErrInvalidCredentials
	}

	token, err := a.generateToken(&user)
	if err != nil {
		return "", nil, err
	}

	return token, &user, nil
}

func (a *AuthService) ValidateToken(tokenStr string) (*models.User, error) {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return a.jwtSecret, nil
	})
	if err != nil || !token.Valid {
		return nil, errors.New("invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("invalid claims")
	}

	userID := int64(claims["user_id"].(float64))
	var user models.User
	err = a.db.QueryRow(
		"SELECT id, username, is_admin, created_at, updated_at FROM users WHERE id = ?",
		userID,
	).Scan(&user.ID, &user.Username, &user.IsAdmin, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return nil, errors.New("user not found")
	}

	return &user, nil
}

func (a *AuthService) generateToken(user *models.User) (string, error) {
	claims := jwt.MapClaims{
		"user_id":  user.ID,
		"username": user.Username,
		"is_admin": user.IsAdmin,
		"exp":      time.Now().Add(24 * time.Hour).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(a.jwtSecret)
}

// ResetPassword changes a user's password (used by CLI recovery and admin panel)
func (a *AuthService) ResetPassword(username, newPassword string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	result, err := a.db.Exec(
		"UPDATE users SET password = ?, updated_at = ? WHERE username = ?",
		string(hash), time.Now(), username,
	)
	if err != nil {
		return err
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		return ErrNoUsers
	}
	return nil
}

// ListUsers returns all users (for admin panel)
func (a *AuthService) ListUsers() ([]models.User, error) {
	rows, err := a.db.Query("SELECT id, username, is_admin, created_at, updated_at FROM users")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		var u models.User
		if err := rows.Scan(&u.ID, &u.Username, &u.IsAdmin, &u.CreatedAt, &u.UpdatedAt); err != nil {
			return nil, err
		}
		users = append(users, u)
	}
	return users, nil
}

// constant time comparison helper (unused but available)
func secureCompare(a, b string) bool {
	return subtle.ConstantTimeCompare([]byte(a), []byte(b)) == 1
}
