package models

import (
	"crypto/sha256"
	"database/sql"
	"encoding/base64"
	"fmt"
	"strings"
	"time"

	"github.com/dmcclung/pixelparade/rand"
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
	DB            *sql.DB
	BytesPerToken int
	Duration      time.Duration
}

func (ps *PasswordResetService) Create(email string) (*PasswordReset, error) {
	email = strings.ToLower(email)

	var userID string
	err := ps.DB.QueryRow(`SELECT id from users where email=$1;`, email).Scan(&userID)
	if err != nil {
		return nil, fmt.Errorf("creating password reset: %w", err)
	}

	bytesPerToken := ps.BytesPerToken
	if bytesPerToken < MinBytesPerToken {
		bytesPerToken = MinBytesPerToken
	}
	token, err := rand.String(bytesPerToken)
	if err != nil {
		return nil, fmt.Errorf("creating password reset token: %w", err)
	}

	duration := ps.Duration
	if duration == 0 {
		duration = DefaulResetDuration
	}

	reset := PasswordReset{
		UserID:    userID,
		Token:     token,
		TokenHash: ps.hash(token),
		ExpiresAt: time.Now().Add(duration),
	}

	err = ps.DB.QueryRow(`
		INSERT INTO password_resets (user_id, token_hash, expires_at)
		VALUES ($1, $2, $3) ON CONFLICT (user_id) DO
		UPDATE
		SET token_hash = $2, expires_at = $3
		RETURNING id;`, reset.UserID, reset.TokenHash, reset.ExpiresAt).Scan(&reset.ID)
	if err != nil {
		return nil, fmt.Errorf("password reset insert: %w", err)
	}

	return &reset, nil
}

func (ps *PasswordResetService) hash(token string) string {
	tokenHash := sha256.Sum256([]byte(token))
	return base64.URLEncoding.EncodeToString(tokenHash[:])
}

func (ps *PasswordResetService) delete(id string) error {
	_, err := ps.DB.Exec(`
		DELETE FROM password_resets
		WHERE id = $1;`, id)
	if err != nil {
		return fmt.Errorf("delete: %w", err)
	}
	return nil
}

func (ps *PasswordResetService) Validate(token string) (*User, error) {
	tokenHash := ps.hash(token)
	var user User
	var pwReset PasswordReset
	row := ps.DB.QueryRow(`
		SELECT password_resets.id,
			password_resets.expires_at,
			users.id,
			users.email,
			users.password
		FROM password_resets
			JOIN users ON users.id = password_resets.user_id
		WHERE password_resets.token_hash = $1;`, tokenHash)
	err := row.Scan(
		&pwReset.ID, &pwReset.ExpiresAt,
		&user.ID, &user.Email, &user.Password)
	if err != nil {
		return nil, fmt.Errorf("consume: %w", err)
	}
	if time.Now().After(pwReset.ExpiresAt) {
		return nil, fmt.Errorf("token expires: %v", token)
	}
	err = ps.delete(pwReset.ID)
	if err != nil {
		return nil, fmt.Errorf("consume: %w", err)
	}
	return &user, nil
}
