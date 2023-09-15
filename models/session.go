package models

import (
	"crypto/sha256"
	"database/sql"
	"encoding/base64"
	"fmt"

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

const createSessionSql = `INSERT INTO sessions (user_id, token_hash) 
	VALUES ($1, $2) RETURNING id;`

const getUserByTokenSql = `SELECT users.*
	FROM users
	JOIN sessions ON users.id = sessions.user_id
	WHERE sessions.tokenHash = $1;`

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
	err = ss.DB.QueryRow(createSessionSql, userID, tokenHash).Scan(&id)
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

func (ss *SessionService) User(tokenHash string) (*User, error) {
	ss.DB.QueryRow(getUserByTokenSql, tokenHash)
	return nil, nil
}

func (ss *SessionService) hash(token string) string {
	tokenHash := sha256.Sum256([]byte(token))
	return base64.URLEncoding.EncodeToString(tokenHash[:])
}