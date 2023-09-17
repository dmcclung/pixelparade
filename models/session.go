package models

import (
	"crypto/sha256"
	"database/sql"
	"encoding/base64"
	"fmt"
	"log"

	"github.com/dmcclung/pixelparade/rand"
)

const (
	MinBytesPerToken = 32
)

type Session struct {
	ID        string
	UserID    string
	Token     string // Token is only set when creating new session
	TokenHash string
}

type SessionService struct {
	DB *sql.DB
	BytesPerToken int
}

const getUserByToken = `SELECT users.*
	FROM users
	JOIN sessions ON users.id = sessions.user_id
	WHERE sessions.token_hash = $1;`

const upsertSession = `INSERT INTO sessions (user_id, token_hash)
VALUES ($1, $2)
ON CONFLICT (user_id)
DO UPDATE SET token_hash = $2 RETURNING id;`

func (ss *SessionService) Create(userID string) (*Session, error) {
	bytesPerToken := ss.BytesPerToken
	if bytesPerToken < MinBytesPerToken {
		bytesPerToken = MinBytesPerToken
	}
	token, err := rand.String(bytesPerToken)
	if err != nil {
		return nil, fmt.Errorf("creating session: %w", err)
	}

	tokenHash := ss.hash(token)

	var id string
	err = ss.DB.QueryRow(upsertSession, userID, tokenHash).Scan(&id)
	if err != nil {
		return nil, fmt.Errorf("creating session: %w", err)
	}

	return &Session{
		ID: id,
		UserID: userID,
		Token: token,
		TokenHash: tokenHash,
	}, nil
}

func (ss *SessionService) Delete(token string) error {
	tokenHash := ss.hash(token)
	log.Printf("deleting token %v", tokenHash)
	_, err := ss.DB.Exec("delete from sessions where token_hash = $1", tokenHash)
	if err != nil {
		return err
	}

	return nil
}

func (ss *SessionService) User(token string) (*User, error) {
	tokenHash := ss.hash(token)	
	log.Printf("user token hash %v", tokenHash)
	user := User{}
	err := ss.DB.QueryRow(getUserByToken, tokenHash).Scan(&user.ID, &user.Email, &user.Password)
	if err != nil {
		return nil, fmt.Errorf("user from session: %w", err)
	}
	return &user, nil
}

func (ss *SessionService) hash(token string) string {
	tokenHash := sha256.Sum256([]byte(token))
	return base64.URLEncoding.EncodeToString(tokenHash[:])
}