package models

import (
	"database/sql"
	"fmt"
	"time"
)

type PasswordReset struct {
	ID     string
	UserID string
	// Token is only set when a PasswordReset is created
	Token     string
	TokenHash string
	ExpiresAt time.Time
}

const (
	DefaulResetDuration = 1 * time.Hour
)

type PasswordResetService struct {
	DB *sql.DB

	BytesPerToken int

	Duration time.Duration
}

func (ps *PasswordResetService) Create(email string) (*PasswordReset, error) {
	return nil, fmt.Errorf("TODO: Implement create")
}

func (ps *PasswordResetService) Validate(token string) (*User, error) {
	return nil, fmt.Errorf("TODO: Implement validate")
}
