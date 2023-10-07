package models

import (
	"database/sql"
	"errors"
	"fmt"
	"log"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgerrcode"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID       string
	Email    string
	Password string
}

type UserService struct {
	DB *sql.DB
}

const (
	createUserSql = `INSERT INTO users (email, password) 
		VALUES ($1, $2) RETURNING id;`

	getUserSql = `SELECT id, email, password FROM users 
		WHERE email = $1`

	deleteUserSql = `DELETE FROM users 
		WHERE email = $1`

	updateUserSql = `UPDATE users SET email = $1, password = $2 
		WHERE email = $3`

	updatePasswordSql = `UPDATE users SET password = $1
		WHERE id = $2`
)

func (u *UserService) UpdatePassword(userID, password string) error {
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("updating password: %w", err)
	}
	passwordHash := string(hashedBytes)
	_, err = u.DB.Exec(updatePasswordSql, passwordHash, userID)
	if err != nil {
		return fmt.Errorf("sql update password: %w", err)
	}
	return nil
}

func (u *UserService) Authenticate(email, password string) (*User, error) {
	user, err := u.Get(email)
	if err != nil {
		return nil, fmt.Errorf("authenticate: %w", err)
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return nil, fmt.Errorf("authenticate: %w", err)
	}
	return user, nil
}

func (u *UserService) Create(email, password string) (*User, error) {
	h, err := HashPassword(password)
	if err != nil {
		return nil, fmt.Errorf("hashing password for create: %w", err)
	}

	var id string
	err = u.DB.QueryRow(createUserSql, email, h).Scan(&id)
	if err != nil {
		var pgError *pgconn.PgError
		if errors.As(err, &pgError) {
			if pgError.Code == pgerrcode.UniqueViolation {
				return nil, ErrEmailTaken
			}
		}
		return nil, fmt.Errorf("creating user: %w", err)
	}

	return &User{
		ID:       id,
		Email:    email,
		Password: h,
	}, nil
}

func (u *UserService) Get(email string) (*User, error) {
	var user User
	err := u.DB.QueryRow(getUserSql, email).Scan(&user.ID, &user.Email, &user.Password)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Fatalf("No user found: %v", err)
			return nil, fmt.Errorf("No user found: %w", err)
		}
		log.Fatalf("Querying user failed: %v", err)
		return nil, fmt.Errorf("querying user: %w", err)
	}
	return &user, nil
}

func (u *UserService) Delete(email string) error {
	result, err := u.DB.Exec(deleteUserSql, email)
	if err != nil {
		return fmt.Errorf("deleting user: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("getting affected rows: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("no user found with email: %s", email)
	}

	return nil
}

func (u *UserService) Update(currentEmail, email, password string) error {
	h, err := HashPassword(password)
	if err != nil {
		return fmt.Errorf("hashing password for update: %w", err)
	}

	result, err := u.DB.Exec(updateUserSql, email, h, currentEmail)
	if err != nil {
		return fmt.Errorf("updating user: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("getting affected rows: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("no user found with email: %s", currentEmail)
	}

	return nil
}
